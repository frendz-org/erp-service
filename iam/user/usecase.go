package user

import (
	"context"

	"github.com/google/uuid"
)

type UserReader interface {
	GetByID(ctx context.Context, callerTenantID *uuid.UUID, id uuid.UUID) (*UserDetailResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*UserDetailResponse, error)
	List(ctx context.Context, tenantID *uuid.UUID, req *ListRequest) (*ListResponse, error)
}

type UserWriter interface {
	Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error)
	Update(ctx context.Context, callerTenantID *uuid.UUID, id uuid.UUID, req *UpdateRequest) (*UserDetailResponse, error)
	UpdateMe(ctx context.Context, userID uuid.UUID, req *UpdateMeRequest) (*UserDetailResponse, error)
	Delete(ctx context.Context, callerTenantID *uuid.UUID, id uuid.UUID) error
}

type UserApproval interface {
	Approve(ctx context.Context, id uuid.UUID, approverID uuid.UUID) (*ApproveResponse, error)
	Reject(ctx context.Context, id uuid.UUID, approverID uuid.UUID, req *RejectRequest) (*RejectResponse, error)
}

type UserSecurity interface {
	Unlock(ctx context.Context, id uuid.UUID) (*UnlockResponse, error)
	ResetPIN(ctx context.Context, id uuid.UUID) (*ResetPINResponse, error)
}

type Usecase interface {
	UserReader
	UserWriter
	UserApproval
	UserSecurity
}
