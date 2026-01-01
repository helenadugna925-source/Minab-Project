package handlers

import (
	"bytes"
	"fmt"
	"io"
	"local-event-backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
		RememberMe  bool   `json:"remember_me"`
	} `json:"input"`
}

type SignupResponse struct {
	Token   string `json:"token"`
	UserID  int    `json:"user_id"`
	Message string `json:"message"` 
}

func SignupHandler(c *gin.Context) {
	var req SignupRequest


	raw, _ := io.ReadAll(c.Request.Body)
	fmt.Println("🔍 HASURA SENT:", string(raw))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON format"})
		return
	}

	email := strings.TrimSpace(req.Input.Email)
	password := req.Input.Password
	fullName := fmt.Sprintf("%s %s", req.Input.FirstName, req.Input.LastName)
	phone := req.Input.PhoneNumber

	
	if email == "" || password == "" {
		fmt.Printf("❌ DATA MISSING: Email='%s', Password='%s'\n", email, password)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email and password are required"})
		return
	}

	
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Encryption failed"})
		return
	}

	var userID int
	query := `INSERT INTO users (email, password, name, phone_number) VALUES ($1, $2, $3, $4) RETURNING id`
	err = utils.DB.QueryRow(query, email, string(hashed), fullName, phone).Scan(&userID)

	if err != nil {
		fmt.Println("❌ DB ERROR:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "A user with this email already exists"})
		return
	}

	
	token, err := utils.GenerateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Token generation failed"})
		return
	}

	
	fmt.Printf("✅ Success! Created User %d\n", userID)
	c.JSON(http.StatusOK, SignupResponse{
		Token:   token,
		UserID:  userID,
		Message: "User registered successfully", // <--- Matches 'message' in Hasura
	})
}