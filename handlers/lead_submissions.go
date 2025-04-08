package handlers

import (
	"net/http"
	"time"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LeadSubmission struct {
	FormID      uint                   `json:"form_id" binding:"required"`
	AffiliateID string                 `json:"affiliate_id"`
	Data        map[string]any `json:"data" binding:"required"`
}

func SubmitLead(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LeadSubmission
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lead submission"})
			return
		}

		var form models.FormConfig
		if err := db.First(&form, input.FormID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Form not found"})
			return
		}

		lead := models.Lead{
			FormID:      form.ID,
			UserID:      form.UserID,
			AffiliateID: input.AffiliateID,
			Data:        models.ToJSON(input.Data),
			Status:      "new",
			IsPaid:      false,
			CreatedAt:   time.Now(),
		}

		if err := db.Create(&lead).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create lead"})
			return
		}

		if input.AffiliateID != "" {
			var link models.AffiliateLink
			if err := db.
				Where("form_config_id = ? AND affiliate_id = ?", form.ID, input.AffiliateID).
				First(&link).Error; err == nil {

				payout := models.AffiliatePayout{
					AffiliateID: input.AffiliateID,
					LeadID:      lead.ID,
					Amount:      link.PayoutAmount,
					Status:      "pending",
					CreatedAt:   time.Now(),
				}
				db.Create(&payout)
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Lead submitted successfully"})
	}
}

func GetLeads(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        formID := c.Query("form_id")
        if formID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "form_id is required"})
            return
        }

        var leads []models.Lead
        if err := db.Where("form_id = ?", formID).Order("created_at desc").Find(&leads).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads"})
            return
        }

        // Return leads as raw array
        c.JSON(http.StatusOK, leads)
    }
}