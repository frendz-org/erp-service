package member

import "github.com/google/uuid"

type RegisterRequest struct {
	TenantID  uuid.UUID `json:"-"`
	ProductID uuid.UUID `json:"-"`
	UserID    uuid.UUID `json:"-"`
}

type ListRequest struct {
	TenantID  uuid.UUID `json:"-"`
	ProductID uuid.UUID `json:"-"`
	Status    *string   `json:"status"`
	Search    string    `json:"search"`
	Page      int       `json:"page" validate:"min=1"`
	PerPage   int       `json:"per_page" validate:"min=1,max=100"`
	SortBy    string    `json:"sort_by"`
	SortOrder string    `json:"sort_order"`
}

type GetMemberRequest struct {
	MemberID  uuid.UUID `json:"-"`
	TenantID  uuid.UUID `json:"-"`
	ProductID uuid.UUID `json:"-"`
}

type ApproveRequest struct {
	MemberID   uuid.UUID `json:"-"`
	TenantID   uuid.UUID `json:"-"`
	ProductID  uuid.UUID `json:"-"`
	ApproverID uuid.UUID `json:"-"`
	RoleCode   string    `json:"role_code" validate:"required"`
}

type RejectRequest struct {
	MemberID   uuid.UUID `json:"-"`
	TenantID   uuid.UUID `json:"-"`
	ProductID  uuid.UUID `json:"-"`
	ApproverID uuid.UUID `json:"-"`
	Reason     string    `json:"reason" validate:"required,min=10,max=500"`
}

type ChangeRoleRequest struct {
	MemberID  uuid.UUID `json:"-"`
	TenantID  uuid.UUID `json:"-"`
	ProductID uuid.UUID `json:"-"`
	ActorID   uuid.UUID `json:"-"`
	RoleCode  string    `json:"role_code" validate:"required"`
}

type DeactivateRequest struct {
	MemberID  uuid.UUID `json:"-"`
	TenantID  uuid.UUID `json:"-"`
	ProductID uuid.UUID `json:"-"`
	ActorID   uuid.UUID `json:"-"`
	Reason    string    `json:"reason"`
}
