package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Seed vendors (users)
	vendors := []models.User{
		{
			Email:        "vendor1@example.com",
			PasswordHash: hashPassword("password123"),
			CompanyName:  "Vendor One LLC",
			Industry:     "Logistics",
			Address:      "123 Main St",
			Role:         "vendor",
			CreatedAt:    time.Now(),
		},
		{
			Email:        "vendor2@example.com",
			PasswordHash: hashPassword("securepass456"),
			CompanyName:  "Vendor Two Inc.",
			Industry:     "Transport",
			Address:      "456 Market Ave",
			Role:         "vendor",
			CreatedAt:    time.Now(),
		},
	}
	for _, v := range vendors {
		db.Create(&v)
	}

	// Seed affiliates
	affiliates := []models.Affiliate{
		{
			ID:             uuid.New().String(),
			Email:          "affiliate1@example.com",
			PasswordHash:   hashPassword("pass123"),
			CompanyName:    "Affiliate Co",
			ContactName:    "Alice",
			Phone:          "555-111-2222",
			Website:        "https://affiliateco.com",
			CommissionRate: 0.10,
			Role:           "affiliate",
			CreatedAt:      time.Now(),
		},
		{
			ID:             uuid.New().String(),
			Email:          "affiliate2@example.com",
			PasswordHash:   hashPassword("hunter2"),
			CompanyName:    "Promo Masters",
			ContactName:    "Bob",
			Phone:          "555-222-3333",
			Website:        "https://promomasters.com",
			CommissionRate: 0.12,
			Role:           "affiliate",
			CreatedAt:      time.Now(),
		},
	}
	for _, a := range affiliates {
		db.Create(&a)
	}

	// Default form fields
	defaultFields := []models.FormField{
		{Label: "Name", Name: "name", Placeholder: "Full Name", Type: "text", Required: true},
		{Label: "Email", Name: "email", Placeholder: "Email Address", Type: "email", Required: true},
		{Label: "Phone", Name: "phone", Placeholder: "Phone Number", Type: "text", Required: false},
	}
	fieldJSON, _ := json.Marshal(defaultFields)

	// Create 1 form for vendor[0]
	form := models.FormConfig{
		UserID:      vendors[0].ID,
		Fields:      fieldJSON,
		FormTitle:   "Request a Transport Quote",
		ButtonText:  "Submit",
		ButtonColor: "#ff6600",
	}
	db.Create(&form)

	// Create affiliate links for each affiliate
	for _, affiliate := range affiliates {
		link := models.AffiliateLink{
			Slug:         uuid.New().String()[0:6],
			AffiliateID:  affiliate.ID,
			FormConfigID: form.ID,
			PayoutAmount: 25.00,
			CreatedAt:    time.Now(),
		}
		db.Create(&link)

		// Create 2 fake leads per affiliate
		for i := 1; i <= 2; i++ {
			leadData := map[string]interface{}{
				"name":  affiliate.ContactName + " Lead " + string(rune(i)),
				"email": affiliate.Email,
				"phone": "123-456-7890",
			}
			dataJSON, _ := json.Marshal(leadData)
			lead := models.Lead{
				AffiliateID: affiliate.ID,
				FormID:      form.ID,
				UserID:      vendors[0].ID,
				Data:        dataJSON,
				Status:      "new",
				IsPaid:      false,
				CreatedAt:   time.Now(),
			}
			db.Create(&lead)

			// Add a payout record
			payout := models.AffiliatePayout{
				AffiliateID: affiliate.ID,
				LeadID:      lead.ID,
				Amount:      link.PayoutAmount,
				Status:      "pending",
				CreatedAt:   time.Now(),
			}
			db.Create(&payout)
		}
	}

	log.Println(" Seeding complete.")
}

func hashPassword(pw string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hash)
}
