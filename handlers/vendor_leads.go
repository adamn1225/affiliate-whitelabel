package handlers

import (
	"net/http"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetVendorLeads(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: replace with real vendor auth
		vendorIDRaw, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		vendorID := vendorIDRaw.(float64)

		var leads []models.Lead
		if err := db.
			Where("user_id = ?", vendorID).
			Order("created_at desc").
			Find(&leads).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"leads": leads})
	}
}
