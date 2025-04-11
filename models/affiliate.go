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
	Slug         string  `gorm:"uniqueIndex;not null"`
	AffiliateID  string  `gorm:"not null"`
	FormConfigID uint    `gorm:"not null"`
	PayoutAmount float64 `gorm:"not null"`
	CreatedAt    time.Time
}

type AffiliatePayout struct {
	ID           uint           `gorm:"primaryKey"`
	AffiliateID  string         `gorm:"index"`
	LeadID       uint           `gorm:"index"`
	Amount       float64        
	AffiliateCut float64 		`json:"affiliate_cut"`
	PlatformFee  float64 		`json:"platform_fee"`
	PaidAt       *time.Time 	`json:"paid_at,omitempty"`
	Status       string         `gorm:"default:'pending'"`
	CreatedAt    time.Time
}

type AffiliateWallet struct {
	ID         uint      `gorm:"primaryKey"`
	AffiliateID string    `gorm:"index;not null"`
	Balance    float64   `gorm:"default:0"`
	UpdatedAt  time.Time
}

type AffiliateCommission struct {
	ID          uint      `gorm:"primaryKey"`
	AffiliateID string    `gorm:"index;not null"`
	Commission  float64   `gorm:"default:0.1"` // default 10%
	CreatedAt   time.Time
}