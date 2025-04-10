package handlers

import (
	"net/http"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SignupInput struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	CompanyName string `json:"company_name"`
	Industry    string `json:"industry"`
	Address     string `json:"address"`
	Role        string `json:"role"` // default to vendor if blank
}

type UpdateUserInput struct {
    Password    *string `json:"password,omitempty"`
    Address     *string `json:"address,omitempty"`
    Website     *string `json:"website,omitempty"`
    Phone       *string `json:"phone,omitempty"`
    FirstName   *string `json:"first_name,omitempty"`
    LastName    *string `json:"last_name,omitempty"`
    CompanyName *string `json:"company_name,omitempty"`
    Industry    *string `json:"industry,omitempty"`
}

func UpdateUser(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get user ID from the token (set by middleware)
        userID, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }

        var input UpdateUserInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
            return
        }

        var user models.User
        if err := db.First(&user, userID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        // Update fields if provided
        if input.Password != nil {
            hash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
                return
            }
            user.PasswordHash = string(hash)
        }
        if input.Address != nil {
            user.Address = *input.Address
        }
        if input.Website != nil {
            user.Website = *input.Website
        }
        if input.Phone != nil {
            user.Phone = *input.Phone
        }
        if input.FirstName != nil {
            user.FirstName = *input.FirstName
        }
        if input.LastName != nil {
            user.LastName = *input.LastName
        }
        if input.CompanyName != nil {
            user.CompanyName = *input.CompanyName
        }
        if input.Industry != nil {
            user.Industry = *input.Industry
        }

        // Save updates
        if err := db.Save(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
    }
}

func Signup(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input SignupInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user := models.User{
			Email:        input.Email,
			PasswordHash: string(hash),
			CompanyName:  input.CompanyName,
			Industry:     input.Industry,
			Address:      input.Address,
			Role:         input.Role,
		}

		if user.Role == "" {
			user.Role = "vendor"
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}
