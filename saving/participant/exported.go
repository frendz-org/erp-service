package participant

import (
	"erp-service/config"
	"erp-service/saving/participant/contract"
	"erp-service/saving/participant/internal"

	"go.uber.org/zap"
)

type Usecase = contract.Usecase

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
) Usecase {
	return internal.NewUsecase(
		cfg,
		logger,
		txManager,
		participantRepo,
		identityRepo,
		addressRepo,
		bankAccountRepo,
		familyMemberRepo,
		employmentRepo,
		pensionRepo,
		beneficiaryRepo,
		statusHistoryRepo,
		fileStorage,
		fileRepo,
		tenantRepo,
		productRepo,
		configRepo,
		utrRepo,
		userProfileRepo,
		masterdataUsecase,
	)
}
