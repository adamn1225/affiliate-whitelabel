package routes

import (
	"github.com/adamn1225/affiliate-whitelabel/controllers"
	"github.com/adamn1225/affiliate-whitelabel/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
    controllers.DB = db

    router.Use(cors.Default())

    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Affiliate API is alive"})
    })

    router.POST("/api/leads", controllers.CreateLead)
    router.GET("/api/affiliates/:id/leads", controllers.GetLeadsByAffiliate)
    router.POST("/api/affiliates", controllers.CreateAffiliate)
	router.GET("/api/leads", middlewares.RequireAdmin(), controllers.GetAllLeads)
	router.GET("/api/affiliates", middlewares.RequireAdmin(), controllers.GetAllAffiliates)
	router.POST("/api/admin/login", controllers.AdminLogin)
	router.GET("/api/form-config/:id", controllers.GetFormConfig)
	router.POST("/api/form-config", controllers.CreateOrUpdateFormConfig)
}
