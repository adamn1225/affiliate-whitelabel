package models

import (
	"time"
)

type Affiliate struct {
	ID             string    `gorm:"primaryKey" json:"id"`
	Role           string    `json:"role" gorm:"default:'affiliate'"`
	PasswordHash   string    `gorm:"not null" json:"-"`
	CompanyName    string    `json:"company_name"`
	ContactName    string    `json:"contact_name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Website        string    `json:"website"`
	CommissionRate float64   `json:"commission_rate"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AffiliateLink struct {
	ID           uint    `gorm:"primaryKey"`
	Slug         string  `gorm:"uniqueIndex;not null"` // e.g. abc123
	AffiliateID  string  `gorm:"not null"`             // User.Email or UUID
	FormConfigID uint    `gorm:"not null"`
	PayoutAmount float64 `gorm:"not null"`
	CreatedAt    time.Time
}

type AffiliatePayout struct {
	ID          uint    `gorm:"primaryKey"`
	AffiliateID string  `gorm:"not null"`
	LeadID      uint    `gorm:"not null"`
	Amount      float64 `gorm:"not null"`
	Status      string  `gorm:"default:'pending'"` // pending, paid
	PaidAt      *time.Time
	CreatedAt   time.Time
}
