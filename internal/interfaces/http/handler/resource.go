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
	Provider string `form:"provider" example:"aws"`
	Type     string `form:"type" example:"ec2_instance"`
	Status   string `form:"status" example:"unused"`
	Region   string `form:"region" example:"us-east-1"`
	Limit    int    `form:"limit,default=50" example:"50"`
	Offset   int    `form:"offset,default=0" example:"0"`
}

// List godoc
//
//	@Summary		List resources
//	@Description	Get a paginated list of cloud resources with optional filters
//	@Tags			Resources
//	@Accept			json
//	@Produce		json
//	@Param			provider	query		string	false	"Filter by cloud provider"	Enums(aws, azure, gcp)
//	@Param			type		query		string	false	"Filter by resource type"
//	@Param			status		query		string	false	"Filter by status"	Enums(active, unused, deleted, excluded)
//	@Param			region		query		string	false	"Filter by region"
//	@Param			limit		query		int		false	"Number of items per page"	default(50)
//	@Param			offset		query		int		false	"Number of items to skip"	default(0)
//	@Success		200			{object}	PaginatedResponse{data=[]ResourceDTO}
//	@Failure		400			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/resources [get]
func (h *ResourceHandler) List(c *gin.Context) {
	var req ListResourcesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch resources"})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:   resources,
		Total:  total,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
}

// Get godoc
//
//	@Summary		Get resource by ID
//	@Description	Get a single cloud resource by its ID
//	@Tags			Resources
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Resource ID"	format(uuid)
//	@Success		200	{object}	map[string]ResourceDTO
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/resources/{id} [get]
func (h *ResourceHandler) Get(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid resource ID"})
		return
	}

	var resource model.Resource
	if err := h.db.First(&resource, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "resource not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch resource"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resource})
}

// Delete godoc
//
//	@Summary		Delete resource
//	@Description	Mark a resource as deleted
//	@Tags			Resources
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Resource ID"	format(uuid)
//	@Success		200	{object}	MessageResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/resources/{id} [delete]
func (h *ResourceHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid resource ID"})
		return
	}

	result := h.db.Model(&model.Resource{}).Where("id = ?", id).Update("status", "deleted")
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to delete resource"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "resource not found"})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "resource deleted"})
}
