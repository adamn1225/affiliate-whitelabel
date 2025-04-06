package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("Wa8wDrAuSHy31FvN7MgZQUUe8tTt5yqmsUraBa7E0XA=")
func init() {
	log.Println("JWT Secret:", os.Getenv("JWT_SECRET"))
}

func GenerateJWT(userID string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 1 week
	})

	return token.SignedString(jwtSecret)
}
