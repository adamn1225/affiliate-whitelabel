package handlers

import (
	"log"
	"net/http"
	"strings"

	"strconv"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/adamn1225/affiliate-whitelabel/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type VendorSignupInput struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	CompanyName string `json:"company_name"`
}

func VendorSignup(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input VendorSignupInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		input.Email = strings.ToLower(strings.TrimSpace(input.Email))

		// Check if email already exists
		var existing models.User
		if err := db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}

		// Hash password
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Create user
		user := models.User{
			Email:        input.Email,
			PasswordHash: string(hashed),
			CompanyName:  input.CompanyName,
			Role:         "vendor",
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		token, err := utils.GenerateJWT(strconv.FormatUint(uint64(user.ID), 10), user.Role)
		if err != nil {
			log.Println("Error generating token:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user":  gin.H{"email": user.Email, "company_name": user.CompanyName, "role": user.Role},
		})
	}
}
