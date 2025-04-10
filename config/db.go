package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	log.Printf("Connecting to DB at %s:%s as user %s", host, port, user)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	log.Println("📡 Connecting to PostgreSQL...")
	var db *gorm.DB
	var err error

	for i := 1; i <= 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Connected to PostgreSQL")
			return db, nil
		}

		log.Printf("Attempt %d: %v", i, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect after 10 attempts: %w", err)
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Couldn't get sql.DB from GORM:", err)
		return
	}
	sqlDB.Close()
}

func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Affiliate{},
		&models.FormConfig{},
		&models.AffiliateLink{},
		&models.AffiliatePayout{},
		&models.Lead{},
		&models.VendorWallet{},
		&models.VendorCommission{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Database migrated")
}
