package handlers

import (
	"log"
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

func GetLeadsByForm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		formID := c.Param("form_id")

		var leads []models.Lead
		if err := db.Where("form_id = ?", formID).Find(&leads).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads"})
			return
		}

		c.JSON(http.StatusOK, leads)
	}
}

func GetMyLeads(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        role := c.GetString("role")
        if role != "vendor" {
            c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
            return
        }

        userID := uint(c.GetFloat64("user_id"))
        formID := c.Query("form_id")

        log.Printf("Fetching leads for user_id: %d, form_id: %s", userID, formID)

        var leads []models.Lead
        if err := db.Where("user_id = ? AND form_id = ?", userID, formID).
            Order("created_at desc").
            Find(&leads).Error; err != nil {
            log.Printf("Database query error: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads"})
            return
        }

        log.Printf("Fetched leads: %+v", leads)
        c.JSON(http.StatusOK, gin.H{"leads": leads})
    }
}

func GetRecentLeadsForZapier(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var leads []models.Lead
		if err := db.Order("created_at desc").Limit(10).Find(&leads).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads"})
			return
		}
		c.JSON(http.StatusOK, leads)
	}
}

