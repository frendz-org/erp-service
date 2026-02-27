package participant

import (
	"context"

	"erp-service/masterdata"
)

type MasterdataValidateUsecase interface {
	ValidateItemCode(ctx context.Context, req *masterdata.ValidateCodeRequest) (*masterdata.ValidateCodeResponse, error)
}

type MasterdataUsecase interface {
	MasterdataValidateUsecase
}

type ParticipantReader interface {
	GetParticipant(ctx context.Context, req *GetParticipantRequest) (*ParticipantResponse, error)
	ListParticipants(ctx context.Context, req *ListParticipantsRequest) (*ListParticipantsResponse, error)
	GetStatusHistory(ctx context.Context, req *GetParticipantRequest) ([]StatusHistoryResponse, error)
	GetMyParticipant(ctx context.Context, req *GetMyParticipantRequest) (*MyParticipantResponse, error)
	GetMyStatusHistory(ctx context.Context, req *GetMyParticipantRequest) ([]StatusHistoryResponse, error)
}

type ParticipantWriter interface {
	CreateParticipant(ctx context.Context, req *CreateParticipantRequest) (*ParticipantResponse, error)
	UpdatePersonalData(ctx context.Context, req *UpdatePersonalDataRequest) (*ParticipantResponse, error)
	DeleteParticipant(ctx context.Context, req *DeleteParticipantRequest) error
}

type IdentityManager interface {
	SaveIdentity(ctx context.Context, req *SaveIdentityRequest) (*IdentityResponse, error)
	DeleteIdentity(ctx context.Context, req *DeleteChildEntityRequest) error
}

type AddressManager interface {
	SaveAddress(ctx context.Context, req *SaveAddressRequest) (*AddressResponse, error)
	SaveAddresses(ctx context.Context, req *SaveAddressesRequest) ([]AddressResponse, error)
	DeleteAddress(ctx context.Context, req *DeleteChildEntityRequest) error
}

type BankAccountManager interface {
	SaveBankAccount(ctx context.Context, req *SaveBankAccountRequest) (*BankAccountResponse, error)
	DeleteBankAccount(ctx context.Context, req *DeleteChildEntityRequest) error
}

type FamilyMemberManager interface {
	SaveFamilyMember(ctx context.Context, req *SaveFamilyMemberRequest) (*FamilyMemberResponse, error)
	SaveFamilyMembers(ctx context.Context, req *SaveFamilyMembersRequest) ([]FamilyMemberResponse, error)
	DeleteFamilyMember(ctx context.Context, req *DeleteChildEntityRequest) error
}

type EmploymentManager interface {
	SaveEmployment(ctx context.Context, req *SaveEmploymentRequest) (*EmploymentResponse, error)
}

type PensionManager interface {
	SavePension(ctx context.Context, req *SavePensionRequest) (*PensionResponse, error)
}

type BeneficiaryManager interface {
	SaveBeneficiary(ctx context.Context, req *SaveBeneficiaryRequest) (*BeneficiaryResponse, error)
	SaveBeneficiaries(ctx context.Context, req *SaveBeneficiariesRequest) ([]BeneficiaryResponse, error)
	DeleteBeneficiary(ctx context.Context, req *DeleteChildEntityRequest) error
}

type FileUploader interface {
	UploadFile(ctx context.Context, req *UploadFileRequest) (*FileUploadResponse, error)
}

type ParticipantWorkflow interface {
	SubmitParticipant(ctx context.Context, req *SubmitParticipantRequest) (*ParticipantResponse, error)
	ApproveParticipant(ctx context.Context, req *ApproveParticipantRequest) (*ParticipantResponse, error)
	RejectParticipant(ctx context.Context, req *RejectParticipantRequest) (*ParticipantResponse, error)
}

type ParticipantRegistration interface {
	SelfRegister(ctx context.Context, req *SelfRegisterRequest) (*SelfRegisterResponse, error)
}

type CsiAmountSummaryReader interface {
	GetCsiAmountSummary(ctx context.Context, req *CsiAmountSummaryRequest) ([]CsiAmountSummaryResponse, error)
}

type CsiLedgerHistoryReader interface {
	GetCsiLedgerHistory(ctx context.Context, req *CsiLedgerHistoryRequest) ([]CsiLedgerHistoryResponse, error)
}

type Usecase interface {
	ParticipantReader
	ParticipantWriter
	IdentityManager
	AddressManager
	BankAccountManager
	FamilyMemberManager
	EmploymentManager
	PensionManager
	BeneficiaryManager
	FileUploader
	ParticipantWorkflow
	ParticipantRegistration
	CsiAmountSummaryReader
	CsiLedgerHistoryReader
}
