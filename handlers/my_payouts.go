package handlers

import (
	"net/http"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMyPayouts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "affiliate" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		affiliateID := c.GetString("user_id")

		var payouts []models.AffiliatePayout
		if err := db.Where("affiliate_id = ?", affiliateID).Order("created_at desc").Find(&payouts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payouts"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"payouts": payouts})
	}
}
