package utils

import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("supersecretkeylocaldev1234567891")

// GenerateToken returns a Hasura-compatible JWT signed with JwtSecret.
func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": fmt.Sprintf("%d", userID),
		"iat": time.Now().Add(-time.Minute * 5).Unix(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"https://hasura.io/jwt/claims": map[string]interface{}{
			"x-hasura-allowed-roles": []string{"user"},
			"x-hasura-default-role":  "user",
			"x-hasura-user-id":       fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}