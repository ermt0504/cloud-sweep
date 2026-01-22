package handler

import (
	"net/http"

	"github.com/cloudsweep/cloudsweep/internal/infrastructure/database/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PolicyHandler handles policy endpoints
type PolicyHandler struct {
	db *gorm.DB
}

// NewPolicyHandler creates a new PolicyHandler
func NewPolicyHandler(db *gorm.DB) *PolicyHandler {
	return &PolicyHandler{db: db}
}

// CreatePolicyRequest represents a request to create a new policy
type CreatePolicyRequest struct {
	OrganizationID string         `json:"organization_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name           string         `json:"name" binding:"required" example:"Delete unused EBS volumes"`
	Description    string         `json:"description" example:"Automatically delete EBS volumes unused for 30 days"`
	Provider       string         `json:"provider" binding:"required,oneof=aws azure gcp" example:"aws"`
	ResourceTypes  []string       `json:"resource_types" example:"ebs_volume,ebs_snapshot"`
	Conditions     map[string]any `json:"conditions"`
	Actions        []string       `json:"actions" binding:"required,min=1" example:"notify,delete"`
	Schedule       string         `json:"schedule" example:"0 0 * * *"`
}

// Create godoc
//
//	@Summary		Create policy
//	@Description	Create a new cleanup policy
//	@Tags			Policies
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreatePolicyRequest	true	"Policy request"
//	@Success		201		{object}	map[string]PolicyDTO
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/policies [post]
func (h *PolicyHandler) Create(c *gin.Context) {
	var req CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid organization ID"})
		return
	}

	policy := model.Policy{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           req.Name,
		Description:    req.Description,
		Provider:       req.Provider,
		ResourceTypes:  req.ResourceTypes,
		Conditions:     req.Conditions,
		Actions:        req.Actions,
		Schedule:       req.Schedule,
		IsEnabled:      true,
	}

	if err := h.db.Create(&policy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create policy"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": policy})
}

// ListPoliciesRequest represents query parameters for listing policies
type ListPoliciesRequest struct {
	Provider  string `form:"provider" example:"aws"`
	IsEnabled *bool  `form:"is_enabled" example:"true"`
	Limit     int    `form:"limit,default=20" example:"20"`
	Offset    int    `form:"offset,default=0" example:"0"`
}

// List godoc
//
//	@Summary		List policies
//	@Description	Get a paginated list of cleanup policies
//	@Tags			Policies
//	@Accept			json
//	@Produce		json
//	@Param			provider	query		string	false	"Filter by cloud provider"	Enums(aws, azure, gcp)
//	@Param			is_enabled	query		boolean	false	"Filter by enabled status"
//	@Param			limit		query		int		false	"Number of items per page"	default(20)
//	@Param			offset		query		int		false	"Number of items to skip"	default(0)
//	@Success		200			{object}	PaginatedResponse{data=[]PolicyDTO}
//	@Failure		400			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/policies [get]
func (h *PolicyHandler) List(c *gin.Context) {
	var req ListPoliciesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	query := h.db.Model(&model.Policy{})

	if req.Provider != "" {
		query = query.Where("provider = ?", req.Provider)
	}
	if req.IsEnabled != nil {
		query = query.Where("is_enabled = ?", *req.IsEnabled)
	}

	var total int64
	query.Count(&total)

	var policies []model.Policy
	if err := query.Limit(req.Limit).Offset(req.Offset).Order("created_at DESC").Find(&policies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch policies"})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:   policies,
		Total:  total,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
}

// Get godoc
//
//	@Summary		Get policy by ID
//	@Description	Get a single policy by its ID
//	@Tags			Policies
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Policy ID"	format(uuid)
//	@Success		200	{object}	map[string]PolicyDTO
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/policies/{id} [get]
func (h *PolicyHandler) Get(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid policy ID"})
		return
	}

	var policy model.Policy
	if err := h.db.First(&policy, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "policy not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch policy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": policy})
}

// Update godoc
//
//	@Summary		Update policy
//	@Description	Update an existing policy
//	@Tags			Policies
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Policy ID"	format(uuid)
//	@Param			request	body		CreatePolicyRequest	true	"Policy update request"
//	@Success		200		{object}	map[string]PolicyDTO
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/policies/{id} [put]
func (h *PolicyHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid policy ID"})
		return
	}

	var req CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updates := map[string]any{
		"name":           req.Name,
		"description":    req.Description,
		"provider":       req.Provider,
		"resource_types": req.ResourceTypes,
		"conditions":     req.Conditions,
		"actions":        req.Actions,
		"schedule":       req.Schedule,
	}

	result := h.db.Model(&model.Policy{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to update policy"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "policy not found"})
		return
	}

	var policy model.Policy
	h.db.First(&policy, "id = ?", id)

	c.JSON(http.StatusOK, gin.H{"data": policy})
}

// Delete godoc
//
//	@Summary		Delete policy
//	@Description	Delete a policy
//	@Tags			Policies
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Policy ID"	format(uuid)
//	@Success		200	{object}	MessageResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/policies/{id} [delete]
func (h *PolicyHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid policy ID"})
		return
	}

	result := h.db.Delete(&model.Policy{}, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to delete policy"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "policy not found"})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "policy deleted"})
}

// Enable godoc
//
//	@Summary		Enable policy
//	@Description	Enable a policy
//	@Tags			Policies
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Policy ID"	format(uuid)
//	@Success		200	{object}	MessageResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/policies/{id}/enable [post]
func (h *PolicyHandler) Enable(c *gin.Context) {
	h.setEnabled(c, true)
}

// Disable godoc
//
//	@Summary		Disable policy
//	@Description	Disable a policy
//	@Tags			Policies
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Policy ID"	format(uuid)
//	@Success		200	{object}	MessageResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/policies/{id}/disable [post]
func (h *PolicyHandler) Disable(c *gin.Context) {
	h.setEnabled(c, false)
}

func (h *PolicyHandler) setEnabled(c *gin.Context, enabled bool) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid policy ID"})
		return
	}

	result := h.db.Model(&model.Policy{}).Where("id = ?", id).Update("is_enabled", enabled)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to update policy"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "policy not found"})
		return
	}

	status := "disabled"
	if enabled {
		status = "enabled"
	}
	c.JSON(http.StatusOK, MessageResponse{Message: "policy " + status})
}
