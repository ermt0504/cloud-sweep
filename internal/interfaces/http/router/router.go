package router

import (
	"github.com/cloudsweep/cloudsweep/internal/infrastructure/config"
	"github.com/cloudsweep/cloudsweep/internal/interfaces/http/handler"
	"github.com/cloudsweep/cloudsweep/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	_ "github.com/cloudsweep/cloudsweep/docs" // Swagger docs
)

// NewRouter creates and configures the Gin router
func NewRouter(db *gorm.DB, queueClient *asynq.Client, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	// Health check
	healthHandler := handler.NewHealthHandler(db)
	r.GET("/health", healthHandler.Check)
	r.GET("/ready", healthHandler.Ready)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Resources
		resourceHandler := handler.NewResourceHandler(db, queueClient)
		resources := v1.Group("/resources")
		{
			resources.GET("", resourceHandler.List)
			resources.GET("/:id", resourceHandler.Get)
			resources.DELETE("/:id", resourceHandler.Delete)
		}

		// Scans
		scanHandler := handler.NewScanHandler(db, queueClient)
		scans := v1.Group("/scans")
		{
			scans.POST("", scanHandler.Create)
			scans.GET("", scanHandler.List)
			scans.GET("/:id", scanHandler.Get)
		}

		// Cleanup
		cleanupHandler := handler.NewCleanupHandler(db, queueClient)
		v1.POST("/cleanup", cleanupHandler.Execute)
		v1.POST("/cleanup/preview", cleanupHandler.Preview)

		// Policies
		policyHandler := handler.NewPolicyHandler(db)
		policies := v1.Group("/policies")
		{
			policies.POST("", policyHandler.Create)
			policies.GET("", policyHandler.List)
			policies.GET("/:id", policyHandler.Get)
			policies.PUT("/:id", policyHandler.Update)
			policies.DELETE("/:id", policyHandler.Delete)
			policies.POST("/:id/enable", policyHandler.Enable)
			policies.POST("/:id/disable", policyHandler.Disable)
		}

		// Dashboard / Stats
		dashboardHandler := handler.NewDashboardHandler(db)
		v1.GET("/dashboard/summary", dashboardHandler.Summary)
		v1.GET("/dashboard/savings", dashboardHandler.Savings)
		v1.GET("/dashboard/carbon", dashboardHandler.Carbon)
	}

	return r
}
