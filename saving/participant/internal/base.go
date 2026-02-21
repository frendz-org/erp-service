package internal

import (
	"erp-service/config"
	"erp-service/saving/participant/contract"
)

type usecase struct {
	cfg               *config.Config
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

	tenantRepo        contract.TenantRepository
	productRepo       contract.ProductRepository
	configRepo        contract.ProductRegistrationConfigRepository
	utrRepo           contract.UserTenantRegistrationRepository
	userProfileRepo   contract.UserProfileRepository
	masterdataUsecase contract.MasterdataUsecase
}

func NewUsecase(
	cfg *config.Config,
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
	tenantRepo contract.TenantRepository,
	productRepo contract.ProductRepository,
	configRepo contract.ProductRegistrationConfigRepository,
	utrRepo contract.UserTenantRegistrationRepository,
	userProfileRepo contract.UserProfileRepository,
	masterdataUsecase contract.MasterdataUsecase,
) contract.Usecase {
	return &usecase{
		cfg:               cfg,
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
		tenantRepo:        tenantRepo,
		productRepo:       productRepo,
		configRepo:        configRepo,
		utrRepo:           utrRepo,
		userProfileRepo:   userProfileRepo,
		masterdataUsecase: masterdataUsecase,
	}
}
