package handlers

import (
	"log"
	"net/http"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAffiliatePayouts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: pull from JWT/session in production
		affiliateIDRaw, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		affiliateID := affiliateIDRaw.(string)

		var payouts []models.AffiliatePayout
		if err := db.
			Where("affiliate_id = ?", affiliateID).
			Order("created_at desc").
			Find(&payouts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payouts"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"payouts": payouts,
		})
	}
}

func GetAffiliateWallet(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		affiliateID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		log.Printf("Affiliate ID from token: %v", affiliateID)

		var wallet models.AffiliateWallet
		if err := db.Where("affiliate_id = ?", affiliateID).First(&wallet).Error; err != nil {
			log.Println("Wallet not found for affiliate:", affiliateID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
			return
		}

		c.JSON(http.StatusOK, wallet)
	}
}

func GetAffiliateCommission(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		affiliateID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var commission models.AffiliateCommission
		if err := db.Where("affiliate_id = ?", affiliateID).First(&commission).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Commission not found"})
			return
		}

		c.JSON(http.StatusOK, commission)
	}
}
