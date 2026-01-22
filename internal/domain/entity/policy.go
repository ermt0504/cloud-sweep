package entity

import (
	"time"

	"github.com/google/uuid"
)

// PolicyAction represents an action to take
type PolicyAction string

const (
	PolicyActionNotify  PolicyAction = "notify"
	PolicyActionTag     PolicyAction = "tag"
	PolicyActionStop    PolicyAction = "stop"
	PolicyActionDelete  PolicyAction = "delete"
)

// Policy represents a cleanup policy
type Policy struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organization_id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Provider       CloudProvider   `json:"provider"`
	ResourceTypes  []ResourceType  `json:"resource_types"`
	Conditions     PolicyConditions `json:"conditions"`
	Actions        []PolicyAction  `json:"actions"`
	IsEnabled      bool            `json:"is_enabled"`
	Schedule       string          `json:"schedule"` // Cron expression
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// PolicyConditions defines when a policy should apply
type PolicyConditions struct {
	UnusedDays       int               `json:"unused_days,omitempty"`
	MinMonthlyCost   float64           `json:"min_monthly_cost,omitempty"`
	MaxMonthlyCost   float64           `json:"max_monthly_cost,omitempty"`
	RequiredTags     map[string]string `json:"required_tags,omitempty"`
	ExcludedTags     map[string]string `json:"excluded_tags,omitempty"`
	Regions          []string          `json:"regions,omitempty"`
	NamePattern      string            `json:"name_pattern,omitempty"`
}

// NewPolicy creates a new Policy
func NewPolicy(orgID uuid.UUID, name, description string, provider CloudProvider) *Policy {
	now := time.Now()
	return &Policy{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
		Description:    description,
		Provider:       provider,
		ResourceTypes:  []ResourceType{},
		Conditions:     PolicyConditions{},
		Actions:        []PolicyAction{},
		IsEnabled:      true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Enable enables the policy
func (p *Policy) Enable() {
	p.IsEnabled = true
	p.UpdatedAt = time.Now()
}

// Disable disables the policy
func (p *Policy) Disable() {
	p.IsEnabled = false
	p.UpdatedAt = time.Now()
}

// HasDeleteAction returns true if the policy includes delete action
func (p *Policy) HasDeleteAction() bool {
	for _, action := range p.Actions {
		if action == PolicyActionDelete {
			return true
		}
	}
	return false
}
