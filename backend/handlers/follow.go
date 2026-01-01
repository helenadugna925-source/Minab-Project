package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"local-event-backend/utils"
)

type FollowEventRequest struct {
	EventID int `json:"event_id"`
}

type FollowEventResponse struct {
	Status string    `json:"status"`
	ID     int       `json:"id"`
	UserID int       `json:"user_id"`
	EventID int      `json:"event_id"`
	CreatedAt time.Time `json:"created_at"`
}


func FollowEventHandler(c *gin.Context) {
	userID, ok := getUserIDFromAuthHeader(c)
	if !ok {
		return
	}

	var req FollowEventRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.EventID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	var id int
	err := utils.DB.QueryRow(`
		INSERT INTO event_followers (event_id, user_id, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT DO NOTHING
		RETURNING id
	`, req.EventID, userID).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, FollowEventResponse{
		Status: "followed",
		ID: id,
		UserID: userID,
		EventID: req.EventID,
		CreatedAt: time.Now(),
	})
}


func UnfollowEventHandler(c *gin.Context) {
	userID, ok := getUserIDFromAuthHeader(c)
	if !ok {
		return
	}

	var req FollowEventRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.EventID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	_, err := utils.DB.Exec(`DELETE FROM event_followers WHERE event_id=$1 AND user_id=$2`, req.EventID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "unfollowed"})
}

