package handlers

import (
	"log"
	"net/http"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllVendors(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var vendors []models.User
        if err := db.Where("role = ? AND public = ?", "vendor", true).Find(&vendors).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vendors"})
            return
        }

        var commissions []models.VendorCommission
        if err := db.Find(&commissions).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch commissions"})
            return
        }

        publicVendors := []models.PublicVendorListing{}
        for _, v := range vendors {
            commission := findCommission(commissions, v.ID)
            publicVendors = append(publicVendors, models.PublicVendorListing{
                ID:          v.ID,
                CompanyName: v.CompanyName,
                Industry:    v.Industry,
                Website:     v.Website,
                Description: v.Description,
                Commission:  commission,
                Public:      v.Public,
            })
        }

        c.JSON(http.StatusOK, gin.H{"vendors": publicVendors})
    }
}

func UpdateVendorVisibility(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input struct {
            Public bool `json:"public"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
            return
        }

        vendorID := c.GetUint("user_id") // Assuming middleware sets user_id in context
        if err := db.Model(&models.User{}).Where("id = ?", vendorID).Update("public", input.Public).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update visibility"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Visibility updated successfully"})
    }
}


func findCommission(commissions []models.VendorCommission, vendorID uint) float64 {
    for _, c := range commissions {
        if c.VendorID == vendorID && c.AffiliateID == nil {
            log.Printf("Found commission for VendorID %d: %f", vendorID, c.Commission)
            return c.Commission
        }
    }
    log.Printf("No commission found for VendorID %d, defaulting to 0.0", vendorID)
    return 0.0 // Default to 0% if no commission is found
}


func GetVendorByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		vendorID := c.Param("id")

		var vendor models.User
		if err := db.Where("id = ? AND role = ?", vendorID, "vendor").First(&vendor).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":           vendor.ID,
			"company_name": vendor.CompanyName,
			"industry":     vendor.Industry,
			"website":      vendor.Website,
			"address":      vendor.Address,
			"description":  vendor.Description, // optional field to add to your User model
		})
	}
}
