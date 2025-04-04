package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB
    

// POST /api/leads
type Lead struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string
	Email        string
	Phone        string
	Year         string
	Manufacturer string
	Model        string
	Message      string
	Origin       string `json:"origin"`
	Destination  string `json:"destination"`
	UtmSource    string `json:"utm_source"`
	UtmMedium    string `json:"utm_medium"`
	UtmCampaign  string `json:"utm_campaign"`
	AffiliateID  string `gorm:"index"`
	CreatedAt    int64  `gorm:"autoCreateTime"`
}

func CreateLead(c *gin.Context) {
    var lead Lead
    if err := c.ShouldBindJSON(&lead); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db.Create(&lead)
    c.JSON(http.StatusOK, gin.H{"status": "success"})
}


