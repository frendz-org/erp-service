package response

import (
	"time"

	"github.com/google/uuid"
)

type MemberRegisterResponse struct {
	ID                uuid.UUID `json:"id"`
	TenantID          uuid.UUID `json:"tenant_id"`
	Status            string    `json:"status"`
	RegistrationType  string    `json:"registration_type"`
	ParticipantNumber string    `json:"participant_number"`
	IdentityNumber    string    `json:"identity_number"`
	OrganizationCode  string    `json:"organization_code"`
	FullName          string    `json:"full_name"`
	CreatedAt         time.Time `json:"created_at"`
}

type MyMemberResponse struct {
	ID                uuid.UUID  `json:"id"`
	UserID            uuid.UUID  `json:"user_id"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Email             string     `json:"email"`
	Status            string     `json:"status"`
	RegistrationType  string     `json:"registration_type"`
	RoleCode          *string    `json:"role_code,omitempty"`
	RoleName          *string    `json:"role_name,omitempty"`
	ParticipantNumber *string    `json:"participant_number,omitempty"`
	IdentityNumber    *string    `json:"identity_number,omitempty"`
	OrganizationCode  *string    `json:"organization_code,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type MemberDetailResponse struct {
	ID                uuid.UUID  `json:"id"`
	UserID            uuid.UUID  `json:"user_id"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Email             string     `json:"email"`
	Status            string     `json:"status"`
	RegistrationType  string     `json:"registration_type"`
	RoleCode          *string    `json:"role_code,omitempty"`
	RoleName          *string    `json:"role_name,omitempty"`
	ParticipantNumber *string    `json:"participant_number,omitempty"`
	IdentityNumber    *string    `json:"identity_number,omitempty"`
	OrganizationCode  *string    `json:"organization_code,omitempty"`
	ApprovedBy        *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt        *time.Time `json:"approved_at,omitempty"`
	Version           int        `json:"version"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type MemberListItemResponse struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	Status            string    `json:"status"`
	RegistrationType  string    `json:"registration_type"`
	RoleCode          *string   `json:"role_code,omitempty"`
	RoleName          *string   `json:"role_name,omitempty"`
	ParticipantNumber *string   `json:"participant_number,omitempty"`
	IdentityNumber    *string   `json:"identity_number,omitempty"`
	OrganizationCode  *string   `json:"organization_code,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type MemberListResponse struct {
	Members    []MemberListItemResponse `json:"members"`
	Pagination Pagination               `json:"pagination"`
}
