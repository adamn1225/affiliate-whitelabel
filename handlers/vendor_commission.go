package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetVendorWallet(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userIDRaw, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }

        var userID uint
        switch v := userIDRaw.(type) {
        case string:
            id, err := strconv.ParseUint(v, 10, 32)
            if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
                return
            }
            userID = uint(id)
        case uint:
            userID = v
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected user ID type"})
            return
        }

        var wallet models.VendorWallet
        if err := db.Where("vendor_id = ?", userID).First(&wallet).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                // Return a default wallet or a more user-friendly error
                c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found. Please contact support."})
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
            }
            return
        }

        c.JSON(http.StatusOK, wallet)
    }
}


func UpdateVendorCommission(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Retrieve vendorID from context
        vendorIDRaw, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }

        // Safely convert vendorID to uint
        var vendorID uint
        switch v := vendorIDRaw.(type) {
        case string:
            id, err := strconv.ParseUint(v, 10, 32)
            if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendor ID"})
                return
            }
            vendorID = uint(id)
        case uint:
            vendorID = v
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected vendor ID type"})
            return
        }

        // Parse input JSON
        var input struct {
            Commission float64 `json:"commission"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
            return
        }

        // Only one global commission per vendor for now
        var vc models.VendorCommission
        err := db.Where("vendor_id = ? AND affiliate_id IS NULL", vendorID).First(&vc).Error
        if err != nil && err != gorm.ErrRecordNotFound {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
            return
        }

        vc.VendorID = vendorID
        vc.Commission = input.Commission
        vc.CreatedAt = time.Now()

        if err == gorm.ErrRecordNotFound {
            if err := db.Create(&vc).Error; err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create commission setting"})
                return
            }
        } else {
            if err := db.Save(&vc).Error; err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update commission"})
                return
            }
        }

        c.JSON(http.StatusOK, gin.H{"message": "Commission updated"})
    }
}

func GetVendorCommission(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        vendorIDRaw, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }

        var vendorID uint
        switch v := vendorIDRaw.(type) {
        case string:
            id, err := strconv.ParseUint(v, 10, 32)
            if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendor ID"})
                return
            }
            vendorID = uint(id)
        case uint:
            vendorID = v
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected vendor ID type"})
            return
        }

        log.Printf("Fetching commission for vendor ID: %d", vendorID)

        var vc models.VendorCommission
        if err := db.Where("vendor_id = ? AND affiliate_id IS NULL", vendorID).First(&vc).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                log.Printf("No commission found for vendor ID: %d", vendorID)
                c.JSON(http.StatusNotFound, gin.H{"error": "Commission not found"})
            } else {
                log.Printf("Database error: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
            }
            return
        }

        log.Printf("Commission found: %+v", vc)
        c.JSON(http.StatusOK, vc)
    }
}