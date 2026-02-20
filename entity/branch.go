package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Branch struct {
	ID       uuid.UUID       `json:"id" db:"id"`
	TenantID uuid.UUID       `json:"tenant_id" db:"tenant_id"`
	Code     string          `json:"code" db:"code"`
	Name     string          `json:"name" db:"name"`
	Address  *string         `json:"address,omitempty" db:"address"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	Status   string          `json:"status" db:"status"`
	Version  int             `json:"version" db:"version"`
	Timestamps
}

func (Branch) TableName() string {
	return "branches"
}

func (b *Branch) IsActive() bool {
	return b.Status == "ACTIVE"
}

type BranchContact struct {
	ID         uuid.UUID `json:"id" db:"id"`
	BranchID   uuid.UUID `json:"branch_id" db:"branch_id"`
	Address    string    `json:"address,omitempty" db:"address"`
	City       string    `json:"city,omitempty" db:"city"`
	Province   string    `json:"province,omitempty" db:"province"`
	PostalCode string    `json:"postal_code,omitempty" db:"postal_code"`
	Phone      string    `json:"phone,omitempty" db:"phone"`
	Email      string    `json:"email,omitempty" db:"email"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
