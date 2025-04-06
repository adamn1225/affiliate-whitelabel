package handlers

import (
	"net/http"

	"strconv"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMyProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawID, exists := c.Get("user_id")
		if !exists {
  		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
  		return
	}

		rawRole, _ := c.Get("role")

		role := rawRole.(string)

if role == "vendor" {
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

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"email":        user.Email,
		"company_name": user.CompanyName,
		"role":         user.Role,
		"created_at":   user.CreatedAt,
	})
	return
}


		if role == "affiliate" {
			affiliateID := rawID.(string)
			var affiliate models.Affiliate
			if err := db.First(&affiliate, "id = ?", affiliateID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate not found"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"id":              affiliate.ID,
				"email":           affiliate.Email,
				"company_name":    affiliate.CompanyName,
				"contact_name":    affiliate.ContactName,
				"role":            affiliate.Role,
				"commission_rate": affiliate.CommissionRate,
				"created_at":      affiliate.CreatedAt,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown role"})
	}
}
