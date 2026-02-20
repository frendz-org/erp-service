package response

import (
	"time"

	"github.com/google/uuid"
)

type MemberRegisterResponse struct {
	ID               uuid.UUID `json:"id"`
	Status           string    `json:"status"`
	RegistrationType string    `json:"registration_type"`
	CreatedAt        time.Time `json:"created_at"`
}

type MemberDetailResponse struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uuid.UUID  `json:"user_id"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Email            string     `json:"email"`
	Status           string     `json:"status"`
	RegistrationType string     `json:"registration_type"`
	RoleCode         *string    `json:"role_code,omitempty"`
	RoleName         *string    `json:"role_name,omitempty"`
	ApprovedBy       *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	Version          int        `json:"version"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type MemberListItemResponse struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Email            string    `json:"email"`
	Status           string    `json:"status"`
	RegistrationType string    `json:"registration_type"`
	RoleCode         *string   `json:"role_code,omitempty"`
	RoleName         *string   `json:"role_name,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type MemberListResponse struct {
	Members    []MemberListItemResponse `json:"members"`
	Pagination Pagination               `json:"pagination"`
}
