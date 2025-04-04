package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var db *gorm.DB
    

// POST /api/leads
type Lead struct {
	ID          uint           `gorm:"primaryKey"`
	AffiliateID string         `gorm:"index"`
	Data        datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt   int64          `gorm:"autoCreateTime"`
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


