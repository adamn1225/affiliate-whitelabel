package controllers

import (
	"net/http"
	"time"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB // this will be set from routes.SetupRoutes()

func CreateLead(c *gin.Context) {
	var lead models.Lead
	if err := c.ShouldBindJSON(&lead); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
    lead.CreatedAt = time.Now().Unix()
	if err := DB.Create(&lead).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "created", "lead": lead})
}


func GetAllLeads(c *gin.Context) {
    var leads []models.Lead

    if err := DB.Find(&leads).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads"})
        return
    }

    c.JSON(http.StatusOK, leads)
}
