package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/adamn1225/affiliate-tracking/config"
	"github.com/adamn1225/affiliate-tracking/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database connection
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	defer config.CloseDB(db)

	// Run DB migrations
	config.MigrateDB(db)

	// Initialize Gin router
	r := gin.Default()

	// Inject DB into routes (if needed)
	routes.SetupRoutes(r, db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
