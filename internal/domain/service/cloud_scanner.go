package service

import (
	"context"

	"github.com/cloudsweep/cloudsweep/internal/domain/entity"
)

// CloudScanner defines the interface for scanning cloud resources
type CloudScanner interface {
	// ScanResources scans for resources of specified types in given regions
	ScanResources(ctx context.Context, regions []string, resourceTypes []entity.ResourceType) ([]*entity.Resource, error)

	// DetectUnused analyzes resources and marks unused ones
	DetectUnused(ctx context.Context, resources []*entity.Resource) error

	// EstimateCost estimates the monthly cost of a resource
	EstimateCost(ctx context.Context, resource *entity.Resource) (float64, error)

	// EstimateCarbonFootprint estimates the carbon footprint of a resource
	EstimateCarbonFootprint(ctx context.Context, resource *entity.Resource) (float64, error)

	// Provider returns the cloud provider
	Provider() entity.CloudProvider
}

// CloudScannerFactory creates cloud scanners based on provider
type CloudScannerFactory interface {
	// Create creates a scanner for the given provider and credentials
	Create(provider entity.CloudProvider, credentials []byte) (CloudScanner, error)
}
