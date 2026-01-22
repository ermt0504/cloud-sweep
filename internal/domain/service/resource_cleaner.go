package service

import (
	"context"

	"github.com/cloudsweep/cloudsweep/internal/domain/entity"
)

// CleanupResult represents the result of a cleanup operation
type CleanupResult struct {
	ResourceID    string
	Success       bool
	Action        entity.PolicyAction
	ErrorMessage  string
	CostSaved     float64
	CarbonSaved   float64
}

// ResourceCleaner defines the interface for cleaning up cloud resources
type ResourceCleaner interface {
	// Delete permanently deletes a resource
	Delete(ctx context.Context, resource *entity.Resource) (*CleanupResult, error)

	// Stop stops a running resource (e.g., EC2 instance)
	Stop(ctx context.Context, resource *entity.Resource) (*CleanupResult, error)

	// Tag adds tags to a resource
	Tag(ctx context.Context, resource *entity.Resource, tags map[string]string) (*CleanupResult, error)

	// Provider returns the cloud provider
	Provider() entity.CloudProvider
}

// ResourceCleanerFactory creates resource cleaners based on provider
type ResourceCleanerFactory interface {
	// Create creates a cleaner for the given provider and credentials
	Create(provider entity.CloudProvider, credentials []byte) (ResourceCleaner, error)
}
