package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Member struct {
	ID                        uuid.UUID    `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	UserTenantRegistrationID  uuid.UUID    `json:"user_tenant_registration_id" gorm:"column:user_tenant_registration_id;type:uuid;not null" db:"user_tenant_registration_id"`
	TenantID                  uuid.UUID    `json:"tenant_id" gorm:"column:tenant_id;type:uuid;not null" db:"tenant_id"`
	ProductID                 uuid.UUID    `json:"product_id" gorm:"column:product_id;type:uuid;not null" db:"product_id"`
	UserID                    uuid.UUID    `json:"user_id" gorm:"column:user_id;type:uuid;not null" db:"user_id"`
	ParticipantNumber         string       `json:"participant_number" gorm:"column:participant_number;type:varchar(50);not null" db:"participant_number"`
	IdentityNumber            string       `json:"identity_number" gorm:"column:identity_number;type:varchar(16);not null" db:"identity_number"`
	OrganizationCode          string       `json:"organization_code" gorm:"column:organization_code;type:varchar(50);not null" db:"organization_code"`
	FullName                  string       `json:"full_name" gorm:"column:full_name;type:varchar(255);not null" db:"full_name"`
	Gender                    *string      `json:"gender,omitempty" gorm:"column:gender;type:varchar(20)" db:"gender"`
	DateOfBirth               *time.Time   `json:"date_of_birth,omitempty" gorm:"column:date_of_birth" db:"date_of_birth"`
	Version                   int          `json:"version" gorm:"column:version;not null;default:1" db:"version"`
	CreatedAt                 time.Time    `json:"created_at" gorm:"column:created_at;not null" db:"created_at"`
	UpdatedAt                 time.Time    `json:"updated_at" gorm:"column:updated_at;not null" db:"updated_at"`
	DeletedAt                 sql.NullTime `json:"deleted_at,omitempty" gorm:"column:deleted_at" db:"deleted_at"`
}

func (Member) TableName() string {
	return "members"
}
