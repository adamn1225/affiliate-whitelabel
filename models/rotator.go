package models

import (
	"time"

	"gorm.io/datatypes"
)

type Rotator struct {
	ID          uint   `gorm:"primaryKey"`
	AffiliateID string `gorm:"index"`
	Name        string
	CreatedAt   time.Time
}

type OfferRotator struct {
	ID          uint      `gorm:"primaryKey"`
	AffiliateID string    `gorm:"index"`
	Name        string    `gorm:"not null"`
	Slug        string    `gorm:"uniqueIndex"`
	Strategy    string    `gorm:"default:'random'"` 
	CreatedAt   time.Time
}


type RotatorLink struct {
	ID         uint           `gorm:"primaryKey"`
	RotatorID  uint           `gorm:"index"`
	URL        string         `gorm:"not null"`
	Weight     int            `gorm:"default:1"`
	Clicks     int            `gorm:"default:0"`
	Conditions datatypes.JSON `gorm:"type:jsonb"` // << NEW
	CreatedAt  time.Time
}

type RotatorClick struct {
	ID         uint      `gorm:"primaryKey"`
	RotatorID  uint      `gorm:"index"`
	LinkID     uint      `gorm:"index"`
	IPAddress  string
	UserAgent  string
	CreatedAt  time.Time
}