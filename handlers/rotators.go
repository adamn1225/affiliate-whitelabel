package handlers

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"encoding/json"

	"github.com/adamn1225/affiliate-whitelabel/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func parseUint(s string) uint {
	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(id)
}

func RotateLink(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")

		var rotator models.OfferRotator
		if err := db.Where("slug = ?", slug).First(&rotator).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rotator not found"})
			return
		}

		var links []models.RotatorLink
		if err := db.Where("rotator_id = ?", rotator.ID).Find(&links).Error; err != nil || len(links) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No links available for this rotator"})
			return
		}

		// Pick link first
		var selected models.RotatorLink
		switch rotator.Strategy {
		case "random":
			rand.Seed(time.Now().UnixNano())
			selected = links[rand.Intn(len(links))]
		default:
			selected = links[0]
		}

		// Increment click count
		db.Model(&selected).UpdateColumn("clicks", gorm.Expr("clicks + ?", 1))

		// Log the click event
		click := models.RotatorClick{
			RotatorID: rotator.ID,
			LinkID:    selected.ID,
			IPAddress: c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
			CreatedAt: time.Now(),
		}
		db.Create(&click)

		var recentClick models.RotatorClick
	if err := db.
    	Where("rotator_id = ? AND ip_address = ?", rotator.ID, c.ClientIP()).
    	Order("created_at desc").
    	First(&recentClick).Error; err == nil {
    
    // Already clicked recently — skip incrementing or recording
    if time.Since(recentClick.CreatedAt) < 10*time.Second {
        log.Printf("Rate-limited click for rotator %s from IP %s", slug, c.ClientIP())
        c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limited"})
        return
    }
}

		log.Printf("Redirecting rotator slug %s to %s", slug, selected.URL)
		c.Redirect(http.StatusFound, selected.URL)
	}
}


func GetRotatorClicks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var clicks []models.RotatorClick
		if err := db.Where("rotator_id = ?", id).Order("created_at desc").Find(&clicks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch clicks"})
			return
		}
		c.JSON(http.StatusOK, clicks)
	}
}

func CreateRotator(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name     string `json:"name"`
			Strategy string `json:"strategy"` // e.g. "random" or "sequential"
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		affiliateIDRaw, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		rotator := models.OfferRotator{
			AffiliateID: affiliateIDRaw.(string),
			Name:        input.Name,
			Slug:        uuid.NewString()[0:8],
			Strategy:    input.Strategy,
			CreatedAt:   time.Now(),
		}

		if err := db.Create(&rotator).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rotator"})
			return
		}

		c.JSON(http.StatusOK, rotator)
	}
}



func GetMyRotators(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		affiliateID, _ := c.Get("user_id")

		var rotators []models.OfferRotator
		if err := db.Where("affiliate_id = ?", affiliateID).Find(&rotators).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rotators"})
			return
		}

		c.JSON(http.StatusOK, rotators)
	}
}

func AddRotatorLink(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rotatorID := c.Param("rotator_id")

		var input struct {
			URL    string `json:"url"`
			Weight int    `json:"weight"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		link := models.RotatorLink{
			RotatorID: parseUint(rotatorID),
			URL:       input.URL,
			Weight:    input.Weight,
			CreatedAt: time.Now(),
		}

		if err := db.Create(&link).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add link"})
			return
		}

		c.JSON(http.StatusOK, link)

	}
}

func GetRotatorByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var rotator models.OfferRotator
		if err := db.First(&rotator, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rotator not found"})
			return
		}

		var links []models.RotatorLink
		if err := db.Where("rotator_id = ?", rotator.ID).Find(&links).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"rotator": rotator,
			"links":   links,
		})
	}
}

func GetRotatorBySlug(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		var rotator models.OfferRotator
		if err := db.Where("slug = ?", slug).First(&rotator).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rotator not found"})
			return
		}

		var links []models.RotatorLink
		if err := db.Where("rotator_id = ?", rotator.ID).Find(&links).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"rotator": rotator,
			"links":   links,
		})
	}
}

func CreateRotatorAuto(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name     string `json:"name"`
			Strategy string `json:"strategy"` // e.g., "rules"
			Links    []struct {
				URL        string                 `json:"url"`
				Weight     int                    `json:"weight"`
				Conditions map[string]interface{} `json:"conditions"` // optional
			} `json:"links"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		affiliateIDRaw, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		affiliateID := affiliateIDRaw.(string)

		// Create rotator
		rotator := models.OfferRotator{
			AffiliateID: affiliateID,
			Name:        input.Name,
			Slug:        uuid.NewString()[:8],
			Strategy:    input.Strategy,
			CreatedAt:   time.Now(),
		}

		if err := db.Create(&rotator).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rotator"})
			return
		}

		// Create links with optional conditions
		for _, l := range input.Links {
			condJSON, err := json.Marshal(l.Conditions)
			if err != nil {
				log.Printf("Error marshaling conditions: %v", err)
				continue // skip bad condition data
			}

			link := models.RotatorLink{
				RotatorID:  rotator.ID,
				URL:        l.URL,
				Weight:     l.Weight,
				Conditions: condJSON,
				CreatedAt:  time.Now(),
			}

			if err := db.Create(&link).Error; err != nil {
				log.Printf("Error creating link: %v", err)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"rotator": rotator,
			"message": "Rotator created with rule-based links",
		})
	}
}


