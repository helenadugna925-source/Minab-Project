package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"local-event-backend/utils"
)

type CreateBookmarkRequest struct {
	EventID int `json:"event_id"`
}

type DeleteBookmarkRequest struct {
	EventID int `json:"event_id"`
}

type BookmarkResponse struct {
	ID        int       `json:"id"`
	EventID   int       `json:"event_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateBookmarkHandler(c *gin.Context) {
	userID, ok := getUserIDFromAuthHeader(c)
	if !ok {
		return
	}

	var req CreateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	var id int
	err := utils.DB.QueryRow(
		`INSERT INTO event_bookmarks (event_id, user_id, created_at)
		 VALUES ($1, $2, NOW())
		 RETURNING id`,
		req.EventID, userID,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, BookmarkResponse{
		ID:        id,
		EventID:   req.EventID,
		UserID:    userID,
		CreatedAt: time.Now(),
	})
}

func DeleteBookmarkHandler(c *gin.Context) {
	userID, ok := getUserIDFromAuthHeader(c)
	if !ok {
		return
	}

	var req DeleteBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	_, err := utils.DB.Exec(
		`DELETE FROM event_bookmarks WHERE event_id = $1 AND user_id = $2`,
		req.EventID, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bookmark removed",
	})
}

func getUserIDFromAuthHeader(c *gin.Context) (int, bool) {
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


func mustAtoi(s string) int {
	var i int
	for _, c := range s {
		i = i*10 + int(c-'0')
	}
	return i
}
