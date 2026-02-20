package entity

import (
	"encoding/json"
	"net"
	"time"

	"github.com/google/uuid"
)

type AuthEventType string

const (
	AuthEventLoginSuccess  AuthEventType = "login_success"
	AuthEventLoginFailed   AuthEventType = "login_failed"
	AuthEventLogout        AuthEventType = "logout"
	AuthEventMFASuccess    AuthEventType = "mfa_success"
	AuthEventMFAFailed     AuthEventType = "mfa_failed"
	AuthEventPINSuccess    AuthEventType = "pin_success"
	AuthEventPINFailed     AuthEventType = "pin_failed"
	AuthEventTokenRefresh  AuthEventType = "token_refresh"
	AuthEventPasswordReset AuthEventType = "password_reset"
	AuthEventPINReset      AuthEventType = "pin_reset"
)

type AuthLog struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	TenantID      *uuid.UUID      `json:"tenant_id,omitempty" db:"tenant_id"`
	UserID        *uuid.UUID      `json:"user_id,omitempty" db:"user_id"`
	EventType     AuthEventType   `json:"event_type" db:"event_type"`
	Email         string          `json:"email,omitempty" db:"email"`
	IPAddress     net.IP          `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent     string          `json:"user_agent,omitempty" db:"user_agent"`
	MFAMethod     *string         `json:"mfa_method,omitempty" db:"mfa_method"`
	FailureReason *string         `json:"failure_reason,omitempty" db:"failure_reason"`
	Metadata      json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
}

type PermissionCheck struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	TenantID       *uuid.UUID      `json:"tenant_id,omitempty" db:"tenant_id"`
	UserID         *uuid.UUID      `json:"user_id,omitempty" db:"user_id"`
	PermissionCode string          `json:"permission_code" db:"permission_code"`
	ResourceID     *uuid.UUID      `json:"resource_id,omitempty" db:"resource_id"`
	ResourceType   *string         `json:"resource_type,omitempty" db:"resource_type"`
	BranchID       *uuid.UUID      `json:"branch_id,omitempty" db:"branch_id"`
	Result         bool            `json:"result" db:"result"`
	Metadata       json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
}

type AdminAction string

const (
	AdminActionCreateUser       AdminAction = "create_user"
	AdminActionUpdateUser       AdminAction = "update_user"
	AdminActionDeleteUser       AdminAction = "delete_user"
	AdminActionCreateRole       AdminAction = "create_role"
	AdminActionUpdateRole       AdminAction = "update_role"
	AdminActionDeleteRole       AdminAction = "delete_role"
	AdminActionAssignRole       AdminAction = "assign_role"
	AdminActionRevokeRole       AdminAction = "revoke_role"
	AdminActionCreateBranch     AdminAction = "create_branch"
	AdminActionUpdateBranch     AdminAction = "update_branch"
	AdminActionDeleteBranch     AdminAction = "delete_branch"
	AdminActionUpdateTenant     AdminAction = "update_tenant"
	AdminActionResetUserPIN     AdminAction = "reset_user_pin"
	AdminActionResetUserPassword AdminAction = "reset_user_password"
)

type EntityType string

const (
	EntityTypeUser       EntityType = "user"
	EntityTypeRole       EntityType = "role"
	EntityTypePermission EntityType = "permission"
	EntityTypeTenant     EntityType = "tenant"
	EntityTypeBranch     EntityType = "branch"
)

type AdminAuditLog struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	TenantID    *uuid.UUID      `json:"tenant_id,omitempty" db:"tenant_id"`
	UserID      uuid.UUID       `json:"user_id" db:"user_id"`
	Action      AdminAction     `json:"action" db:"action"`
	EntityType  EntityType      `json:"entity_type" db:"entity_type"`
	EntityID    *uuid.UUID      `json:"entity_id,omitempty" db:"entity_id"`
	BeforeState json.RawMessage `json:"before_state,omitempty" db:"before_state"`
	AfterState  json.RawMessage `json:"after_state,omitempty" db:"after_state"`
	IPAddress   net.IP          `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent   string          `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
}
