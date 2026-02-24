package user

import (
	"time"

	"github.com/google/uuid"
)

type GenderResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CreateResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	RoleCode string    `json:"role_code"`
	TenantID uuid.UUID `json:"tenant_id"`
}

type BranchInfo struct {
	BranchID uuid.UUID `json:"branch_id"`
	Name     string    `json:"name"`
	Code     string    `json:"code"`
}

type RoleInfo struct {
	RoleID uuid.UUID `json:"role_id"`
	Name   string    `json:"name"`
	Code   string    `json:"code"`
}

type UserDetailResponse struct {
	ID                uuid.UUID    `json:"id"`
	Email             string       `json:"email"`
	FirstName         string       `json:"first_name"`
	LastName          string       `json:"last_name"`
	FullName          string       `json:"full_name"`
	PhoneNumber       *string      `json:"phone_number,omitempty"`
	DateOfBirth       *string      `json:"date_of_birth,omitempty"`
	Gender            *GenderResponse `json:"gender,omitempty"`
	Address           *string         `json:"address,omitempty"`
	ProfilePictureURL *string         `json:"profile_picture_url,omitempty"`
	EmailVerified     bool         `json:"email_verified"`
	PINSet            bool         `json:"pin_set"`
	Status            string       `json:"status"`
	IsActive          bool         `json:"is_active"`
	LastLoginAt       *time.Time   `json:"last_login_at,omitempty"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
	Branches          []BranchInfo `json:"branches,omitempty"`
	Roles             []RoleInfo   `json:"roles,omitempty"`
}

type UserListItem struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	FullName      string     `json:"full_name"`
	PhoneNumber   *string    `json:"phone_number,omitempty"`
	EmailVerified bool       `json:"email_verified"`
	Status        string     `json:"status"`
	IsActive      bool       `json:"is_active"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type Pagination struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_pages"`
}

type ListResponse struct {
	Users      []UserListItem `json:"users"`
	Pagination Pagination     `json:"pagination"`
}

type UpdateResponse struct {
	UserDetailResponse
	Message string `json:"message"`
}

type ApproveResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}

type RejectResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}

type UnlockResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}

type ResetPINResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}

type DeleteResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}
