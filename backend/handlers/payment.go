package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
)

type HasuraActionPayload struct {
	Input struct {
		EventID  int    `json:"event_id"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	} `json:"input"`
}

func InitializePaymentHandler(c *gin.Context) {
	var payload HasuraActionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"checkout_url": "", "status": "failed", "message": "Invalid input"})
		return
	}

	// --- THE ULTIMATE EMAIL FIX ---
	email := strings.TrimSpace(payload.Input.Email)
	if email == "" || !strings.Contains(email, "@") {
		email = "customer.success@gmail.com" // Force valid email for Chapa
	}

	apiKey := os.Getenv("CHAPA_SECRET_KEY")
	txRef := fmt.Sprintf("tx-minab-1-%d-%d", payload.Input.EventID, time.Now().UnixNano())

	chapaBody := map[string]interface{}{
		"amount": "250", "currency": "ETB", "email": email, 
		"first_name": payload.Input.FullName, "tx_ref": txRef,
		"callback_url": "https://webhook.site/minab",
		"return_url":   "http://localhost:3000/events/verify",
		"customization": map[string]interface{}{ "title": "Minab Pass" },
	}

	jsonData, _ := json.Marshal(chapaBody)
	req, _ := http.NewRequest("POST", "https://api.chapa.co/v1/transaction/initialize", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		c.JSON(500, gin.H{"message": "Chapa connection failed"})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// Force status/message to be Strings
	statusStr := fmt.Sprintf("%v", result["status"])
	msgStr := fmt.Sprintf("%v", result["message"])

	if statusStr == "success" {
		data := result["data"].(map[string]interface{})
		c.JSON(200, gin.H{"checkout_url": fmt.Sprintf("%v", data["checkout_url"]), "status": "success", "message": "Success"})
	} else {
		c.JSON(400, gin.H{"checkout_url": "", "status": "failed", "message": msgStr})
	}
}