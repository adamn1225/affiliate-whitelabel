package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	CompanyName  string
	Industry     string
	Address      string
	Role         string `gorm:"default:vendor"`
	CreatedAt    time.Time
}
