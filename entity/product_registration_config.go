package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProductRegistrationConfig struct {
	ID               uuid.UUID  `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	ProductID        uuid.UUID  `json:"product_id" gorm:"column:product_id;type:uuid;not null" db:"product_id"`
	RegistrationType string     `json:"registration_type" gorm:"column:registration_type;type:varchar(20);not null" db:"registration_type"`
	AutoGrantRoleID  *uuid.UUID `json:"auto_grant_role_id,omitempty" gorm:"column:auto_grant_role_id;type:uuid" db:"auto_grant_role_id"`
	RequiresApproval bool       `json:"requires_approval" gorm:"column:requires_approval;not null;default:true" db:"requires_approval"`
	IsActive         bool       `json:"is_active" gorm:"column:is_active;not null;default:true" db:"is_active"`
	CreatedAt        time.Time  `json:"created_at" gorm:"column:created_at;not null" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"column:updated_at;not null" db:"updated_at"`
}

func (ProductRegistrationConfig) TableName() string {
	return "product_registration_configs"
}
