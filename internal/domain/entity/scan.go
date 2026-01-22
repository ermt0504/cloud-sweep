package entity

import (
	"time"

	"github.com/google/uuid"
)

// ScanStatus represents the status of a scan
type ScanStatus string

const (
	ScanStatusPending    ScanStatus = "pending"
	ScanStatusRunning    ScanStatus = "running"
	ScanStatusCompleted  ScanStatus = "completed"
	ScanStatusFailed     ScanStatus = "failed"
	ScanStatusCancelled  ScanStatus = "cancelled"
)

// Scan represents a cloud resource scan
type Scan struct {
	ID               uuid.UUID       `json:"id"`
	OrganizationID   uuid.UUID       `json:"organization_id"`
	Provider         CloudProvider   `json:"provider"`
	Regions          []string        `json:"regions"`
	ResourceTypes    []ResourceType  `json:"resource_types"`
	Status           ScanStatus      `json:"status"`
	ResourcesFound   int             `json:"resources_found"`
	UnusedFound      int             `json:"unused_found"`
	EstimatedSavings float64         `json:"estimated_savings"`
	CarbonSavings    float64         `json:"carbon_savings_kg"`
	ErrorMessage     string          `json:"error_message,omitempty"`
	StartedAt        *time.Time      `json:"started_at,omitempty"`
	CompletedAt      *time.Time      `json:"completed_at,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// NewScan creates a new Scan
func NewScan(orgID uuid.UUID, provider CloudProvider, regions []string, resourceTypes []ResourceType) *Scan {
	now := time.Now()
	return &Scan{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Provider:       provider,
		Regions:        regions,
		ResourceTypes:  resourceTypes,
		Status:         ScanStatusPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Start marks the scan as running
func (s *Scan) Start() {
	now := time.Now()
	s.Status = ScanStatusRunning
	s.StartedAt = &now
	s.UpdatedAt = now
}

// Complete marks the scan as completed
func (s *Scan) Complete(resourcesFound, unusedFound int, estimatedSavings, carbonSavings float64) {
	now := time.Now()
	s.Status = ScanStatusCompleted
	s.ResourcesFound = resourcesFound
	s.UnusedFound = unusedFound
	s.EstimatedSavings = estimatedSavings
	s.CarbonSavings = carbonSavings
	s.CompletedAt = &now
	s.UpdatedAt = now
}

// Fail marks the scan as failed
func (s *Scan) Fail(errMsg string) {
	now := time.Now()
	s.Status = ScanStatusFailed
	s.ErrorMessage = errMsg
	s.CompletedAt = &now
	s.UpdatedAt = now
}

// IsRunning returns true if the scan is running
func (s *Scan) IsRunning() bool {
	return s.Status == ScanStatusRunning
}

// IsCompleted returns true if the scan is completed
func (s *Scan) IsCompleted() bool {
	return s.Status == ScanStatusCompleted
}
