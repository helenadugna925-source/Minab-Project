package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"local-event-backend/utils"
)

type PurchaseTicketRequest struct {
	EventID  int     `json:"event_id" binding:"required"`
	TicketID int     `json:"ticket_id" binding:"required"`
	Quantity int     `json:"quantity" binding:"required"`
	Amount   float64 `json:"amount" binding:"required"`
}

type PurchaseTicketResponse struct {
	Status        string    `json:"status"`
	ReservationID int       `json:"reservation_id"`
	CheckoutURL   string    `json:"checkout_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

func PurchaseTicketHandler(c *gin.Context) {
	userID, ok := getUserIDFromAuthHeaderTicket(c)
	if !ok {
		return
	}

	var req PurchaseTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	if req.Quantity <= 0 || req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Quantity/Amount must be greater than zero"})
		return
	}

	var reservationID int
	var userEmail string
	_ = utils.DB.QueryRow(`SELECT email FROM users WHERE id=$1`, userID).Scan(&userEmail)

	err := utils.DB.QueryRow(`
		INSERT INTO ticket_sales (event_id, ticket_id, user_id, quantity, amount, status, created_at)
		VALUES ($1, $2, $3, $4, $5, 'pending', NOW())
		RETURNING id
	`, req.EventID, req.TicketID, userID, req.Quantity, req.Amount).Scan(&reservationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	checkoutURL, txRef, err := initiateChapaPayment(req.Amount, userEmail, fmt.Sprintf("user-%d", userID), "buyer", reservationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	_, _ = utils.DB.Exec(`UPDATE ticket_sales SET tx_ref=$1, checkout_url=$2 WHERE id=$3`, txRef, checkoutURL, reservationID)

	c.JSON(http.StatusOK, PurchaseTicketResponse{
		Status:        "pending",
		ReservationID: reservationID,
		CheckoutURL:   checkoutURL,
		CreatedAt:     time.Now(),
	})
}

func getUserIDFromAuthHeaderTicket(c *gin.Context) (int, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing Authorization header"})
		return 0, false
	}
	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtSecret, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return 0, false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
		return 0, false
	}
	hasura, ok := claims["https://hasura.io/jwt/claims"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing Hasura claims"})
		return 0, false
	}
	userIDStr, ok := hasura["x-hasura-user-id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing user id claim"})
		return 0, false
	}
	return mustAtoi(userIDStr), true
}

type chapaInitResponse struct {
	Status string `json:"status"`
	Data   struct {
		CheckoutURL string `json:"checkout_url"`
		Reference   string `json:"reference"`
	} `json:"data"`
}

func initiateChapaPayment(amount float64, email, firstName, lastName string, reservationID int) (string, string, error) {
	secret := os.Getenv("CHAPA_SECRET")
	callbackURL := os.Getenv("CHAPA_CALLBACK_URL")
	returnURL := os.Getenv("CHAPA_RETURN_URL")
	if secret == "" {
		return "", "", fmt.Errorf("CHAPA_SECRET is not set")
	}
	if callbackURL == "" {
		callbackURL = "http://localhost:8082/chapa/webhook"
	}
	if returnURL == "" {
		returnURL = "http://localhost:3000/events/reserved"
	}

	txRef := fmt.Sprintf("tx-%d-%d", reservationID, time.Now().Unix())
	payload := map[string]interface{}{
		"amount":       fmt.Sprintf("%.2f", amount),
		"currency":     "ETB",
		"email":        email,
		"first_name":   firstName,
		"last_name":    lastName,
		"tx_ref":       txRef,
		"callback_url": callbackURL,
		"return_url":   returnURL,
		"customization": map[string]string{
			"title":       "Event Ticket",
			"description": "Local event ticket purchase",
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.chapa.co/v1/transaction/initialize", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+secret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return "", "", fmt.Errorf("chapa init failed: %s", res.Status)
	}

	var resp chapaInitResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return "", "", err
	}
	if resp.Status != "success" {
		return "", "", fmt.Errorf("chapa init status: %s", resp.Status)
	}
	return resp.Data.CheckoutURL, resp.Data.Reference, nil
}

