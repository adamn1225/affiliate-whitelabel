package models

import "time"

type Affiliate struct {
	ID             string    `gorm:"primaryKey" json:"id"`
	Role 		   string    `json:"role" gorm:"default:'affiliate'"`
	CompanyName    string    `json:"company_name"`
	ContactName    string    `json:"contact_name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Website        string    `json:"website"`
	CommissionRate float64   `json:"commission_rate"`
	CreatedAt      time.Time `json:"created_at"`
}
