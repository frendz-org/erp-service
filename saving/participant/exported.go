package participant

import (
	"iam-service/config"
	"iam-service/saving/participant/contract"
	"iam-service/saving/participant/internal"
)

type Usecase = contract.Usecase

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
) Usecase {
	return internal.NewUsecase(
		cfg,
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
		tenantRepo,
		productRepo,
		configRepo,
		utrRepo,
		userProfileRepo,
		masterdataUsecase,
	)
}
