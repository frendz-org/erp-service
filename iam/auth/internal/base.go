package internal

import (
	"iam-service/config"
	"iam-service/iam/auth/contract"
	"iam-service/pkg/logger"
)

type usecase struct {
	TxManager             contract.TransactionManager
	Config                *config.Config
	UserRepo              contract.UserRepository
	UserProfileRepo       contract.UserProfileRepository
	UserAuthMethodRepo    contract.UserAuthMethodRepository
	UserSecurityStateRepo contract.UserSecurityStateRepository
	TenantRepo            contract.TenantRepository
	RoleRepo              contract.RoleRepository
	RefreshTokenRepo      contract.RefreshTokenRepository
	UserRoleRepo          contract.UserRoleRepository
	ProductRepo           contract.ProductRepository
	PermissionRepo        contract.PermissionRepository
	EmailService          contract.EmailService
	InMemoryStore         contract.InMemoryStore
	UserSessionRepo       contract.UserSessionRepository
	UserTenantRegRepo     contract.UserTenantRegistrationRepository
	ProductsByTenantRepo  contract.ProductsByTenantRepository
	AuditLogger      logger.AuditLogger
	MasterdataUsecase contract.MasterdataUsecase
}

func NewUsecase(
	txManager contract.TransactionManager,
	cfg *config.Config,
	userRepo contract.UserRepository,
	userProfileRepo contract.UserProfileRepository,
	userAuthMethodRepo contract.UserAuthMethodRepository,
	userSecurityStateRepo contract.UserSecurityStateRepository,
	tenantRepo contract.TenantRepository,
	roleRepository contract.RoleRepository,
	refreshTokenRepo contract.RefreshTokenRepository,
	userRoleRepo contract.UserRoleRepository,
	productRepo contract.ProductRepository,
	permissionRepo contract.PermissionRepository,
	emailService contract.EmailService,
	inMemoryStore contract.InMemoryStore,
	userSessionRepo contract.UserSessionRepository,
	userTenantRegRepo contract.UserTenantRegistrationRepository,
	productsByTenantRepo contract.ProductsByTenantRepository,
	auditLogger logger.AuditLogger,
	masterdataUsecase contract.MasterdataUsecase,
) *usecase {
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
		AuditLogger:      auditLogger,
		MasterdataUsecase: masterdataUsecase,
	}
}
