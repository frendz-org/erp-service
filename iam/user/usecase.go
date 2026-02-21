package user

import (
	"context"

	"erp-service/config"

	"github.com/google/uuid"
)

type Usecase interface {
	Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error)
	GetByID(ctx context.Context, callerTenantID *uuid.UUID, id uuid.UUID) (*UserDetailResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*UserDetailResponse, error)
	UpdateMe(ctx context.Context, userID uuid.UUID, req *UpdateMeRequest) (*UserDetailResponse, error)
	List(ctx context.Context, tenantID *uuid.UUID, req *ListRequest) (*ListResponse, error)
	Update(ctx context.Context, callerTenantID *uuid.UUID, id uuid.UUID, req *UpdateRequest) (*UserDetailResponse, error)
	Delete(ctx context.Context, callerTenantID *uuid.UUID, id uuid.UUID) error
	Approve(ctx context.Context, id uuid.UUID, approverID uuid.UUID) (*ApproveResponse, error)
	Reject(ctx context.Context, id uuid.UUID, approverID uuid.UUID, req *RejectRequest) (*RejectResponse, error)
	Unlock(ctx context.Context, id uuid.UUID) (*UnlockResponse, error)
	ResetPIN(ctx context.Context, id uuid.UUID) (*ResetPINResponse, error)
}

func NewUsecase(
	txManager TransactionManager,
	cfg *config.Config,
	userRepo UserRepository,
	userProfileRepo UserProfileRepository,
	userAuthMethodRepo UserAuthMethodRepository,
	userSecurityStateRepo UserSecurityStateRepository,
	tenantRepo TenantRepository,
	roleRepo RoleRepository,
	userRoleRepo UserRoleRepository,
) Usecase {
	return newUsecase(txManager, cfg, userRepo, userProfileRepo, userAuthMethodRepo, userSecurityStateRepo, tenantRepo, roleRepo, userRoleRepo)
}
