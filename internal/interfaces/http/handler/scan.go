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
	OrganizationID string   `json:"organization_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Provider       string   `json:"provider" binding:"required,oneof=aws azure gcp" example:"aws"`
	Regions        []string `json:"regions" binding:"required,min=1" example:"us-east-1,eu-west-1"`
	ResourceTypes  []string `json:"resource_types" example:"ec2_instance,ebs_volume"`
}

// CreateScanResponse represents the response after creating a scan
type CreateScanResponse struct {
	Data    ScanDTO `json:"data"`
	Message string  `json:"message" example:"scan created and queued for processing"`
}

// Create godoc
//
//	@Summary		Create a new scan
//	@Description	Create a new cloud resource scan and queue it for processing
//	@Tags			Scans
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateScanRequest	true	"Scan request"
//	@Success		201		{object}	CreateScanResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/scans [post]
func (h *ScanHandler) Create(c *gin.Context) {
	var req CreateScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid organization ID"})
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create scan"})
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to enqueue scan task"})
		return
	}

	c.JSON(http.StatusCreated, CreateScanResponse{
		Data: ScanDTO{
			ID:             scan.ID.String(),
			OrganizationID: scan.OrganizationID.String(),
			Provider:       scan.Provider,
			Regions:        scan.Regions,
			ResourceTypes:  scan.ResourceTypes,
			Status:         scan.Status,
			CreatedAt:      scan.CreatedAt,
			UpdatedAt:      scan.UpdatedAt,
		},
		Message: "scan created and queued for processing",
	})
}

// ListScansRequest represents query parameters for listing scans
type ListScansRequest struct {
	Provider string `form:"provider" example:"aws"`
	Status   string `form:"status" example:"completed"`
	Limit    int    `form:"limit,default=20" example:"20"`
	Offset   int    `form:"offset,default=0" example:"0"`
}

// List godoc
//
//	@Summary		List scans
//	@Description	Get a paginated list of scans with optional filters
//	@Tags			Scans
//	@Accept			json
//	@Produce		json
//	@Param			provider	query		string	false	"Filter by cloud provider"	Enums(aws, azure, gcp)
//	@Param			status		query		string	false	"Filter by status"	Enums(pending, running, completed, failed, cancelled)
//	@Param			limit		query		int		false	"Number of items per page"	default(20)
//	@Param			offset		query		int		false	"Number of items to skip"	default(0)
//	@Success		200			{object}	PaginatedResponse{data=[]ScanDTO}
//	@Failure		400			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/scans [get]
func (h *ScanHandler) List(c *gin.Context) {
	var req ListScansRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch scans"})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:   scans,
		Total:  total,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
}

// Get godoc
//
//	@Summary		Get scan by ID
//	@Description	Get a single scan by its ID
//	@Tags			Scans
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Scan ID"	format(uuid)
//	@Success		200	{object}	map[string]ScanDTO
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/scans/{id} [get]
func (h *ScanHandler) Get(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid scan ID"})
		return
	}

	var scan model.Scan
	if err := h.db.First(&scan, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "scan not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch scan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": scan})
}
