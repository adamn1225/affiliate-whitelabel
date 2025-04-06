package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/adamn1225/affiliate-whitelabel/middlewares"
	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AffiliateSignupInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Company  string `json:"company_name"`
	Name     string `json:"contact_name"`
}

func AffiliateSignup(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input AffiliateSignupInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		input.Email = strings.ToLower(strings.TrimSpace(input.Email))

		// Check if affiliate already exists
		var existing models.Affiliate
		if err := db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Affiliate already exists"})
			return
		}

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
			return
		}

		// Create new affiliate
		affiliate := models.Affiliate{
			ID:           uuid.New().String(),
			Email:        input.Email,
			PasswordHash: string(hash),
			CompanyName:  input.Company,
			ContactName:  input.Name,
			Role:         "affiliate",
			CreatedAt:    time.Now(),
		}

		if err := db.Create(&affiliate).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create affiliate"})
			return
		}

		// Generate JWT
		token, err := middlewares.GenerateToken(affiliate.ID, affiliate.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id":    affiliate.ID,
				"email": affiliate.Email,
				"role":  affiliate.Role,
			},
		})
	}
}
