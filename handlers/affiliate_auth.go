package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/adamn1225/affiliate-whitelabel/middlewares"
	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AffiliateAuthInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AffiliateLogin handles login for existing affiliate accounts
func AffiliateLogin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input AffiliateAuthInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		input.Email = strings.ToLower(strings.TrimSpace(input.Email))

		var affiliate models.Affiliate
		if err := db.Where("email = ?", input.Email).First(&affiliate).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(affiliate.PasswordHash), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		token, err := middlewares.GenerateToken(affiliate.ID, "affiliate")
		if err != nil {
			log.Println("JWT generation error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id":    affiliate.ID,
				"email": affiliate.Email,
				"role":  "affiliate",
			},
		})
	}
}
