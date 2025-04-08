package controllers

import (
	"encoding/json"
	"net/http"

	"strings"

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
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
        return
    }

    if config.AffiliateID == nil || *config.AffiliateID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Affiliate ID is required"})
        return
    }

    err := DB.Where(models.FormConfig{AffiliateID: config.AffiliateID}).Assign(config).FirstOrCreate(&config).Error
    if err != nil {
        if strings.Contains(err.Error(), "duplicate key value") {
            c.JSON(http.StatusConflict, gin.H{"error": "Affiliate ID already exists"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Save failed"})
        return
    }

    c.JSON(http.StatusOK, config)
}

func GetVendorFormConfigs(c *gin.Context) {
    var forms []models.FormConfig

    // Fetch forms from the database
    if err := DB.Find(&forms).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch forms"})
        return
    }

    // Unmarshal the fields for each form
    for i := range forms {
        var fields []models.FormField
        if err := json.Unmarshal(forms[i].Fields, &fields); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse form fields"})
            return
        }
        marshaledFields, err := json.Marshal(fields)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to re-marshal form fields"})
            return
        }
        forms[i].Fields = datatypes.JSON(marshaledFields) // Replace the raw JSON with the re-marshaled fields
    }

    c.JSON(http.StatusOK, forms)
}