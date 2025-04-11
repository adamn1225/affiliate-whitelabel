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
	r.GET("/api/public/vendors", handlers.GetAllVendors(db))
	r.GET("/api/public/vendors/:id", handlers.GetVendorByID(db))
	r.GET("/api/public/rotators/:slug", handlers.GetRotatorBySlug(db))
	r.GET("/r/:slug", handlers.RotateLink(db)) 
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
	r.GET("/api/zapier/leads", middlewares.RequireRole("admin"), handlers.GetRecentLeadsForZapier(db))

	// FormConfig (consider renaming everything to match new structure later)
	r.GET("/api/formconfig/:affiliateId", controllers.GetFormConfig) // legacy for now
	r.POST("/api/formconfig", middlewares.RequireRole("vendor"), handlers.CreateForm(db))

	r.GET("/api/vendor/formconfigs", middlewares.RequireRole("vendor"), handlers.GetVendorFormConfigs(db))
	//offer rotator
	r.PATCH("/api/affiliate/profile", middlewares.RequireRole("affiliate"), handlers.UpdateAffiliateProfile(db))
	r.GET("/api/affiliate/rotators/:id", middlewares.RequireRole("affiliate"), handlers.GetRotatorByID(db))
	r.GET("/api/affiliate/rotators", middlewares.RequireRole("affiliate"), handlers.GetMyRotators(db))
	r.POST("/api/affiliate/rotators", middlewares.RequireRole("affiliate"), handlers.CreateRotator(db))
	r.POST("/api/affiliate/rotators/:id/links", middlewares.RequireRole("affiliate"), handlers.AddRotatorLink(db))
	r.POST("/api/affiliate/rotators/auto", middlewares.RequireRole("affiliate"), handlers.CreateRotatorAuto(db))
	// Affiliate
	r.GET("/api/affiliate-form/:slug", handlers.GetAffiliateForm(db))
	r.GET("/api/affiliate/wallet", middlewares.RequireRole("affiliate"), handlers.GetAffiliateWallet(db))
	r.GET("/api/affiliate/commissions", middlewares.RequireRole("affiliate"), handlers.GetAffiliateCommission(db))
	r.GET("/api/my-payouts", middlewares.RequireAuth(), handlers.GetMyPayouts(db))
	
	r.GET("/api/affiliate/payouts", middlewares.RequireRole("affiliate"), handlers.GetAffiliatePayouts(db))
	
	r.GET("/api/affiliate/payout-test", handlers.GetAffiliatePayouts(db))
	// Admin
	r.POST("/api/admin/payouts/trigger", middlewares.RequireRole("admin"), handlers.TriggerPayouts(db))
	
	// Vendor
	// Lead submission
	r.GET("/api/vendor/commissions", middlewares.RequireRole("vendor"), handlers.GetVendorCommission(db))
	r.PATCH("/api/vendor/commission", middlewares.RequireRole("vendor"), handlers.UpdateVendorCommission(db))
	r.GET("/api/vendor/wallet", middlewares.RequireRole("vendor"), handlers.GetVendorWallet(db))
	r.GET("/api/vendor/leads", handlers.GetVendorLeads(db))
	r.GET("/api/leads/:form_id", middlewares.RequireRole("vendor"), handlers.GetLeadsByForm(db))
	r.POST("/api/leads", handlers.SubmitLead(db))
	r.GET("/api/leads", middlewares.RequireRole("vendor"), handlers.GetLeads(db))}

