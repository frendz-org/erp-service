package participant

import (
	"erp-service/config"

	"go.uber.org/zap"
)

type usecase struct {
	cfg               *config.Config
	logger            *zap.Logger
	txManager         TransactionManager
	participantRepo   ParticipantRepository
	identityRepo      ParticipantIdentityRepository
	addressRepo       ParticipantAddressRepository
	bankAccountRepo   ParticipantBankAccountRepository
	familyMemberRepo  ParticipantFamilyMemberRepository
	employmentRepo    ParticipantEmploymentRepository
	pensionRepo       ParticipantPensionRepository
	beneficiaryRepo   ParticipantBeneficiaryRepository
	statusHistoryRepo ParticipantStatusHistoryRepository
	fileStorage       FileStorageAdapter
	fileRepo          FileRepository

	tenantRepo        TenantRepository
	productRepo       ProductRepository
	configRepo        ProductRegistrationConfigRepository
	utrRepo           UserTenantRegistrationRepository
	userProfileRepo   UserProfileRepository
	masterdataUsecase MasterdataUsecase
}

func NewUsecase(
	cfg *config.Config,
	logger *zap.Logger,
	txManager TransactionManager,
	participantRepo ParticipantRepository,
	identityRepo ParticipantIdentityRepository,
	addressRepo ParticipantAddressRepository,
	bankAccountRepo ParticipantBankAccountRepository,
	familyMemberRepo ParticipantFamilyMemberRepository,
	employmentRepo ParticipantEmploymentRepository,
	pensionRepo ParticipantPensionRepository,
	beneficiaryRepo ParticipantBeneficiaryRepository,
	statusHistoryRepo ParticipantStatusHistoryRepository,
	fileStorage FileStorageAdapter,
	fileRepo FileRepository,
	tenantRepo TenantRepository,
	productRepo ProductRepository,
	configRepo ProductRegistrationConfigRepository,
	utrRepo UserTenantRegistrationRepository,
	userProfileRepo UserProfileRepository,
	masterdataUsecase MasterdataUsecase,
) Usecase {
	return &usecase{
		cfg:               cfg,
		logger:            logger,
		txManager:         txManager,
		participantRepo:   participantRepo,
		identityRepo:      identityRepo,
		addressRepo:       addressRepo,
		bankAccountRepo:   bankAccountRepo,
		familyMemberRepo:  familyMemberRepo,
		employmentRepo:    employmentRepo,
		pensionRepo:       pensionRepo,
		beneficiaryRepo:   beneficiaryRepo,
		statusHistoryRepo: statusHistoryRepo,
		fileStorage:       fileStorage,
		fileRepo:          fileRepo,
		tenantRepo:        tenantRepo,
		productRepo:       productRepo,
		configRepo:        configRepo,
		utrRepo:           utrRepo,
		userProfileRepo:   userProfileRepo,
		masterdataUsecase: masterdataUsecase,
	}
}
