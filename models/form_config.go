package models

import (
	"gorm.io/datatypes"
)

type FormField struct {
	Label       string `json:"label"`
	Name        string `json:"name"`
	Placeholder string `json:"placeholder"`
	Type        string `json:"type"` // e.g. text, email, number, etc.
	Required    bool   `json:"required"`
}

type FormConfig struct {
	ID           uint           `gorm:"primaryKey"`
	AffiliateID  string         `gorm:"uniqueIndex"`
	Fields       datatypes.JSON // JSON array of FormField
	FormTitle    string
	ButtonText   string
	ButtonColor  string
}
