package handlers

import (
	"net/http"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMyLeads(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "vendor" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		userID := uint(c.GetFloat64("user_id"))

		var leads []models.Lead
		if err := db.Where("user_id = ?", userID).Order("created_at desc").Find(&leads).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"leads": leads})
	}
}
