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

db.Exec("DELETE FROM form_configs")
db.Exec("DELETE FROM affiliates")
db.Exec("DELETE FROM users")
db.Exec("DELETE FROM affiliate_links")
db.Exec("DELETE FROM leads")
db.Exec("DELETE FROM affiliate_payouts")
db.Exec("DELETE FROM vendor_commissions")
db.Exec("DELETE FROM vendor_wallets")
db.Exec("SELECT setval(pg_get_serial_sequence('users', 'id'), 1, false)")
db.Exec("SELECT setval(pg_get_serial_sequence('vendor_wallets', 'id'), 1, false)")
db.Exec("SELECT setval(pg_get_serial_sequence('vendor_commissions', 'id'), 1, false)")

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

    wallet := models.VendorWallet{
        VendorID:  v.ID,
        Balance:   100.0,
        UpdatedAt: time.Now(),
    }
    if err := db.Create(&wallet).Error; err != nil {
        log.Fatal("Error creating vendor wallet:", err)
    }
}

if vendors[0].ID == 0 {
    log.Fatal("Vendor ID not set after creation. Check database connection or model configuration.")
}

// Create a VendorCommission for the first vendor
if err := db.Create(&models.VendorCommission{
    VendorID:   vendors[0].ID, 
    Commission: 0.15,   
    CreatedAt:  time.Now(),
}).Error; err != nil {
    log.Fatal("Error creating vendor commission:", err)
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

        wallet := models.AffiliateWallet{
    AffiliateID: a.ID,
    Balance:     42.00,
    UpdatedAt:   time.Now(),
}
if err := db.Create(&wallet).Error; err != nil {
    log.Fatalf("Error creating affiliate wallet for %s: %v", a.Email, err)
}

commission := models.AffiliateCommission{
    AffiliateID: a.ID,
    Commission:  a.CommissionRate,
    CreatedAt:   time.Now(),
}
if err := db.Create(&commission).Error; err != nil {
    log.Fatalf("Error creating affiliate commission for %s: %v", a.Email, err)
}
    }

     for _, affiliate := range affiliates {
    rotator := models.OfferRotator{
        AffiliateID: affiliate.ID,
        Name:        "My First Rotator",
        Slug:        uuid.New().String()[0:8],
        Strategy:    "random", // or "weighted" later
        CreatedAt:   time.Now(),
    }

    if err := db.Create(&rotator).Error; err != nil {
        log.Fatalf("Error creating rotator for affiliate %s: %v", affiliate.Email, err)
    }

    // Create wallet for affiliate affiliate1@example.com
if err := db.Create(&models.AffiliateWallet{
    AffiliateID: "7cc60aa5-d01c-4e87-949b-7ee9936bd8e7",
    Balance:     50.0,
    UpdatedAt:   time.Now(),
}).Error; err != nil {
    log.Fatalf("Error creating wallet for affiliate %s: %v", "affiliate1@example.com", err)
}

// Create commission record for affiliate affiliate1@example.com
if err := db.Create(&models.AffiliateCommission{
    AffiliateID: "7cc60aa5-d01c-4e87-949b-7ee9936bd8e7",
    Commission:  0.10,
    CreatedAt:   time.Now(),
}).Error; err != nil {
    log.Fatalf("Error creating commission for affiliate %s: %v", "affiliate1@example.com", err)
}

// Create wallet for affiliate affiliate2@example.com
if err := db.Create(&models.AffiliateWallet{
    AffiliateID: "c90fb43e-95ac-4a01-b42b-d684832ac38c",
    Balance:     50.0,
    UpdatedAt:   time.Now(),
}).Error; err != nil {
    log.Fatalf("Error creating wallet for affiliate %s: %v", "affiliate2@example.com", err)
}

// Create commission record for affiliate affiliate2@example.com
if err := db.Create(&models.AffiliateCommission{
    AffiliateID: "c90fb43e-95ac-4a01-b42b-d684832ac38c",
    Commission:  0.12,
    CreatedAt:   time.Now(),
}).Error; err != nil {
    log.Fatalf("Error creating commission for affiliate %s: %v", "affiliate2@example.com", err)
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

        // Seed rotators for each affiliate
   


    // Add sample links
    links := []models.RotatorLink{
        {
            RotatorID: rotator.ID,
            URL:       "https://example.com/offer1",
            Weight:    1,
            CreatedAt: time.Now(),
        },
        {
            RotatorID: rotator.ID,
            URL:       "https://example.com/offer2",
            Weight:    2,
            CreatedAt: time.Now(),
        },
    }

    for _, link := range links {
        if err := db.Create(&link).Error; err != nil {
            log.Fatalf("Error creating link for rotator %d: %v", rotator.ID, err)
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