package auth

import (
	"erp-service/config"
	"erp-service/pkg/logger"
)

type usecase struct {
	TxManager             TransactionManager
	Config                *config.Config
	UserRepo              UserRepository
	UserProfileRepo       UserProfileRepository
	UserAuthMethodRepo    UserAuthMethodRepository
	UserSecurityStateRepo UserSecurityStateRepository
	TenantRepo            TenantRepository
	RoleRepo              RoleRepository
	RefreshTokenRepo      RefreshTokenRepository
	UserRoleRepo          UserRoleRepository
	ProductRepo           ProductRepository
	PermissionRepo        PermissionRepository
	EmailService          EmailService
	InMemoryStore         InMemoryStore
	UserSessionRepo       UserSessionRepository
	UserTenantRegRepo     UserTenantRegistrationRepository
	ProductsByTenantRepo  ProductsByTenantRepository
	AuditLogger           logger.AuditLogger
	MasterdataUsecase     MasterdataUsecase
}

func NewUsecase(
	txManager TransactionManager,
	cfg *config.Config,
	userRepo UserRepository,
	userProfileRepo UserProfileRepository,
	userAuthMethodRepo UserAuthMethodRepository,
	userSecurityStateRepo UserSecurityStateRepository,
	tenantRepo TenantRepository,
	roleRepository RoleRepository,
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
	return &usecase{
		TxManager:             txManager,
		Config:                cfg,
		UserRepo:              userRepo,
		UserProfileRepo:       userProfileRepo,
		UserAuthMethodRepo:    userAuthMethodRepo,
		UserSecurityStateRepo: userSecurityStateRepo,
		TenantRepo:            tenantRepo,
		RoleRepo:              roleRepository,
		RefreshTokenRepo:      refreshTokenRepo,
		UserRoleRepo:          userRoleRepo,
		ProductRepo:           productRepo,
		PermissionRepo:        permissionRepo,
		EmailService:          emailService,
		InMemoryStore:         inMemoryStore,
		UserSessionRepo:       userSessionRepo,
		UserTenantRegRepo:     userTenantRegRepo,
		ProductsByTenantRepo:  productsByTenantRepo,
		AuditLogger:           auditLogger,
		MasterdataUsecase:     masterdataUsecase,
	}
}
