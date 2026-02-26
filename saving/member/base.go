package member

import "erp-service/config"

type usecase struct {
	cfg               *config.Config
	txManager         TransactionManager
	utrRepo           UserTenantRegistrationRepository
	userRole          UserRoleRepository
	productRepo       ProductRepository
	roleRepo          RoleRepository
	configRepo        ProductRegistrationConfigRepository
	profileRepo       UserProfileRepository
	userRepo          UserRepository
	memberRepo        MemberRepository
	employeeDataRepo  EmployeeDataRepository
	tenantRepo        TenantRepository
	masterdataUsecase MasterdataUsecase
}

func NewUsecase(
	cfg *config.Config,
	txManager TransactionManager,
	utrRepo UserTenantRegistrationRepository,
	userRole UserRoleRepository,
	productRepo ProductRepository,
	roleRepo RoleRepository,
	configRepo ProductRegistrationConfigRepository,
	profileRepo UserProfileRepository,
	userRepo UserRepository,
	memberRepo MemberRepository,
	employeeDataRepo EmployeeDataRepository,
	tenantRepo TenantRepository,
	masterdataUsecase MasterdataUsecase,
) Usecase {
	return &usecase{
		cfg:               cfg,
		txManager:         txManager,
		utrRepo:           utrRepo,
		userRole:          userRole,
		productRepo:       productRepo,
		roleRepo:          roleRepo,
		configRepo:        configRepo,
		profileRepo:       profileRepo,
		userRepo:          userRepo,
		memberRepo:        memberRepo,
		employeeDataRepo:  employeeDataRepo,
		tenantRepo:        tenantRepo,
		masterdataUsecase: masterdataUsecase,
	}
}
