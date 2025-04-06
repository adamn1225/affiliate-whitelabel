package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/adamn1225/affiliate-whitelabel/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		input.Email = strings.ToLower(strings.TrimSpace(input.Email))

		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Compare hash
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		token, err := utils.GenerateJWT(strconv.FormatUint(uint64(user.ID), 10), user.Role)
		if err != nil {
			log.Println("JWT generation error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		c.SetCookie("token", token, 86400, "/", "", false, false)


		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user":  gin.H{"email": user.Email, "role": user.Role, "company_name": user.CompanyName},
		})
	}
}
