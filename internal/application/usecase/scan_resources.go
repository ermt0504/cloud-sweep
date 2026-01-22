package usecase

import (
	"context"
	"fmt"

	"github.com/cloudsweep/cloudsweep/internal/domain/entity"
	"github.com/cloudsweep/cloudsweep/internal/domain/repository"
	"github.com/cloudsweep/cloudsweep/internal/domain/service"
	"github.com/google/uuid"
)

// ScanResourcesUseCase handles resource scanning operations
type ScanResourcesUseCase struct {
	scanRepo       repository.ScanRepository
	resourceRepo   repository.ResourceRepository
	scannerFactory service.CloudScannerFactory
}

// NewScanResourcesUseCase creates a new ScanResourcesUseCase
func NewScanResourcesUseCase(
	scanRepo repository.ScanRepository,
	resourceRepo repository.ResourceRepository,
	scannerFactory service.CloudScannerFactory,
) *ScanResourcesUseCase {
	return &ScanResourcesUseCase{
		scanRepo:       scanRepo,
		resourceRepo:   resourceRepo,
		scannerFactory: scannerFactory,
	}
}

// ScanResourcesInput represents input for scanning resources
type ScanResourcesInput struct {
	OrganizationID uuid.UUID
	Provider       entity.CloudProvider
	Regions        []string
	ResourceTypes  []entity.ResourceType
	Credentials    []byte
}

// ScanResourcesOutput represents output from scanning resources
type ScanResourcesOutput struct {
	ScanID           uuid.UUID
	ResourcesFound   int
	UnusedFound      int
	EstimatedSavings float64
	CarbonSavings    float64
}

// Execute executes the scan resources use case
func (uc *ScanResourcesUseCase) Execute(ctx context.Context, input ScanResourcesInput) (*ScanResourcesOutput, error) {
	// Create scan record
	scan := entity.NewScan(input.OrganizationID, input.Provider, input.Regions, input.ResourceTypes)
	if err := uc.scanRepo.Create(ctx, scan); err != nil {
		return nil, fmt.Errorf("failed to create scan: %w", err)
	}

	// Start scan
	scan.Start()
	if err := uc.scanRepo.Update(ctx, scan); err != nil {
		return nil, fmt.Errorf("failed to update scan status: %w", err)
	}

	// Create scanner
	scanner, err := uc.scannerFactory.Create(input.Provider, input.Credentials)
	if err != nil {
		scan.Fail(err.Error())
		uc.scanRepo.Update(ctx, scan)
		return nil, fmt.Errorf("failed to create scanner: %w", err)
	}

	// Scan resources
	resources, err := scanner.ScanResources(ctx, input.Regions, input.ResourceTypes)
	if err != nil {
		scan.Fail(err.Error())
		uc.scanRepo.Update(ctx, scan)
		return nil, fmt.Errorf("failed to scan resources: %w", err)
	}

	// Set organization ID for all resources
	for _, r := range resources {
		r.OrganizationID = input.OrganizationID
	}

	// Detect unused resources
	if err := scanner.DetectUnused(ctx, resources); err != nil {
		scan.Fail(err.Error())
		uc.scanRepo.Update(ctx, scan)
		return nil, fmt.Errorf("failed to detect unused resources: %w", err)
	}

	// Calculate costs and carbon footprint
	var totalSavings, totalCarbon float64
	unusedCount := 0
	for _, r := range resources {
		cost, _ := scanner.EstimateCost(ctx, r)
		carbon, _ := scanner.EstimateCarbonFootprint(ctx, r)
		r.MonthlyCost = cost
		r.CarbonFootprint = carbon

		if r.IsUnused() {
			unusedCount++
			totalSavings += cost
			totalCarbon += carbon
		}
	}

	// Save resources
	if err := uc.resourceRepo.BulkCreate(ctx, resources); err != nil {
		scan.Fail(err.Error())
		uc.scanRepo.Update(ctx, scan)
		return nil, fmt.Errorf("failed to save resources: %w", err)
	}

	// Complete scan
	scan.Complete(len(resources), unusedCount, totalSavings, totalCarbon)
	if err := uc.scanRepo.Update(ctx, scan); err != nil {
		return nil, fmt.Errorf("failed to complete scan: %w", err)
	}

	return &ScanResourcesOutput{
		ScanID:           scan.ID,
		ResourcesFound:   len(resources),
		UnusedFound:      unusedCount,
		EstimatedSavings: totalSavings,
		CarbonSavings:    totalCarbon,
	}, nil
}
