package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cloudsweep/cloudsweep/internal/infrastructure/database/model"
	"github.com/cloudsweep/cloudsweep/internal/infrastructure/queue"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// CleanupHandler handles cleanup endpoints
type CleanupHandler struct {
	db          *gorm.DB
	queueClient *asynq.Client
}

// NewCleanupHandler creates a new CleanupHandler
func NewCleanupHandler(db *gorm.DB, queueClient *asynq.Client) *CleanupHandler {
	return &CleanupHandler{
		db:          db,
		queueClient: queueClient,
	}
}

// ExecuteCleanupRequest represents a request to execute cleanup
type ExecuteCleanupRequest struct {
	OrganizationID string   `json:"organization_id" binding:"required"`
	ResourceIDs    []string `json:"resource_ids" binding:"required,min=1"`
	Action         string   `json:"action" binding:"required,oneof=delete stop tag notify"`
	DryRun         bool     `json:"dry_run"`
}

// Execute executes a cleanup operation
func (h *CleanupHandler) Execute(c *gin.Context) {
	var req ExecuteCleanupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate resource IDs
	for _, id := range req.ResourceIDs {
		if _, err := uuid.Parse(id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource ID: " + id})
			return
		}
	}

	// Enqueue cleanup task
	payload, _ := json.Marshal(queue.CleanupResourcesPayload{
		OrganizationID: req.OrganizationID,
		ResourceIDs:    req.ResourceIDs,
		Action:         req.Action,
		DryRun:         req.DryRun,
	})

	task := asynq.NewTask(queue.TaskTypeCleanupResources, payload)
	info, err := h.queueClient.Enqueue(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue cleanup task"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "cleanup task queued",
		"task_id": info.ID,
		"dry_run": req.DryRun,
	})
}

// Preview returns a preview of what would be cleaned up
func (h *CleanupHandler) Preview(c *gin.Context) {
	var req ExecuteCleanupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert string IDs to UUIDs
	var uuids []uuid.UUID
	for _, id := range req.ResourceIDs {
		u, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource ID: " + id})
			return
		}
		uuids = append(uuids, u)
	}

	// Fetch resources
	var resources []model.Resource
	if err := h.db.Where("id IN ?", uuids).Find(&resources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch resources"})
		return
	}

	// Calculate totals
	var totalCost, totalCarbon float64
	for _, r := range resources {
		totalCost += r.MonthlyCost
		totalCarbon += r.CarbonFootprint
	}

	c.JSON(http.StatusOK, gin.H{
		"resources":              resources,
		"count":                  len(resources),
		"estimated_monthly_savings": totalCost,
		"estimated_carbon_savings":  totalCarbon,
		"action":                 req.Action,
	})
}
