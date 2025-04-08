package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func GetFormConfig(c *gin.Context) {
	affiliateID := c.Param("affiliateId")
	var config models.FormConfig

	if err := DB.Where("affiliate_id = ?", affiliateID).First(&config).Error; err != nil {
		// Fallback: return default config
		defaultFields := []models.FormField{
			{Label: "Name", Name: "name", Placeholder: "Enter your name", Type: "text", Required: true},
			{Label: "Phone", Name: "phone", Placeholder: "Enter your phone", Type: "text", Required: true},
			{Label: "Message", Name: "message", Placeholder: "Enter details", Type: "text", Required: false},
		}

		fieldsJSON, err := json.Marshal(defaultFields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode default fields"})
			return
		}

defaultConfig := models.FormConfig{
    AffiliateID:  &affiliateID, // Use & to convert string to *string
    Fields:       datatypes.JSON(fieldsJSON),
    FormTitle:    "Request a Transport Quote",
    ButtonText:   "Submit Request",
    ButtonColor:  "#000000",
}

		c.JSON(http.StatusOK, defaultConfig)
		return
	}

	c.JSON(http.StatusOK, config)
}

func CreateOrUpdateFormConfig(c *gin.Context) {
	var config models.FormConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

err := DB.Where(models.FormConfig{AffiliateID: config.AffiliateID}).Assign(config).FirstOrCreate(&config).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Save failed"})
		return
	}

	c.JSON(http.StatusOK, config)
}