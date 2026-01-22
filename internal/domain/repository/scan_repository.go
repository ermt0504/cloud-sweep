package repository

import (
	"context"

	"github.com/cloudsweep/cloudsweep/internal/domain/entity"
	"github.com/google/uuid"
)

// ScanRepository defines the interface for scan persistence
type ScanRepository interface {
	// Create creates a new scan
	Create(ctx context.Context, scan *entity.Scan) error

	// Update updates an existing scan
	Update(ctx context.Context, scan *entity.Scan) error

	// GetByID retrieves a scan by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Scan, error)

	// List retrieves scans with filters
	List(ctx context.Context, filter ScanFilter) ([]*entity.Scan, error)

	// GetLatestByOrg retrieves the latest scan for an organization
	GetLatestByOrg(ctx context.Context, orgID uuid.UUID) (*entity.Scan, error)
}

// ScanFilter defines filters for scan queries
type ScanFilter struct {
	OrganizationID *uuid.UUID
	Provider       *entity.CloudProvider
	Status         *entity.ScanStatus
	Limit          int
	Offset         int
}
