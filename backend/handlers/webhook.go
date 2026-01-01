package handlers

import (
	"fmt"
	"local-event-backend/utils"
	"net/http"
	"strconv" 
	"strings"

	"github.com/gin-gonic/gin"
)

type ChapaWebhookPayload struct {
	Status    string  `json:"status"`
	Amount    float64 `json:"amount"`
	TxRef     string  `json:"tx_ref"`
	Reference string  `json:"reference"`
}

func ChapaWebhookHandler(c *gin.Context) {
	var payload ChapaWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Println("❌ Webhook JSON Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid webhook data"})
		return
	}

	fmt.Printf("🔔 Received Webhook for Ref: %s (Status: %s)\n", payload.TxRef, payload.Status)

	if payload.Status != "success" {
		c.JSON(http.StatusOK, gin.H{"message": "Payment not successful, no ticket issued"})
		return
	}

	parts := strings.Split(payload.TxRef, "-")
	if len(parts) < 5 {
		fmt.Println("❌ Error: Invalid tx_ref format received")
		c.JSON(http.StatusOK, gin.H{"message": "Invalid reference format"})
		return
	}

	eventID, errEv := strconv.Atoi(parts[2]) 
	userID, errUsr := strconv.Atoi(parts[4]) 

	if errEv != nil || errUsr != nil {
		fmt.Println("❌ Error converting IDs to integers")
		c.JSON(http.StatusOK, gin.H{"message": "ID conversion error"})
		return
	}

	tx, err := utils.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer tx.Rollback()

	var saleID int
	err = tx.QueryRow(`
		INSERT INTO ticket_sales (user_id, event_id, amount, status, tx_ref)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		userID, eventID, payload.Amount, "completed", payload.TxRef,
	).Scan(&saleID)

	if err != nil {
		fmt.Println("❌ Error saving to ticket_sales:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save sale"})
		return
	}

	_, err = tx.Exec(`
		INSERT INTO tickets (user_id, event_id, ticket_number, status)
		VALUES ($1, $2, $3, $4)`,
		userID, eventID, "TKT-"+payload.TxRef, "active",
	)

	if err != nil {
		fmt.Println("❌ Error saving to tickets:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to issue ticket"})
		return
	}

	if err = tx.Commit(); err != nil {
		fmt.Println("❌ Transaction commit failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "DB Commit Error"})
		return
	}

	fmt.Printf("✅ SUCCESS! Ticket issued to User %d for Event %d\n", userID, eventID)
	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}