package userdto

import "github.com/google/uuid"

type CreateRequest struct {
	TenantID uuid.UUID `json:"tenant_id" validate:"required"`
	RoleCode string    `json:"role_code" validate:"required"`
	Email    string    `json:"email" validate:"required,email,max=255"`
	Password string    `json:"password" validate:"required,min=8,max=128"`
	FirstName string   `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string   `json:"last_name" validate:"required,min=2,max=100"`
}

type ListRequest struct {
	Page      int        `query:"page" validate:"min=1"`
	PerPage   int        `query:"per_page" validate:"min=1,max=100"`
	SortBy    string     `query:"sort_by" validate:"omitempty,oneof=created_at updated_at email first_name last_name"`
	SortOrder string     `query:"sort_order" validate:"omitempty,oneof=asc desc"`
	Status    string     `query:"status" validate:"omitempty"`
	Search    string     `query:"search" validate:"omitempty,max=100"`
	RoleID    *uuid.UUID `query:"role_id" validate:"omitempty"`
}

func (r *ListRequest) SetDefaults() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PerPage <= 0 {
		r.PerPage = 20
	}
	if r.PerPage > 100 {
		r.PerPage = 100
	}
	if r.SortBy == "" {
		r.SortBy = "created_at"
	}
	if r.SortOrder == "" {
		r.SortOrder = "desc"
	}
}

type UpdateMeRequest struct {
	FirstName   *string `json:"first_name,omitempty" validate:"omitempty,min=2,max=100"`
	LastName    *string `json:"last_name,omitempty" validate:"omitempty,min=2,max=100"`
	PhoneNumber *string `json:"phone_number,omitempty" validate:"omitempty,max=50"`
	Address     *string `json:"address,omitempty" validate:"omitempty,max=500"`
}

type UpdateRequest struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,max=50"`
	Address   *string `json:"address,omitempty" validate:"omitempty,max=500"`
	Status    *string `json:"status,omitempty" validate:"omitempty"`
}

type RejectRequest struct {
	Reason string `json:"reason" validate:"required,max=500"`
}
