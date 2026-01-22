package repository

import (
	"context"

	"github.com/cloudsweep/cloudsweep/internal/domain/entity"
	"github.com/google/uuid"
)

// ResourceRepository defines the interface for resource persistence
type ResourceRepository interface {
	// Create creates a new resource
	Create(ctx context.Context, resource *entity.Resource) error

	// Update updates an existing resource
	Update(ctx context.Context, resource *entity.Resource) error

	// Delete deletes a resource by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByID retrieves a resource by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Resource, error)

	// GetByResourceID retrieves a resource by cloud resource ID
	GetByResourceID(ctx context.Context, orgID uuid.UUID, provider entity.CloudProvider, resourceID string) (*entity.Resource, error)

	// List retrieves resources with filters
	List(ctx context.Context, filter ResourceFilter) ([]*entity.Resource, error)

	// Count counts resources with filters
	Count(ctx context.Context, filter ResourceFilter) (int64, error)

	// BulkCreate creates multiple resources
	BulkCreate(ctx context.Context, resources []*entity.Resource) error

	// BulkUpdate updates multiple resources
	BulkUpdate(ctx context.Context, resources []*entity.Resource) error
}

// ResourceFilter defines filters for resource queries
type ResourceFilter struct {
	OrganizationID *uuid.UUID
	Provider       *entity.CloudProvider
	Type           *entity.ResourceType
	Status         *entity.ResourceStatus
	Region         *string
	Limit          int
	Offset         int
}
