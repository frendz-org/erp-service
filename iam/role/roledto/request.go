package roledto

import "github.com/google/uuid"

type CreateRequest struct {
	TenantID    uuid.UUID   `json:"tenant_id" validate:"required"`
	Code        string      `json:"code" validate:"required,min=2,max=100,uppercase"`
	Name        string      `json:"name" validate:"required,min=2,max=255"`
	Description *string     `json:"description,omitempty" validate:"omitempty,max=1000"`
	ScopeLevel  string      `json:"scope_level" validate:"required,oneof=tenant branch self"`
	Permissions []uuid.UUID `json:"permissions,omitempty" validate:"omitempty,dive,uuid"`
}
