package handlers

import (
	"net/http"
	"time"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TriggerPayouts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only admins
		role, _ := c.Get("role")
		if role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin access required"})
			return
		}

		var pending []models.AffiliatePayout
		db.Where("status = ?", "pending").Find(&pending)

		for _, payout := range pending {
			// Find vendor wallet by payout.LeadID -> leads.user_id
			var lead models.Lead
			if err := db.First(&lead, payout.LeadID).Error; err != nil {
				continue
			}
			var wallet models.VendorWallet
			if err := db.Where("vendor_id = ?", lead.UserID).First(&wallet).Error; err != nil {
				continue
			}

			if wallet.Balance >= payout.Amount {
				wallet.Balance -= payout.Amount
				payout.Status = "paid"
				now := time.Now()
				payout.PaidAt = &now

				db.Save(&wallet)
				db.Save(&payout)
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": "processed payouts"})
	}
}

func MarkPayoutPaid(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		payoutID := c.Param("id")
		var payout models.AffiliatePayout
		if err := db.First(&payout, payoutID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payout not found"})
			return
		}

		role, _ := c.Get("role")
		if role != "vendor" && role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		payout.Status = "manual"
		now := time.Now()
		payout.PaidAt = &now
		if err := db.Save(&payout).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payout"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Payout marked as manually paid"})
	}
}

