package entity

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type UserTenantRegistrationStatus string

const (
	UTRStatusPendingApproval UserTenantRegistrationStatus = "PENDING_APPROVAL"
	UTRStatusActive          UserTenantRegistrationStatus = "ACTIVE"
	UTRStatusRejected        UserTenantRegistrationStatus = "REJECTED"
	UTRStatusInactive        UserTenantRegistrationStatus = "INACTIVE"
)

type UserTenantRegistration struct {
	ID                   uuid.UUID                    `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	UserID               uuid.UUID                    `json:"user_id" gorm:"column:user_id;type:uuid;not null" db:"user_id"`
	TenantID             uuid.UUID                    `json:"tenant_id" gorm:"column:tenant_id;type:uuid;not null" db:"tenant_id"`
	ProductID            *uuid.UUID                   `json:"product_id,omitempty" gorm:"column:product_id;type:uuid" db:"product_id"`
	RegistrationType     string                       `json:"registration_type" gorm:"column:registration_type;type:varchar(20);not null" db:"registration_type"`
	IdentificationNumber *string                      `json:"identification_number,omitempty" gorm:"column:identification_number;type:varchar(100)" db:"identification_number"`
	Status               UserTenantRegistrationStatus `json:"status" gorm:"column:status;type:varchar(20);not null;default:PENDING_APPROVAL" db:"status"`
	ApprovedBy           *uuid.UUID                   `json:"approved_by,omitempty" gorm:"column:approved_by;type:uuid" db:"approved_by"`
	ApprovedAt           *time.Time                   `json:"approved_at,omitempty" gorm:"column:approved_at" db:"approved_at"`
	Metadata             json.RawMessage              `json:"metadata" gorm:"column:metadata;type:jsonb;not null;default:'{}'" db:"metadata"`
	CreatedAt            time.Time                    `json:"created_at" gorm:"column:created_at;not null" db:"created_at"`
	UpdatedAt            time.Time                    `json:"updated_at" gorm:"column:updated_at;not null" db:"updated_at"`
	DeletedAt            sql.NullTime                 `json:"deleted_at,omitempty" gorm:"column:deleted_at" db:"deleted_at"`
	Version              int                          `json:"version" gorm:"column:version;not null;default:1" db:"version"`
}

func (UserTenantRegistration) TableName() string {
	return "user_tenant_registrations"
}
