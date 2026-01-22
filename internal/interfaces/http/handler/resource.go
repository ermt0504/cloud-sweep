package handler

import (
	"net/http"

	"github.com/cloudsweep/cloudsweep/internal/infrastructure/database/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// ResourceHandler handles resource endpoints
type ResourceHandler struct {
	db          *gorm.DB
	queueClient *asynq.Client
}

// NewResourceHandler creates a new ResourceHandler
func NewResourceHandler(db *gorm.DB, queueClient *asynq.Client) *ResourceHandler {
	return &ResourceHandler{
		db:          db,
		queueClient: queueClient,
	}
}

// ListResourcesRequest represents query parameters for listing resources
type ListResourcesRequest struct {
	Provider string `form:"provider"`
	Type     string `form:"type"`
	Status   string `form:"status"`
	Region   string `form:"region"`
	Limit    int    `form:"limit,default=50"`
	Offset   int    `form:"offset,default=0"`
}

// List returns a list of resources
func (h *ResourceHandler) List(c *gin.Context) {
	var req ListResourcesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build query
	query := h.db.Model(&model.Resource{})

	if req.Provider != "" {
		query = query.Where("provider = ?", req.Provider)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Region != "" {
		query = query.Where("region = ?", req.Region)
	}

	// Count total
	var total int64
	query.Count(&total)

	// Fetch resources
	var resources []model.Resource
	if err := query.Limit(req.Limit).Offset(req.Offset).Order("created_at DESC").Find(&resources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch resources"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   resources,
		"total":  total,
		"limit":  req.Limit,
		"offset": req.Offset,
	})
}

// Get returns a single resource by ID
func (h *ResourceHandler) Get(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource ID"})
		return
	}

	var resource model.Resource
	if err := h.db.First(&resource, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch resource"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resource})
}

// Delete deletes a resource (marks as deleted)
func (h *ResourceHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource ID"})
		return
	}

	result := h.db.Model(&model.Resource{}).Where("id = ?", id).Update("status", "deleted")
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete resource"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "resource deleted"})
}
