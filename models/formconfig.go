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
	ID          uint           `gorm:"primaryKey"`
	AffiliateID string         `gorm:"uniqueIndex" json:"affiliate_id"`
	Fields      datatypes.JSON `json:"fields"`
	FormTitle   string         `json:"form_title"`
	ButtonText  string         `json:"button_text"`
	ButtonColor string         `json:"button_color"`
	UserID      uint           `json:"user_id"`
}
