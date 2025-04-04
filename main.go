package main

import (
	"fmt"
	"log"
	"os"

	"github.com/adamn1225/affiliate-whitelabel/config"
	"github.com/adamn1225/affiliate-whitelabel/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
		}
	defer config.CloseDB(db)
		config.MigrateDB(db) 

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file loaded (ok in prod)")
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
	AllowOrigins:     []string{"*"}, // DEV ONLY — don't push to prod like this
	AllowMethods:     []string{"GET", "POST", "OPTIONS"},
	AllowHeaders:     []string{"Origin", "Content-Type"},
}))

	routes.SetupRoutes(router, db)

	// Add other routes here (POST /api/leads, etc.)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
