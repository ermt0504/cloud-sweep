package usecase

import (
	"context"
	"fmt"

	"github.com/cloudsweep/cloudsweep/internal/domain/entity"
	"github.com/cloudsweep/cloudsweep/internal/domain/repository"
	"github.com/cloudsweep/cloudsweep/internal/domain/service"
	"github.com/google/uuid"
)

// CleanupResourcesUseCase handles resource cleanup operations
type CleanupResourcesUseCase struct {
	resourceRepo   repository.ResourceRepository
	policyRepo     repository.PolicyRepository
	cleanerFactory service.ResourceCleanerFactory
}

// NewCleanupResourcesUseCase creates a new CleanupResourcesUseCase
func NewCleanupResourcesUseCase(
	resourceRepo repository.ResourceRepository,
	policyRepo repository.PolicyRepository,
	cleanerFactory service.ResourceCleanerFactory,
) *CleanupResourcesUseCase {
	return &CleanupResourcesUseCase{
		resourceRepo:   resourceRepo,
		policyRepo:     policyRepo,
		cleanerFactory: cleanerFactory,
	}
}

// CleanupResourcesInput represents input for cleaning up resources
type CleanupResourcesInput struct {
	OrganizationID uuid.UUID
	ResourceIDs    []uuid.UUID
	Action         entity.PolicyAction
	Credentials    []byte
	DryRun         bool
}

// CleanupResourcesOutput represents output from cleaning up resources
type CleanupResourcesOutput struct {
	Results       []*service.CleanupResult
	TotalCostSaved   float64
	TotalCarbonSaved float64
	SuccessCount     int
	FailureCount     int
}

// Execute executes the cleanup resources use case
func (uc *CleanupResourcesUseCase) Execute(ctx context.Context, input CleanupResourcesInput) (*CleanupResourcesOutput, error) {
	output := &CleanupResourcesOutput{
		Results: make([]*service.CleanupResult, 0, len(input.ResourceIDs)),
	}

	// Get resources
	var resources []*entity.Resource
	for _, id := range input.ResourceIDs {
		resource, err := uc.resourceRepo.GetByID(ctx, id)
		if err != nil {
			output.Results = append(output.Results, &service.CleanupResult{
				ResourceID:   id.String(),
				Success:      false,
				ErrorMessage: fmt.Sprintf("resource not found: %v", err),
			})
			output.FailureCount++
			continue
		}
		resources = append(resources, resource)
	}

	if len(resources) == 0 {
		return output, nil
	}

	// Group resources by provider
	resourcesByProvider := make(map[entity.CloudProvider][]*entity.Resource)
	for _, r := range resources {
		resourcesByProvider[r.Provider] = append(resourcesByProvider[r.Provider], r)
	}

	// Process each provider
	for provider, providerResources := range resourcesByProvider {
		cleaner, err := uc.cleanerFactory.Create(provider, input.Credentials)
		if err != nil {
			for _, r := range providerResources {
				output.Results = append(output.Results, &service.CleanupResult{
					ResourceID:   r.ID.String(),
					Success:      false,
					ErrorMessage: fmt.Sprintf("failed to create cleaner: %v", err),
				})
				output.FailureCount++
			}
			continue
		}

		// Process each resource
		for _, resource := range providerResources {
			if input.DryRun {
				output.Results = append(output.Results, &service.CleanupResult{
					ResourceID:  resource.ID.String(),
					Success:     true,
					Action:      input.Action,
					CostSaved:   resource.MonthlyCost,
					CarbonSaved: resource.CarbonFootprint,
				})
				output.TotalCostSaved += resource.MonthlyCost
				output.TotalCarbonSaved += resource.CarbonFootprint
				output.SuccessCount++
				continue
			}

			var result *service.CleanupResult
			switch input.Action {
			case entity.PolicyActionDelete:
				result, err = cleaner.Delete(ctx, resource)
			case entity.PolicyActionStop:
				result, err = cleaner.Stop(ctx, resource)
			case entity.PolicyActionTag:
				result, err = cleaner.Tag(ctx, resource, map[string]string{
					"cloudsweep:marked-for-deletion": "true",
				})
			default:
				result = &service.CleanupResult{
					ResourceID:   resource.ID.String(),
					Success:      false,
					ErrorMessage: "unsupported action",
				}
			}

			if err != nil {
				result = &service.CleanupResult{
					ResourceID:   resource.ID.String(),
					Success:      false,
					ErrorMessage: err.Error(),
				}
			}

			output.Results = append(output.Results, result)
			if result.Success {
				output.TotalCostSaved += result.CostSaved
				output.TotalCarbonSaved += result.CarbonSaved
				output.SuccessCount++

				// Update resource status
				resource.MarkAsDeleted()
				uc.resourceRepo.Update(ctx, resource)
			} else {
				output.FailureCount++
			}
		}
	}

	return output, nil
}
