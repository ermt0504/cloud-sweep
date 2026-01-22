package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// JSONB represents a JSONB field
type JSONB map[string]any

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

// StringArray represents a PostgreSQL text array
type StringArray []string

// Value implements the driver.Valuer interface
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface
func (a *StringArray) Scan(value any) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, a)
}

// Organization represents the organizations table
type Organization struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Slug      string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Plan      string    `gorm:"type:varchar(50);default:'free'"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// CloudAccount represents the cloud_accounts table
type CloudAccount struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID `gorm:"type:uuid;index;not null"`
	Provider       string    `gorm:"type:varchar(20);not null"`
	AccountID      string    `gorm:"type:varchar(255);not null"`
	Name           string    `gorm:"type:varchar(255)"`
	Credentials    []byte    `gorm:"type:bytea"`
	IsActive       bool      `gorm:"default:true"`
	LastSyncAt     *time.Time
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`

	Organization Organization `gorm:"foreignKey:OrganizationID"`
}

// Resource represents the resources table
type Resource struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID  uuid.UUID `gorm:"type:uuid;index;not null"`
	Provider        string    `gorm:"type:varchar(20);index;not null"`
	Type            string    `gorm:"type:varchar(50);index;not null"`
	ResourceID      string    `gorm:"type:varchar(255);index;not null"`
	Region          string    `gorm:"type:varchar(50);index"`
	Name            string    `gorm:"type:varchar(255)"`
	Status          string    `gorm:"type:varchar(20);index;default:'active'"`
	Tags            JSONB     `gorm:"type:jsonb"`
	Metadata        JSONB     `gorm:"type:jsonb"`
	MonthlyCost     float64   `gorm:"type:decimal(10,2);default:0"`
	CarbonFootprint float64   `gorm:"type:decimal(10,4);default:0"`
	LastSeenAt      time.Time
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`

	Organization Organization `gorm:"foreignKey:OrganizationID"`
}

// Scan represents the scans table
type Scan struct {
	ID               uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID   uuid.UUID   `gorm:"type:uuid;index;not null"`
	Provider         string      `gorm:"type:varchar(20);not null"`
	Regions          StringArray `gorm:"type:jsonb"`
	ResourceTypes    StringArray `gorm:"type:jsonb"`
	Status           string      `gorm:"type:varchar(20);index;default:'pending'"`
	ResourcesFound   int         `gorm:"default:0"`
	UnusedFound      int         `gorm:"default:0"`
	EstimatedSavings float64     `gorm:"type:decimal(10,2);default:0"`
	CarbonSavings    float64     `gorm:"type:decimal(10,4);default:0"`
	ErrorMessage     string      `gorm:"type:text"`
	StartedAt        *time.Time
	CompletedAt      *time.Time
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`

	Organization Organization `gorm:"foreignKey:OrganizationID"`
}

// Policy represents the policies table
type Policy struct {
	ID             uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID   `gorm:"type:uuid;index;not null"`
	Name           string      `gorm:"type:varchar(255);not null"`
	Description    string      `gorm:"type:text"`
	Provider       string      `gorm:"type:varchar(20);not null"`
	ResourceTypes  StringArray `gorm:"type:jsonb"`
	Conditions     JSONB       `gorm:"type:jsonb"`
	Actions        StringArray `gorm:"type:jsonb"`
	IsEnabled      bool        `gorm:"default:true"`
	Schedule       string      `gorm:"type:varchar(100)"`
	CreatedAt      time.Time   `gorm:"autoCreateTime"`
	UpdatedAt      time.Time   `gorm:"autoUpdateTime"`

	Organization Organization `gorm:"foreignKey:OrganizationID"`
}

// TableName overrides
func (Organization) TableName() string  { return "organizations" }
func (CloudAccount) TableName() string  { return "cloud_accounts" }
func (Resource) TableName() string      { return "resources" }
func (Scan) TableName() string          { return "scans" }
func (Policy) TableName() string        { return "policies" }
