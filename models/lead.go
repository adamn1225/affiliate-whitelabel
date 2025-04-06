package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var db *gorm.DB

// POST /api/leads
type Lead struct {
	ID          uint           `gorm:"primaryKey"`
	AffiliateID string         `gorm:"index"`
	FormID      uint           `gorm:"index"` // new
	UserID      uint           `gorm:"index"` // vendor ID
	Data        datatypes.JSON `gorm:"type:jsonb"`
	Status      string         `gorm:"default:'new'"` // new, verified, rejected
	IsPaid      bool           `gorm:"default:false"`
	CreatedAt   time.Time
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

func ToJSON(data interface{}) datatypes.JSON {
	bytes, _ := json.Marshal(data)
	return datatypes.JSON(bytes)
}
