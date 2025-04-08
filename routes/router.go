package routes

import (
	"github.com/adamn1225/affiliate-whitelabel/controllers"
	"github.com/adamn1225/affiliate-whitelabel/handlers" // Import the handlers package
	"github.com/adamn1225/affiliate-whitelabel/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	controllers.DB = db
	var r = router

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Affiliate API is alive"})
	})

	// Existing routes
	r.POST("/api/signup", handlers.Signup(db))
	r.POST("/api/vendor/signup", handlers.VendorSignup(db))
	r.POST("/api/affiliate/signup", handlers.AffiliateSignup(db))
	r.GET("/api/me", middlewares.RequireAuth(), handlers.GetMyProfile(db))
	r.GET("/api/my-leads", middlewares.RequireRole("vendor"), handlers.GetMyLeads(db))
	// Auth
	r.POST("/api/admin/login", controllers.AdminLogin)
	r.POST("/api/vendor/login", handlers.Login(db))
	r.POST("/api/affiliate/login", handlers.AffiliateLogin(db))
	r.POST("/api/logout", handlers.Logout)

	// FormConfig (consider renaming everything to match new structure later)
	r.GET("/api/formconfig/:affiliateId", controllers.GetFormConfig) // legacy for now
	r.POST("/api/formconfig", middlewares.RequireRole("vendor"), handlers.CreateForm(db))

	r.GET("/api/vendor/formconfigs", middlewares.RequireRole("vendor"), handlers.GetVendorFormConfigs(db))
	

	// Affiliate
	r.GET("/api/affiliate-form/:slug", handlers.GetAffiliateForm(db))
	r.GET("/api/affiliate/payouts", handlers.GetAffiliatePayouts(db))

	r.GET("/api/my-payouts", middlewares.RequireAuth(), handlers.GetMyPayouts(db))

	// Vendor
	r.GET("/api/vendor/leads", handlers.GetVendorLeads(db))

	// Lead submission
	r.GET("/api/leads/:form_id", middlewares.RequireRole("vendor"), handlers.GetLeadsByForm(db))
	r.POST("/api/leads", handlers.SubmitLead(db))
	r.GET("/api/leads", middlewares.RequireRole("vendor"), handlers.GetLeads(db))}
