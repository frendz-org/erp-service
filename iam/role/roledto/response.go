package roledto

import (
	"time"

	"github.com/google/uuid"
)

type CreateResponse struct {
	RoleID      uuid.UUID `json:"role_id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	ScopeLevel  string    `json:"scope_level"`
	IsSystem    bool      `json:"is_system"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}
