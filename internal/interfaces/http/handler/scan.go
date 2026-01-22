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

// ScanHandler handles scan endpoints
type ScanHandler struct {
	db          *gorm.DB
	queueClient *asynq.Client
}

// NewScanHandler creates a new ScanHandler
func NewScanHandler(db *gorm.DB, queueClient *asynq.Client) *ScanHandler {
	return &ScanHandler{
		db:          db,
		queueClient: queueClient,
	}
}

// CreateScanRequest represents a request to create a new scan
type CreateScanRequest struct {
	OrganizationID string   `json:"organization_id" binding:"required"`
	Provider       string   `json:"provider" binding:"required,oneof=aws azure gcp"`
	Regions        []string `json:"regions" binding:"required,min=1"`
	ResourceTypes  []string `json:"resource_types"`
}

// Create creates a new scan and enqueues it for processing
func (h *ScanHandler) Create(c *gin.Context) {
	var req CreateScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
		return
	}

	// Create scan record
	scan := model.Scan{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Provider:       req.Provider,
		Regions:        req.Regions,
		ResourceTypes:  req.ResourceTypes,
		Status:         "pending",
	}

	if err := h.db.Create(&scan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create scan"})
		return
	}

	// Enqueue scan task
	payload, _ := json.Marshal(queue.ScanResourcesPayload{
		OrganizationID: req.OrganizationID,
		Provider:       req.Provider,
		Regions:        req.Regions,
		ResourceTypes:  req.ResourceTypes,
	})

	task := asynq.NewTask(queue.TaskTypeScanResources, payload)
	if _, err := h.queueClient.Enqueue(task); err != nil {
		// Update scan status to failed
		h.db.Model(&scan).Update("status", "failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue scan task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    scan,
		"message": "scan created and queued for processing",
	})
}

// ListScansRequest represents query parameters for listing scans
type ListScansRequest struct {
	Provider string `form:"provider"`
	Status   string `form:"status"`
	Limit    int    `form:"limit,default=20"`
	Offset   int    `form:"offset,default=0"`
}

// List returns a list of scans
func (h *ScanHandler) List(c *gin.Context) {
	var req ListScansRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := h.db.Model(&model.Scan{})

	if req.Provider != "" {
		query = query.Where("provider = ?", req.Provider)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	var total int64
	query.Count(&total)

	var scans []model.Scan
	if err := query.Limit(req.Limit).Offset(req.Offset).Order("created_at DESC").Find(&scans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch scans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   scans,
		"total":  total,
		"limit":  req.Limit,
		"offset": req.Offset,
	})
}

// Get returns a single scan by ID
func (h *ScanHandler) Get(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scan ID"})
		return
	}

	var scan model.Scan
	if err := h.db.First(&scan, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "scan not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch scan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": scan})
}
