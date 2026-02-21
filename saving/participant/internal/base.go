package internal

import (
	"erp-service/config"
	"erp-service/saving/participant/contract"

	"go.uber.org/zap"
)

type usecase struct {
	cfg               *config.Config
	logger            *zap.Logger
	txManager         contract.TransactionManager
	participantRepo   contract.ParticipantRepository
	identityRepo      contract.ParticipantIdentityRepository
	addressRepo       contract.ParticipantAddressRepository
	bankAccountRepo   contract.ParticipantBankAccountRepository
	familyMemberRepo  contract.ParticipantFamilyMemberRepository
	employmentRepo    contract.ParticipantEmploymentRepository
	pensionRepo       contract.ParticipantPensionRepository
	beneficiaryRepo   contract.ParticipantBeneficiaryRepository
	statusHistoryRepo contract.ParticipantStatusHistoryRepository
	fileStorage       contract.FileStorageAdapter
	fileRepo          contract.FileRepository

	tenantRepo        contract.TenantRepository
	productRepo       contract.ProductRepository
	configRepo        contract.ProductRegistrationConfigRepository
	utrRepo           contract.UserTenantRegistrationRepository
	userProfileRepo   contract.UserProfileRepository
	masterdataUsecase contract.MasterdataUsecase
}

func NewUsecase(
	cfg *config.Config,
	logger *zap.Logger,
	txManager contract.TransactionManager,
	participantRepo contract.ParticipantRepository,
	identityRepo contract.ParticipantIdentityRepository,
	addressRepo contract.ParticipantAddressRepository,
	bankAccountRepo contract.ParticipantBankAccountRepository,
	familyMemberRepo contract.ParticipantFamilyMemberRepository,
	employmentRepo contract.ParticipantEmploymentRepository,
	pensionRepo contract.ParticipantPensionRepository,
	beneficiaryRepo contract.ParticipantBeneficiaryRepository,
	statusHistoryRepo contract.ParticipantStatusHistoryRepository,
	fileStorage contract.FileStorageAdapter,
	fileRepo contract.FileRepository,
	tenantRepo contract.TenantRepository,
	productRepo contract.ProductRepository,
	configRepo contract.ProductRegistrationConfigRepository,
	utrRepo contract.UserTenantRegistrationRepository,
	userProfileRepo contract.UserProfileRepository,
	masterdataUsecase contract.MasterdataUsecase,
) contract.Usecase {
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
