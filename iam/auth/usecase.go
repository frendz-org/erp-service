package auth

import (
	"context"

	"erp-service/config"
	"erp-service/pkg/logger"

	"github.com/google/uuid"
)

type Usecase interface {
	Logout(ctx context.Context, req *LogoutRequest) error
	LogoutAll(ctx context.Context, req *LogoutAllRequest) error
	RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error)

	InitiateRegistration(ctx context.Context, req *InitiateRegistrationRequest) (*InitiateRegistrationResponse, error)
	VerifyRegistrationOTP(ctx context.Context, req *VerifyRegistrationOTPRequest) (*VerifyRegistrationOTPResponse, error)
	ResendRegistrationOTP(ctx context.Context, req *ResendRegistrationOTPRequest) (*ResendRegistrationOTPResponse, error)
	CompleteRegistration(ctx context.Context, req *CompleteRegistrationRequest) (*CompleteRegistrationResponse, error)
	SetPassword(ctx context.Context, req *SetPasswordRequest) (*SetPasswordResponse, error)
	CompleteProfileRegistration(ctx context.Context, req *CompleteProfileRegistrationRequest) (*CompleteProfileRegistrationResponse, error)
	GetRegistrationStatus(ctx context.Context, registrationID uuid.UUID, email string) (*RegistrationStatusResponse, error)

	InitiateLogin(ctx context.Context, req *InitiateLoginRequest) (*UnifiedLoginResponse, error)
	VerifyLoginOTP(ctx context.Context, req *VerifyLoginOTPRequest) (*VerifyLoginOTPResponse, error)
	ResendLoginOTP(ctx context.Context, req *ResendLoginOTPRequest) (*ResendLoginOTPResponse, error)
	GetLoginStatus(ctx context.Context, req *GetLoginStatusRequest) (*LoginStatusResponse, error)
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
	refreshTokenRepo RefreshTokenRepository,
	userRoleRepo UserRoleRepository,
	productRepo ProductRepository,
	permissionRepo PermissionRepository,
	emailService EmailService,
	inMemoryStore InMemoryStore,
	userSessionRepo UserSessionRepository,
	userTenantRegRepo UserTenantRegistrationRepository,
	productsByTenantRepo ProductsByTenantRepository,
	auditLogger logger.AuditLogger,
	masterdataUsecase MasterdataUsecase,
) Usecase {
	return newUsecase(
		txManager,
		cfg,
		userRepo,
		userProfileRepo,
		userAuthMethodRepo,
		userSecurityStateRepo,
		tenantRepo,
		roleRepo,
		refreshTokenRepo,
		userRoleRepo,
		productRepo,
		permissionRepo,
		emailService,
		inMemoryStore,
		userSessionRepo,
		userTenantRegRepo,
		productsByTenantRepo,
		auditLogger,
		masterdataUsecase,
	)
}
