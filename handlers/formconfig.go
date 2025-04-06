package handlers

import (
	"net/http"

	"strconv"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateForm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.FormConfig

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}

		// Grab user ID from JWT claims (set by middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		input.UserID = userID.(uint)

		if err := db.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save form"})
			return
		}

		c.JSON(http.StatusOK, input)
	}
}

func GetVendorFormConfigs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userIDStr, ok := rawID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			return
		}

		userIDUint64, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
			return
		}

		userID := uint(userIDUint64)

		var forms []models.FormConfig
		if err := db.Where("user_id = ?", userID).Find(&forms).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch forms"})
			return
		}

		c.JSON(http.StatusOK, forms)
	}
}

