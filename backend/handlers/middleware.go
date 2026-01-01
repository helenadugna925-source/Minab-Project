package handlers

import (
    "fmt"
    "net/http"
    "strings"
    "local-event-backend/utils"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }

        tokenStr := parts[1]
        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
			return utils.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if hasuraClaims, ok := claims["https://hasura.io/jwt/claims"].(map[string]interface{}); ok {
				if userID, ok := hasuraClaims["x-hasura-user-id"].(string); ok {
					c.Set("userID", userID)
					c.Next()
					return
				}
			}
			
			if userID, ok := claims["user_id"].(string); ok {
				c.Set("userID", userID)
				c.Next()
				return
			}
		}
	}
}