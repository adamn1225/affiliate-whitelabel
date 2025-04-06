package handlers

import (
	"net/http"

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
