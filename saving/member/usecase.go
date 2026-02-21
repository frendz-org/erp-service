package member

import (
	"context"

	"erp-service/config"
)

type Usecase interface {
	RegisterMember(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	ListMembers(ctx context.Context, req *ListRequest) (*ListResponse, error)
	GetMember(ctx context.Context, req *GetMemberRequest) (*MemberDetailResponse, error)
	ApproveMember(ctx context.Context, req *ApproveRequest) (*MemberDetailResponse, error)
	RejectMember(ctx context.Context, req *RejectRequest) (*MemberDetailResponse, error)
	ChangeRole(ctx context.Context, req *ChangeRoleRequest) (*MemberDetailResponse, error)
	DeactivateMember(ctx context.Context, req *DeactivateRequest) (*MemberDetailResponse, error)
}

func NewUsecase(
	cfg *config.Config,
	txManager TransactionManager,
	utrRepo UserTenantRegistrationRepository,
	userRoleRepo UserRoleRepository,
	productRepo ProductRepository,
	roleRepo RoleRepository,
	configRepo ProductRegistrationConfigRepository,
	profileRepo UserProfileRepository,
	userRepo UserRepository,
) Usecase {
	return newUsecase(
		cfg,
		txManager,
		utrRepo,
		userRoleRepo,
		productRepo,
		roleRepo,
		configRepo,
		profileRepo,
		userRepo,
	)
}
