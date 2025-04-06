package handlers

import (
	"net/http"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAffiliateForm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		if slug == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing slug"})
			return
		}

		var link models.AffiliateLink
		if err := db.Where("slug = ?", slug).First(&link).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Affiliate link not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving affiliate link"})
			}
			return
		}

		var form models.FormConfig
		if err := db.First(&form, link.FormConfigID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Form not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"form": gin.H{
				"id":           form.ID,
				"title":        form.FormTitle,
				"fields":       form.Fields,
				"button_text":  form.ButtonText,
				"button_color": form.ButtonColor,
			},
			"affiliate_id":  link.AffiliateID,
			"payout_amount": link.PayoutAmount,
		})
	}
}
