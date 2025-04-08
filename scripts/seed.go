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
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file loaded (ok in prod)")
    }

    // Database connection
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

    // Clear existing data
    db.Exec("DELETE FROM form_configs")
    db.Exec("DELETE FROM affiliates")
    db.Exec("DELETE FROM users")
    db.Exec("DELETE FROM affiliate_links")
    db.Exec("DELETE FROM leads")
    db.Exec("DELETE FROM affiliate_payouts")

    // Seed vendors (users)
    vendors := []*models.User{
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
        if err := db.Create(v).Error; err != nil {
            log.Fatal("Error creating vendor:", err)
        }
    }

    // Seed affiliates
    affiliates := []*models.Affiliate{
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
        if err := db.Create(a).Error; err != nil {
            log.Fatal("Error creating affiliate:", err)
        }
    }

    // Default form fields
    defaultFields := []models.FormField{
        {Label: "Name", Name: "name", Placeholder: "Full Name", Type: "text", Required: true},
        {Label: "Email", Name: "email", Placeholder: "Email Address", Type: "email", Required: true},
        {Label: "Phone", Name: "phone", Placeholder: "Phone Number", Type: "text", Required: false},
    }
    fieldJSON, _ := json.Marshal(defaultFields)

    // Define multiple form templates
    formTemplates := []struct {
        Title      string
        ButtonText string
        ButtonColor string
    }{
        {"Request a Transport Quote", "Submit", "#ff6600"},
        {"Schedule a Pickup", "Schedule", "#007bff"},
        {"Get a Delivery Estimate", "Get Estimate", "#28a745"},
    }

    // Iterate over the form templates and create forms
    for _, tpl := range formTemplates {
        form := models.FormConfig{
            UserID:      vendors[0].ID,
            Fields:      fieldJSON,
            FormTitle:   tpl.Title,
            ButtonText:  tpl.ButtonText,
            ButtonColor: tpl.ButtonColor,
        }
        if err := db.Create(&form).Error; err != nil {
            log.Fatalf("Error creating form config for vendor[0]: %v", err)
        }

        // Create affiliate links and leads for each form
        for _, affiliate := range affiliates {
            link := models.AffiliateLink{
                Slug:         uuid.New().String()[0:6],
                AffiliateID:  affiliate.ID,
                FormConfigID: form.ID,
                PayoutAmount: 25.00,
                CreatedAt:    time.Now(),
            }
            if err := db.Create(&link).Error; err != nil {
                log.Fatalf("Error creating affiliate link for affiliate %s: %v", affiliate.Email, err)
            }

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
                if err := db.Create(&lead).Error; err != nil {
                    log.Fatalf("Error creating lead for affiliate %s: %v", affiliate.Email, err)
                }

                payout := models.AffiliatePayout{
                    AffiliateID: affiliate.ID,
                    LeadID:      lead.ID,
                    Amount:      link.PayoutAmount,
                    Status:      "pending",
                    CreatedAt:   time.Now(),
                }
                if err := db.Create(&payout).Error; err != nil {
                    log.Fatalf("Error creating payout for affiliate %s: %v", affiliate.Email, err)
                }
            }
        }
    }

    log.Println("✅ Seeding complete.")
}

// Helper function to hash passwords
func hashPassword(pw string) string {
    hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
    return string(hash)
}