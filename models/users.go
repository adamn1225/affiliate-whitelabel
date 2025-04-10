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
    Website      string
    Phone        string
    FirstName    string
    LastName     string
	Role         string `gorm:"default:vendor"`
	CreatedAt    time.Time
}

type VendorWallet struct {
    ID        uint    `gorm:"primaryKey"`
    VendorID  uint    `gorm:"uniqueIndex"`
    Balance   float64
    BillingAddress string
    BillingCity    string
    BillingState   string
    BillingZip     string
    BillingCountry  string
    UpdatedAt time.Time
}

type VendorCommission struct {
    ID             uint    `gorm:"primaryKey"`
    VendorID       uint    `gorm:"index"`
    AffiliateID    *string `gorm:"index;default:null"`
    Commission     float64
    CreatedAt      time.Time
}

func GetCommissionRate(db *gorm.DB, vendorID uint, affiliateID string) float64 {
    var vc VendorCommission

    // Try affiliate-specific override first
    err := db.Where("vendor_id = ? AND affiliate_id = ?", vendorID, affiliateID).First(&vc).Error
    if err == nil {
        return vc.Commission
    }

    // Fallback to default
    err = db.Where("vendor_id = ? AND affiliate_id IS NULL", vendorID).First(&vc).Error
    if err == nil {
        return vc.Commission
    }

    // Fallback to platform default (e.g. 10%)
    return 0.10
}
