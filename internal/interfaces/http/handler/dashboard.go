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

// Summary returns dashboard summary statistics
func (h *DashboardHandler) Summary(c *gin.Context) {
	var stats struct {
		TotalResources  int64   `json:"total_resources"`
		UnusedResources int64   `json:"unused_resources"`
		TotalCost       float64 `json:"total_monthly_cost"`
		PotentialSavings float64 `json:"potential_monthly_savings"`
		TotalCarbon     float64 `json:"total_carbon_kg"`
		CarbonSavings   float64 `json:"potential_carbon_savings_kg"`
	}

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

// Savings returns savings breakdown by provider and resource type
func (h *DashboardHandler) Savings(c *gin.Context) {
	// By provider
	var byProvider []struct {
		Provider string  `json:"provider"`
		Cost     float64 `json:"monthly_cost"`
		Savings  float64 `json:"potential_savings"`
		Count    int64   `json:"unused_count"`
	}

	h.db.Model(&model.Resource{}).
		Select("provider, SUM(monthly_cost) as cost, COUNT(*) as count").
		Where("status = ?", "unused").
		Group("provider").
		Scan(&byProvider)

	// By resource type
	var byType []struct {
		Type    string  `json:"resource_type"`
		Cost    float64 `json:"monthly_cost"`
		Count   int64   `json:"unused_count"`
	}

	h.db.Model(&model.Resource{}).
		Select("type, SUM(monthly_cost) as cost, COUNT(*) as count").
		Where("status = ?", "unused").
		Group("type").
		Order("cost DESC").
		Limit(10).
		Scan(&byType)

	c.JSON(http.StatusOK, gin.H{
		"by_provider":      byProvider,
		"by_resource_type": byType,
	})
}

// Carbon returns carbon footprint breakdown
func (h *DashboardHandler) Carbon(c *gin.Context) {
	// By provider
	var byProvider []struct {
		Provider string  `json:"provider"`
		Carbon   float64 `json:"carbon_kg"`
		Savings  float64 `json:"potential_savings_kg"`
	}

	h.db.Model(&model.Resource{}).
		Select("provider, SUM(carbon_footprint) as carbon").
		Where("status = ?", "unused").
		Group("provider").
		Scan(&byProvider)

	// By region
	var byRegion []struct {
		Region string  `json:"region"`
		Carbon float64 `json:"carbon_kg"`
	}

	h.db.Model(&model.Resource{}).
		Select("region, SUM(carbon_footprint) as carbon").
		Where("status = ?", "unused").
		Group("region").
		Order("carbon DESC").
		Limit(10).
		Scan(&byRegion)

	c.JSON(http.StatusOK, gin.H{
		"by_provider": byProvider,
		"by_region":   byRegion,
	})
}
