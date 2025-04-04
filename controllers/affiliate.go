package controllers

import (
	"net/http"

	"github.com/adamn1225/affiliate-tracking/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


func CreateAffiliate(c *gin.Context) {
	var affiliate models.Affiliate
	if err := c.ShouldBindJSON(&affiliate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if affiliate.ID == "" {
		affiliate.ID = uuid.NewString()
	}

	if err := DB.Create(&affiliate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "created",
		"affiliate": affiliate,
	})
}

func GetLeadsByAffiliate(c *gin.Context) {
	affiliateID := c.Param("id")
	var leads []models.Lead

	if err := DB.Where("affiliate_id = ?", affiliateID).Order("created_at desc").Find(&leads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads"})
		return
	}

	c.JSON(http.StatusOK, leads)
}

func GetAllAffiliates(c *gin.Context) {
	var affiliates []models.Affiliate

	if err := DB.Find(&affiliates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch affiliates"})
		return
	}

	c.JSON(http.StatusOK, affiliates)
}

