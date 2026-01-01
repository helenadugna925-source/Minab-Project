package handlers

import (
	"database/sql"
	"fmt"
	"local-event-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Input struct {
		Email      string `json:"email"`       
		Password   string `json:"password"`    
		RememberMe bool   `json:"remember_me"` 
	} `json:"input"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	UserID  int    `json:"user_id"`
	Message string `json:"message"` 
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Login JSON Bind Error:", err)
		c.JSON(http.StatusBadRequest, LoginResponse{
			Message: "Invalid request format",
		})
		return
	}

	email := req.Input.Email
	password := req.Input.Password

	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, LoginResponse{
			Message: "Email and password are required",
		})
		return
	}

	var storedHash string
	var userID int

	err := utils.DB.QueryRow("SELECT id, password FROM users WHERE email = $1", email).Scan(&userID, &storedHash)

	if err == sql.ErrNoRows {
		fmt.Printf("❌ Login Failed: User %s not found\n", email)
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Message: "Invalid email or password",
		})
		return
	} else if err != nil {
		fmt.Println("❌ Database Error during login:", err)
		c.JSON(http.StatusInternalServerError, LoginResponse{
			Message: "Internal server error",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		fmt.Printf("❌ Login Failed: Incorrect password for %s\n", email)
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Message: "Invalid email or password",
		})
		return
	}

	token, err := utils.GenerateToken(userID)
	if err != nil {
		fmt.Println("❌ Token Generation Failed:", err)
		c.JSON(http.StatusInternalServerError, LoginResponse{
			Message: "Failed to generate security token",
		})
		return
	}

	fmt.Printf("✅ User %s logged in successfully (ID: %d)\n", email, userID)
	c.JSON(http.StatusOK, LoginResponse{
		Token:   token,
		UserID:  userID,
		Message: "Login successful",
	})
}