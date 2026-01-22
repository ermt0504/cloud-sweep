package handler

import (
	"net/http"

	"github.com/cloudsweep/cloudsweep/internal/infrastructure/database/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DashboardHandler handles dashboard endpoints
type DashboardHandler struct {
	db *gorm.DB
}

// NewDashboardHandler creates a new DashboardHandler
func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// SummaryStats represents dashboard summary statistics
type SummaryStats struct {
	TotalResources   int64   `json:"total_resources" example:"500"`
	UnusedResources  int64   `json:"unused_resources" example:"75"`
	TotalCost        float64 `json:"total_monthly_cost" example:"15000.00"`
	PotentialSavings float64 `json:"potential_monthly_savings" example:"2500.00"`
	TotalCarbon      float64 `json:"total_carbon_kg" example:"1200.50"`
	CarbonSavings    float64 `json:"potential_carbon_savings_kg" example:"180.25"`
}

// ProviderSavings represents savings by provider
type ProviderSavings struct {
	Provider string  `json:"provider" example:"aws"`
	Cost     float64 `json:"monthly_cost" example:"1500.00"`
	Savings  float64 `json:"potential_savings" example:"250.00"`
	Count    int64   `json:"unused_count" example:"25"`
}

// TypeSavings represents savings by resource type
type TypeSavings struct {
	Type  string  `json:"resource_type" example:"ec2_instance"`
	Cost  float64 `json:"monthly_cost" example:"800.00"`
	Count int64   `json:"unused_count" example:"10"`
}

// SavingsResponse represents savings breakdown response
type SavingsResponse struct {
	ByProvider     []ProviderSavings `json:"by_provider"`
	ByResourceType []TypeSavings     `json:"by_resource_type"`
}

// ProviderCarbon represents carbon by provider
type ProviderCarbon struct {
	Provider string  `json:"provider" example:"aws"`
	Carbon   float64 `json:"carbon_kg" example:"450.25"`
	Savings  float64 `json:"potential_savings_kg" example:"75.50"`
}

// RegionCarbon represents carbon by region
type RegionCarbon struct {
	Region string  `json:"region" example:"us-east-1"`
	Carbon float64 `json:"carbon_kg" example:"250.00"`
}

// CarbonResponse represents carbon breakdown response
type CarbonResponse struct {
	ByProvider []ProviderCarbon `json:"by_provider"`
	ByRegion   []RegionCarbon   `json:"by_region"`
}

// Summary godoc
//
//	@Summary		Dashboard summary
//	@Description	Get dashboard summary statistics including total resources, unused resources, costs and carbon footprint
//	@Tags			Dashboard
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]SummaryStats
//	@Failure		500	{object}	ErrorResponse
//	@Router			/dashboard/summary [get]
func (h *DashboardHandler) Summary(c *gin.Context) {
	var stats SummaryStats

	// Total resources
	h.db.Model(&model.Resource{}).Where("status != ?", "deleted").Count(&stats.TotalResources)

	// Unused resources
	h.db.Model(&model.Resource{}).Where("status = ?", "unused").Count(&stats.UnusedResources)

	// Total cost
	h.db.Model(&model.Resource{}).
		Where("status != ?", "deleted").
		Select("COALESCE(SUM(monthly_cost), 0)").
		Scan(&stats.TotalCost)

	// Potential savings (unused resources cost)
	h.db.Model(&model.Resource{}).
		Where("status = ?", "unused").
		Select("COALESCE(SUM(monthly_cost), 0)").
		Scan(&stats.PotentialSavings)

	// Total carbon
	h.db.Model(&model.Resource{}).
		Where("status != ?", "deleted").
		Select("COALESCE(SUM(carbon_footprint), 0)").
		Scan(&stats.TotalCarbon)

	// Carbon savings
	h.db.Model(&model.Resource{}).
		Where("status = ?", "unused").
		Select("COALESCE(SUM(carbon_footprint), 0)").
		Scan(&stats.CarbonSavings)

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// Savings godoc
//
//	@Summary		Savings breakdown
//	@Description	Get potential savings breakdown by provider and resource type
//	@Tags			Dashboard
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	SavingsResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/dashboard/savings [get]
func (h *DashboardHandler) Savings(c *gin.Context) {
	// By provider
	var byProvider []ProviderSavings

	h.db.Model(&model.Resource{}).
		Select("provider, SUM(monthly_cost) as cost, COUNT(*) as count").
		Where("status = ?", "unused").
		Group("provider").
		Scan(&byProvider)

	// By resource type
	var byType []TypeSavings

	h.db.Model(&model.Resource{}).
		Select("type, SUM(monthly_cost) as cost, COUNT(*) as count").
		Where("status = ?", "unused").
		Group("type").
		Order("cost DESC").
		Limit(10).
		Scan(&byType)

	c.JSON(http.StatusOK, SavingsResponse{
		ByProvider:     byProvider,
		ByResourceType: byType,
	})
}

// Carbon godoc
//
//	@Summary		Carbon footprint breakdown
//	@Description	Get carbon footprint breakdown by provider and region
//	@Tags			Dashboard
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	CarbonResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/dashboard/carbon [get]
func (h *DashboardHandler) Carbon(c *gin.Context) {
	// By provider
	var byProvider []ProviderCarbon

	h.db.Model(&model.Resource{}).
		Select("provider, SUM(carbon_footprint) as carbon").
		Where("status = ?", "unused").
		Group("provider").
		Scan(&byProvider)

	// By region
	var byRegion []RegionCarbon

	h.db.Model(&model.Resource{}).
		Select("region, SUM(carbon_footprint) as carbon").
		Where("status = ?", "unused").
		Group("region").
		Order("carbon DESC").
		Limit(10).
		Scan(&byRegion)

	c.JSON(http.StatusOK, CarbonResponse{
		ByProvider: byProvider,
		ByRegion:   byRegion,
	})
}
