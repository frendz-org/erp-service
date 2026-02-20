package contract

import (
	"context"
	"io"

	"iam-service/masterdata/masterdatadto"
	"iam-service/saving/participant/participantdto"
)

type MasterdataValidateUsecase interface {
	ValidateItemCode(ctx context.Context, req *masterdatadto.ValidateCodeRequest) (*masterdatadto.ValidateCodeResponse, error)
}

type MasterdataUsecase interface {
	MasterdataValidateUsecase
}
type Usecase interface {
	CreateParticipant(ctx context.Context, req *participantdto.CreateParticipantRequest) (*participantdto.ParticipantResponse, error)
	UpdatePersonalData(ctx context.Context, req *participantdto.UpdatePersonalDataRequest) (*participantdto.ParticipantResponse, error)
	GetParticipant(ctx context.Context, req *participantdto.GetParticipantRequest) (*participantdto.ParticipantResponse, error)
	ListParticipants(ctx context.Context, req *participantdto.ListParticipantsRequest) (*participantdto.ListParticipantsResponse, error)
	DeleteParticipant(ctx context.Context, req *participantdto.DeleteParticipantRequest) error

	SaveIdentity(ctx context.Context, req *participantdto.SaveIdentityRequest) (*participantdto.IdentityResponse, error)
	DeleteIdentity(ctx context.Context, req *participantdto.DeleteChildEntityRequest) error

	SaveAddress(ctx context.Context, req *participantdto.SaveAddressRequest) (*participantdto.AddressResponse, error)
	DeleteAddress(ctx context.Context, req *participantdto.DeleteChildEntityRequest) error

	SaveBankAccount(ctx context.Context, req *participantdto.SaveBankAccountRequest) (*participantdto.BankAccountResponse, error)
	DeleteBankAccount(ctx context.Context, req *participantdto.DeleteChildEntityRequest) error

	SaveFamilyMember(ctx context.Context, req *participantdto.SaveFamilyMemberRequest) (*participantdto.FamilyMemberResponse, error)
	DeleteFamilyMember(ctx context.Context, req *participantdto.DeleteChildEntityRequest) error

	SaveEmployment(ctx context.Context, req *participantdto.SaveEmploymentRequest) (*participantdto.EmploymentResponse, error)

	SavePension(ctx context.Context, req *participantdto.SavePensionRequest) (*participantdto.PensionResponse, error)

	SaveBeneficiary(ctx context.Context, req *participantdto.SaveBeneficiaryRequest) (*participantdto.BeneficiaryResponse, error)
	DeleteBeneficiary(ctx context.Context, req *participantdto.DeleteChildEntityRequest) error

	UploadFile(ctx context.Context, req *participantdto.UploadFileRequest, file io.Reader, fileSize int64, contentType, filename string) (*participantdto.FileUploadResponse, error)

	SubmitParticipant(ctx context.Context, req *participantdto.SubmitParticipantRequest) (*participantdto.ParticipantResponse, error)
	ApproveParticipant(ctx context.Context, req *participantdto.ApproveParticipantRequest) (*participantdto.ParticipantResponse, error)
	RejectParticipant(ctx context.Context, req *participantdto.RejectParticipantRequest) (*participantdto.ParticipantResponse, error)

	GetStatusHistory(ctx context.Context, req *participantdto.GetParticipantRequest) ([]participantdto.StatusHistoryResponse, error)

	SelfRegister(ctx context.Context, req *participantdto.SelfRegisterRequest) (*participantdto.SelfRegisterResponse, error)
}
