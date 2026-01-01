package handlers

import (
	"fmt"
	"local-event-backend/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)

type CreateEventRequest struct {
	Input struct {
		Title         string   `json:"title"`
		Description   string   `json:"description"`
		Date          string   `json:"date"`
		Price         float64  `json:"price"`
		LocationLat   float64  `json:"location_lat"`
		LocationLng   float64  `json:"location_lng"`
		VenueName     string   `json:"venue_name"`
		Address       string   `json:"address"`
		CategoryID    int      `json:"category_id"`
		Tags          []string `json:"tags"`
		ImageURLs     []string `json:"image_urls"`
		FeaturedImage string   `json:"featured_image"`
	} `json:"input"`
	SessionVariables map[string]interface{} `json:"session_variables"`
}

type EventResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

func CreateEventHandler(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ JSON Bind Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON payload"})
		return
	}

	var userIDStr string

	if req.SessionVariables != nil {
		if val, exists := req.SessionVariables["x-hasura-user-id"]; exists {
			userIDStr = fmt.Sprintf("%v", val)
		}
	}

	if userIDStr == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
			if err == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					// Check common Hasura claim paths
					if hasuraClaims, ok := claims["https://hasura.io/jwt/claims"].(map[string]interface{}); ok {
						userIDStr = fmt.Sprintf("%v", hasuraClaims["x-hasura-user-id"])
					} else if sub, ok := claims["sub"].(string); ok {
						userIDStr = sub
					}
				}
			}
		}
	}

	if userIDStr == "" || userIDStr == "<nil>" {
		fmt.Println("❌ Auth Error: No User ID found in SessionVariables or Token")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Please log in again. Your identity could not be verified."})
		return
	}

	userIDInt, _ := strconv.Atoi(userIDStr)
	fmt.Printf("✅ Authorized: Creating event for User ID %d\n", userIDInt)

	var eventID int
	query := `
		INSERT INTO events 
		(title, description, event_date, price, location_lat, location_lng, venue_name, address, category_id, user_id, featured_image, event_images)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`
	
	err := utils.DB.QueryRow(query, 
		req.Input.Title, 
		req.Input.Description, 
		req.Input.Date, 
		req.Input.Price,
		req.Input.LocationLat, 
		req.Input.LocationLng, 
		req.Input.VenueName, 
		req.Input.Address,
		req.Input.CategoryID, 
		userIDInt, 
		req.Input.FeaturedImage, 
		pq.Array(req.Input.ImageURLs),
	).Scan(&eventID)

	if err != nil {
		fmt.Println("❌ SQL Insert Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error: " + err.Error()})
		return
	}

	for _, tag := range req.Input.Tags {
		if tag == "" || tag == "General" { continue }
		_, err := utils.DB.Exec("INSERT INTO event_tags (event_id, tag_name) VALUES ($1, $2)", eventID, tag)
		if err != nil {
			fmt.Println("⚠️ Tag Insert Warning:", err)
		}
	}

	fmt.Printf("🎉 Success! Event %d published by User %d\n", eventID, userIDInt)
	
	
	c.JSON(http.StatusOK, EventResponse{
		ID:      eventID,
		Message: "Event created successfully",
	})
}