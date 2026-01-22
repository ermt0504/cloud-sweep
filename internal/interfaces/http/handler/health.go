package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Service string `json:"service" example:"cloudsweep-api"`
}

// ReadyResponse represents a readiness check response
type ReadyResponse struct {
	Status string            `json:"status" example:"ready"`
	Checks map[string]string `json:"checks"`
}

// Check godoc
//
//	@Summary		Health check
//	@Description	Basic health check endpoint
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	HealthResponse
//	@Router			/health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "ok",
		Service: "cloudsweep-api",
	})
}

// Ready godoc
//
//	@Summary		Readiness check
//	@Description	Readiness check with dependency verification
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ReadyResponse
//	@Failure		503	{object}	ErrorResponse
//	@Router			/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error: "database connection unavailable",
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error: "database ping failed",
		})
		return
	}

	c.JSON(http.StatusOK, ReadyResponse{
		Status: "ready",
		Checks: map[string]string{
			"database": "ok",
		},
	})
}
