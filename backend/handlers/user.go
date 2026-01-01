package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"local-event-backend/utils"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUserHandler(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Password error"})
		return
	}

	_, err = utils.DB.Exec(
		`INSERT INTO users (email, password) VALUES ($1, $2)`,
		req.Email,
		string(hashedPassword),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
