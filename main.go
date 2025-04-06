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
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file loaded (ok in prod)")
	}
	db, err := config.ConnectDB()

	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	defer config.CloseDB(db)
	config.MigrateDB(db)

	// Load environment variables from .env file

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // or "*" in dev
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.Use(func(c *gin.Context) {
	log.Printf("[CORS DEBUG] Method: %s, Path: %s, Origin: %s", c.Request.Method, c.Request.URL.Path, c.GetHeader("Origin"))
	c.Next()
})

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
