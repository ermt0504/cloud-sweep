package entity

import (
	"time"

	"github.com/google/uuid"
)

// Organization represents a customer organization
type Organization struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Plan      string    `json:"plan"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewOrganization creates a new Organization
func NewOrganization(name, slug string) *Organization {
	now := time.Now()
	return &Organization{
		ID:        uuid.New(),
		Name:      name,
		Slug:      slug,
		Plan:      "free",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CloudAccount represents a connected cloud account
type CloudAccount struct {
	ID             uuid.UUID     `json:"id"`
	OrganizationID uuid.UUID     `json:"organization_id"`
	Provider       CloudProvider `json:"provider"`
	AccountID      string        `json:"account_id"`
	Name           string        `json:"name"`
	Credentials    []byte        `json:"-"` // Encrypted credentials
	IsActive       bool          `json:"is_active"`
	LastSyncAt     *time.Time    `json:"last_sync_at,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// NewCloudAccount creates a new CloudAccount
func NewCloudAccount(orgID uuid.UUID, provider CloudProvider, accountID, name string) *CloudAccount {
	now := time.Now()
	return &CloudAccount{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Provider:       provider,
		AccountID:      accountID,
		Name:           name,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
