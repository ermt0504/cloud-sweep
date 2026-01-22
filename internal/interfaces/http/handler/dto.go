package handler

import "time"

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"invalid request"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message" example:"operation successful"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data   any   `json:"data"`
	Total  int64 `json:"total" example:"100"`
	Limit  int   `json:"limit" example:"50"`
	Offset int   `json:"offset" example:"0"`
}

// ResourceDTO represents a cloud resource
type ResourceDTO struct {
	ID              string            `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrganizationID  string            `json:"organization_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Provider        string            `json:"provider" example:"aws" enums:"aws,azure,gcp"`
	Type            string            `json:"type" example:"ec2_instance"`
	ResourceID      string            `json:"resource_id" example:"i-1234567890abcdef0"`
	Region          string            `json:"region" example:"us-east-1"`
	Name            string            `json:"name" example:"my-instance"`
	Status          string            `json:"status" example:"unused" enums:"active,unused,deleted,excluded"`
	Tags            map[string]string `json:"tags"`
	MonthlyCost     float64           `json:"monthly_cost" example:"45.50"`
	CarbonFootprint float64           `json:"carbon_footprint_kg" example:"12.5"`
	LastSeenAt      time.Time         `json:"last_seen_at"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// ScanDTO represents a scan
type ScanDTO struct {
	ID               string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrganizationID   string    `json:"organization_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Provider         string    `json:"provider" example:"aws" enums:"aws,azure,gcp"`
	Regions          []string  `json:"regions" example:"us-east-1,eu-west-1"`
	ResourceTypes    []string  `json:"resource_types" example:"ec2_instance,ebs_volume"`
	Status           string    `json:"status" example:"completed" enums:"pending,running,completed,failed,cancelled"`
	ResourcesFound   int       `json:"resources_found" example:"150"`
	UnusedFound      int       `json:"unused_found" example:"23"`
	EstimatedSavings float64   `json:"estimated_savings" example:"1250.00"`
	CarbonSavings    float64   `json:"carbon_savings_kg" example:"45.5"`
	ErrorMessage     string    `json:"error_message,omitempty"`
	StartedAt        *time.Time `json:"started_at,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// PolicyDTO represents a cleanup policy
type PolicyDTO struct {
	ID             string         `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrganizationID string         `json:"organization_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Name           string         `json:"name" example:"Delete unused EBS volumes"`
	Description    string         `json:"description" example:"Automatically delete EBS volumes unused for 30 days"`
	Provider       string         `json:"provider" example:"aws" enums:"aws,azure,gcp"`
	ResourceTypes  []string       `json:"resource_types" example:"ebs_volume"`
	Conditions     map[string]any `json:"conditions"`
	Actions        []string       `json:"actions" example:"notify,delete" enums:"notify,tag,stop,delete"`
	IsEnabled      bool           `json:"is_enabled" example:"true"`
	Schedule       string         `json:"schedule" example:"0 0 * * *"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// DashboardSummaryDTO represents dashboard summary
type DashboardSummaryDTO struct {
	TotalResources   int64   `json:"total_resources" example:"500"`
	UnusedResources  int64   `json:"unused_resources" example:"75"`
	TotalMonthlyCost float64 `json:"total_monthly_cost" example:"15000.00"`
	PotentialSavings float64 `json:"potential_monthly_savings" example:"2500.00"`
	TotalCarbonKg    float64 `json:"total_carbon_kg" example:"1200.50"`
	CarbonSavingsKg  float64 `json:"potential_carbon_savings_kg" example:"180.25"`
}

// CleanupPreviewDTO represents a cleanup preview response
type CleanupPreviewDTO struct {
	Resources             []ResourceDTO `json:"resources"`
	Count                 int           `json:"count" example:"5"`
	EstimatedMonthlySavings float64     `json:"estimated_monthly_savings" example:"250.00"`
	EstimatedCarbonSavings  float64     `json:"estimated_carbon_savings" example:"35.5"`
	Action                string        `json:"action" example:"delete"`
}
