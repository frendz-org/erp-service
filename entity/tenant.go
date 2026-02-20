package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "ACTIVE"
	TenantStatusInactive  TenantStatus = "INACTIVE"
	TenantStatusSuspended TenantStatus = "SUSPENDED"
)

type Tenant struct {
	ID         uuid.UUID       `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	Code       string          `json:"code" gorm:"column:code;uniqueIndex" db:"code"`
	Name       string          `json:"name" gorm:"column:name" db:"name"`
	Settings   json.RawMessage `json:"settings" gorm:"column:settings;type:jsonb;not null;default:'{}'" db:"settings"`
	TenantType *string         `json:"tenant_type,omitempty" gorm:"column:tenant_type" db:"tenant_type"`
	Status     TenantStatus    `json:"status" gorm:"column:status;default:ACTIVE" db:"status"`
	Version    int             `json:"version" gorm:"column:version;default:1" db:"version"`
	Timestamps
}

func (Tenant) TableName() string {
	return "tenants"
}

func (t *Tenant) IsActive() bool {
	return t.Status == TenantStatusActive
}

type TenantSettings struct {
	ID               uuid.UUID `json:"id" gorm:"column:id;primaryKey" db:"id"`
	TenantID         uuid.UUID `json:"tenant_id" gorm:"column:tenant_id;uniqueIndex" db:"tenant_id"`
	SubscriptionTier string    `json:"subscription_tier" gorm:"column:subscription_tier;default:standard" db:"subscription_tier"`
	MaxBranches      int       `json:"max_branches" gorm:"column:max_branches;default:10" db:"max_branches"`
	MaxEmployees     int       `json:"max_employees" gorm:"column:max_employees;default:10000" db:"max_employees"`
	ContactEmail     string    `json:"contact_email,omitempty" gorm:"column:contact_email" db:"contact_email"`
	ContactPhone     string    `json:"contact_phone,omitempty" gorm:"column:contact_phone" db:"contact_phone"`
	ContactAddress   string    `json:"contact_address,omitempty" gorm:"column:contact_address" db:"contact_address"`
	DefaultLanguage  string    `json:"default_language" gorm:"column:default_language;default:en" db:"default_language"`
	Timezone         string    `json:"timezone" gorm:"column:timezone;default:Asia/Jakarta" db:"timezone"`
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
}

func (TenantSettings) TableName() string {
	return "tenant_settings"
}
