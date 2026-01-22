package repository

import (
	"context"

	"github.com/cloudsweep/cloudsweep/internal/domain/entity"
	"github.com/google/uuid"
)

// PolicyRepository defines the interface for policy persistence
type PolicyRepository interface {
	// Create creates a new policy
	Create(ctx context.Context, policy *entity.Policy) error

	// Update updates an existing policy
	Update(ctx context.Context, policy *entity.Policy) error

	// Delete deletes a policy by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByID retrieves a policy by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Policy, error)

	// List retrieves policies with filters
	List(ctx context.Context, filter PolicyFilter) ([]*entity.Policy, error)

	// GetEnabledByOrg retrieves all enabled policies for an organization
	GetEnabledByOrg(ctx context.Context, orgID uuid.UUID) ([]*entity.Policy, error)
}

// PolicyFilter defines filters for policy queries
type PolicyFilter struct {
	OrganizationID *uuid.UUID
	Provider       *entity.CloudProvider
	IsEnabled      *bool
	Limit          int
	Offset         int
}
