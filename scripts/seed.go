package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"gorm.io/datatypes"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=require"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	affiliates := []models.Affiliate{
		{
			ID:             "soflotrailerco",
			CompanyName:    "SoFlo Trailer Co, LLC",
			ContactName:    "Raphael Lantigua",
			Email:          "sofloco.ralphy@gmail.com",
			Phone:          "561-480-8596",
			Website:        "https://soflotrailerco.com/",
			CommissionRate: 0.1,
			CreatedAt:      time.Now(),
		},
		{
			ID:             "haulheroes",
			CompanyName:    "Haul Heroes Inc.",
			ContactName:    "Maggie Ryder",
			Email:          "maggie@haulheroes.com",
			Phone:          "555-555-0199",
			Website:        "https://haulheroes.com",
			CommissionRate: 0.12,
			CreatedAt:      time.Now(),
		},
	}

	previewConfig := models.FormConfig{
	AffiliateID: "preview-mode",
	FormTitle:   "Preview a Demo Form",
	ButtonText:  "Send Inquiry",
	ButtonColor: "#ff5733",
	Fields: datatypes.JSON(`[
		{"label": "Full Name", "name": "name", "placeholder": "John Doe", "type": "text", "required": true},
		{"label": "Email Address", "name": "email", "placeholder": "you@example.com", "type": "email", "required": true},
		{"label": "Phone Number", "name": "phone", "placeholder": "(123) 456-7890", "type": "text", "required": false}
	]`),
}

db.Clauses(clause.OnConflict{DoNothing: true}).Create(&previewConfig)

	

	for _, affiliate := range affiliates {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&affiliate)
	}
		// Default form fields to use for each affiliate
	defaultFields := []models.FormField{
		{Label: "Name", Name: "name", Placeholder: "Full Name", Type: "text", Required: true},
		{Label: "Email", Name: "email", Placeholder: "Email Address", Type: "email", Required: true},
		{Label: "Phone", Name: "phone", Placeholder: "Phone Number", Type: "tel", Required: false},
	}

	fieldJSON, err := json.Marshal(defaultFields)
	if err != nil {
		log.Fatal("Failed to marshal default fields:", err)
	}

	// One FormConfig per affiliate
	for _, affiliate := range affiliates {
		formConfig := models.FormConfig{
			AffiliateID: affiliate.ID,
			Fields:      fieldJSON,
			FormTitle:   "Request a Transport Quote",
			ButtonText:  "Submit Request",
			ButtonColor: "#000000",
		}

		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&formConfig)
	}

}
