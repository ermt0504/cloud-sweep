package entity

import (
	"time"

	"github.com/google/uuid"
)

// CloudProvider represents a cloud provider
type CloudProvider string

const (
	CloudProviderAWS   CloudProvider = "aws"
	CloudProviderAzure CloudProvider = "azure"
	CloudProviderGCP   CloudProvider = "gcp"
)

// ResourceType represents a type of cloud resource
type ResourceType string

const (
	ResourceTypeEC2Instance   ResourceType = "ec2_instance"
	ResourceTypeEBSVolume     ResourceType = "ebs_volume"
	ResourceTypeEBSSnapshot   ResourceType = "ebs_snapshot"
	ResourceTypeElasticIP     ResourceType = "elastic_ip"
	ResourceTypeLoadBalancer  ResourceType = "load_balancer"
	ResourceTypeS3Bucket      ResourceType = "s3_bucket"
	ResourceTypeRDSInstance   ResourceType = "rds_instance"
	ResourceTypeAzureVM       ResourceType = "azure_vm"
	ResourceTypeAzureDisk     ResourceType = "azure_disk"
	ResourceTypeGCEInstance   ResourceType = "gce_instance"
	ResourceTypeGCEDisk       ResourceType = "gce_disk"
)

// ResourceStatus represents the status of a resource
type ResourceStatus string

const (
	ResourceStatusActive   ResourceStatus = "active"
	ResourceStatusUnused   ResourceStatus = "unused"
	ResourceStatusDeleted  ResourceStatus = "deleted"
	ResourceStatusExcluded ResourceStatus = "excluded"
)

// Resource represents a cloud resource
type Resource struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organization_id"`
	Provider       CloudProvider   `json:"provider"`
	Type           ResourceType    `json:"type"`
	ResourceID     string          `json:"resource_id"`
	Region         string          `json:"region"`
	Name           string          `json:"name"`
	Status         ResourceStatus  `json:"status"`
	Tags           map[string]string `json:"tags"`
	Metadata       map[string]any  `json:"metadata"`
	MonthlyCost    float64         `json:"monthly_cost"`
	CarbonFootprint float64        `json:"carbon_footprint_kg"`
	LastSeenAt     time.Time       `json:"last_seen_at"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// NewResource creates a new Resource
func NewResource(orgID uuid.UUID, provider CloudProvider, resourceType ResourceType, resourceID, region, name string) *Resource {
	now := time.Now()
	return &Resource{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Provider:       provider,
		Type:           resourceType,
		ResourceID:     resourceID,
		Region:         region,
		Name:           name,
		Status:         ResourceStatusActive,
		Tags:           make(map[string]string),
		Metadata:       make(map[string]any),
		LastSeenAt:     now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// MarkAsUnused marks the resource as unused
func (r *Resource) MarkAsUnused() {
	r.Status = ResourceStatusUnused
	r.UpdatedAt = time.Now()
}

// MarkAsDeleted marks the resource as deleted
func (r *Resource) MarkAsDeleted() {
	r.Status = ResourceStatusDeleted
	r.UpdatedAt = time.Now()
}

// IsUnused returns true if the resource is unused
func (r *Resource) IsUnused() bool {
	return r.Status == ResourceStatusUnused
}
