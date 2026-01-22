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
	OrganizationID string              `json:"organization_id" binding:"required"`
	Name           string              `json:"name" binding:"required"`
	Description    string              `json:"description"`
	Provider       string              `json:"provider" binding:"required,oneof=aws azure gcp"`
	ResourceTypes  []string            `json:"resource_types"`
	Conditions     map[string]any      `json:"conditions"`
	Actions        []string            `json:"actions" binding:"required,min=1"`
	Schedule       string              `json:"schedule"`
}

// Create creates a new policy
func (h *PolicyHandler) Create(c *gin.Context) {
	var req CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create policy"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": policy})
}

// ListPoliciesRequest represents query parameters for listing policies
type ListPoliciesRequest struct {
	Provider  string `form:"provider"`
	IsEnabled *bool  `form:"is_enabled"`
	Limit     int    `form:"limit,default=20"`
	Offset    int    `form:"offset,default=0"`
}

// List returns a list of policies
func (h *PolicyHandler) List(c *gin.Context) {
	var req ListPoliciesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch policies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   policies,
		"total":  total,
		"limit":  req.Limit,
		"offset": req.Offset,
	})
}

// Get returns a single policy by ID
func (h *PolicyHandler) Get(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy ID"})
		return
	}

	var policy model.Policy
	if err := h.db.First(&policy, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch policy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": policy})
}

// Update updates a policy
func (h *PolicyHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy ID"})
		return
	}

	var req CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update policy"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
		return
	}

	var policy model.Policy
	h.db.First(&policy, "id = ?", id)

	c.JSON(http.StatusOK, gin.H{"data": policy})
}

// Delete deletes a policy
func (h *PolicyHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy ID"})
		return
	}

	result := h.db.Delete(&model.Policy{}, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete policy"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "policy deleted"})
}

// Enable enables a policy
func (h *PolicyHandler) Enable(c *gin.Context) {
	h.setEnabled(c, true)
}

// Disable disables a policy
func (h *PolicyHandler) Disable(c *gin.Context) {
	h.setEnabled(c, false)
}

func (h *PolicyHandler) setEnabled(c *gin.Context, enabled bool) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy ID"})
		return
	}

	result := h.db.Model(&model.Policy{}).Where("id = ?", id).Update("is_enabled", enabled)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update policy"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
		return
	}

	status := "disabled"
	if enabled {
		status = "enabled"
	}
	c.JSON(http.StatusOK, gin.H{"message": "policy " + status})
}
