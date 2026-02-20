package participantdto

import (
	"time"

	"github.com/google/uuid"
)

type CreateParticipantRequest struct {
	TenantID      uuid.UUID `json:"-"`
	ProductID uuid.UUID `json:"-"`
	UserID        uuid.UUID `json:"-"`
	FullName      string    `json:"full_name" validate:"required,min=2,max=255"`
}

type UpdatePersonalDataRequest struct {
	TenantID      uuid.UUID  `json:"-"`
	ProductID uuid.UUID  `json:"-"`
	ParticipantID uuid.UUID  `json:"-"`
	UserID        uuid.UUID  `json:"-"`
	FullName      string     `json:"full_name" validate:"required,min=2,max=255"`
	Gender        *string    `json:"gender,omitempty" validate:"omitempty,oneof=MALE FEMALE"`
	PlaceOfBirth  *string    `json:"place_of_birth,omitempty" validate:"omitempty,max=255"`
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty"`
	MaritalStatus *string    `json:"marital_status,omitempty" validate:"omitempty,max=50"`
	Citizenship   *string    `json:"citizenship,omitempty" validate:"omitempty,max=10"`
	Religion      *string    `json:"religion,omitempty" validate:"omitempty,max=50"`
	KTPNumber     *string    `json:"ktp_number,omitempty" validate:"omitempty,len=16,numeric"`
	EmployeeNumber *string   `json:"employee_number,omitempty" validate:"omitempty,max=50"`
	PhoneNumber   *string    `json:"phone_number,omitempty" validate:"omitempty,max=20"`
}

type SaveIdentityRequest struct {
	ID                *uuid.UUID `json:"id,omitempty"`
	TenantID          uuid.UUID  `json:"-"`
	ProductID     uuid.UUID  `json:"-"`
	ParticipantID     uuid.UUID  `json:"-"`
	IdentityType      string     `json:"identity_type" validate:"required,max=50"`
	IdentityNumber    string     `json:"identity_number" validate:"required,max=100"`
	IdentityAuthority *string    `json:"identity_authority,omitempty" validate:"omitempty,max=255"`
	IssueDate         *time.Time `json:"issue_date,omitempty"`
	ExpiryDate        *time.Time `json:"expiry_date,omitempty"`
	PhotoFilePath     *string    `json:"photo_file_path,omitempty" validate:"omitempty,max=500"`
}

type SaveAddressRequest struct {
	ID              *uuid.UUID `json:"id,omitempty"`
	TenantID        uuid.UUID  `json:"-"`
	ProductID   uuid.UUID  `json:"-"`
	ParticipantID   uuid.UUID  `json:"-"`
	AddressType     string     `json:"address_type" validate:"required,max=50"`
	CountryCode     *string    `json:"country_code,omitempty" validate:"omitempty,max=10"`
	ProvinceCode    *string    `json:"province_code,omitempty" validate:"omitempty,max=10"`
	CityCode        *string    `json:"city_code,omitempty" validate:"omitempty,max=10"`
	DistrictCode    *string    `json:"district_code,omitempty" validate:"omitempty,max=10"`
	SubdistrictCode *string    `json:"subdistrict_code,omitempty" validate:"omitempty,max=10"`
	PostalCode      *string    `json:"postal_code,omitempty" validate:"omitempty,max=10"`
	RT              *string    `json:"rt,omitempty" validate:"omitempty,max=5"`
	RW              *string    `json:"rw,omitempty" validate:"omitempty,max=5"`
	AddressLine     *string    `json:"address_line,omitempty" validate:"omitempty,max=500"`
	IsPrimary       bool       `json:"is_primary"`
}

type SaveBankAccountRequest struct {
	ID                *uuid.UUID `json:"id,omitempty"`
	TenantID          uuid.UUID  `json:"-"`
	ProductID     uuid.UUID  `json:"-"`
	ParticipantID     uuid.UUID  `json:"-"`
	BankCode          string     `json:"bank_code" validate:"required,max=10"`
	AccountNumber     string     `json:"account_number" validate:"required,max=50"`
	AccountHolderName string     `json:"account_holder_name" validate:"required,max=255"`
	AccountType       *string    `json:"account_type,omitempty" validate:"omitempty,max=50"`
	CurrencyCode      string     `json:"currency_code" validate:"required,len=3"`
	IsPrimary         bool       `json:"is_primary"`
	IssueDate         *time.Time `json:"issue_date,omitempty"`
	ExpiryDate        *time.Time `json:"expiry_date,omitempty"`
}

type SaveFamilyMemberRequest struct {
	ID                    *uuid.UUID `json:"id,omitempty"`
	TenantID              uuid.UUID  `json:"-"`
	ProductID         uuid.UUID  `json:"-"`
	ParticipantID         uuid.UUID  `json:"-"`
	FullName              string     `json:"full_name" validate:"required,max=255"`
	RelationshipType      string     `json:"relationship_type" validate:"required,max=50"`
	IsDependent           bool       `json:"is_dependent"`
	SupportingDocFilePath *string    `json:"supporting_doc_file_path,omitempty" validate:"omitempty,max=500"`
}

type SaveEmploymentRequest struct {
	ID                 *uuid.UUID `json:"id,omitempty"`
	TenantID           uuid.UUID  `json:"-"`
	ProductID      uuid.UUID  `json:"-"`
	ParticipantID      uuid.UUID  `json:"-"`
	PersonnelNumber    *string    `json:"personnel_number,omitempty" validate:"omitempty,max=50"`
	DateOfHire         *time.Time `json:"date_of_hire,omitempty"`
	CorporateGroupName *string    `json:"corporate_group_name,omitempty" validate:"omitempty,max=255"`
	LegalEntityCode    *string    `json:"legal_entity_code,omitempty" validate:"omitempty,max=50"`
	LegalEntityName    *string    `json:"legal_entity_name,omitempty" validate:"omitempty,max=255"`
	BusinessUnitCode   *string    `json:"business_unit_code,omitempty" validate:"omitempty,max=50"`
	BusinessUnitName   *string    `json:"business_unit_name,omitempty" validate:"omitempty,max=255"`
	TenantName         *string    `json:"tenant_name,omitempty" validate:"omitempty,max=255"`
	EmploymentStatus   *string    `json:"employment_status,omitempty" validate:"omitempty,max=50"`
	PositionName       *string    `json:"position_name,omitempty" validate:"omitempty,max=255"`
	JobLevel           *string    `json:"job_level,omitempty" validate:"omitempty,max=50"`
	LocationCode       *string    `json:"location_code,omitempty" validate:"omitempty,max=50"`
	LocationName       *string    `json:"location_name,omitempty" validate:"omitempty,max=255"`
	SubLocationName    *string    `json:"sub_location_name,omitempty" validate:"omitempty,max=255"`
	RetirementDate     *time.Time `json:"retirement_date,omitempty"`
	RetirementTypeCode *string    `json:"retirement_type_code,omitempty" validate:"omitempty,max=50"`
}

type SavePensionRequest struct {
	ID                      *uuid.UUID `json:"id,omitempty"`
	TenantID                uuid.UUID  `json:"-"`
	ProductID           uuid.UUID  `json:"-"`
	ParticipantID           uuid.UUID  `json:"-"`
	ParticipantNumber       *string    `json:"participant_number,omitempty" validate:"omitempty,max=50"`
	PensionCategory         *string    `json:"pension_category,omitempty" validate:"omitempty,max=50"`
	PensionStatus           *string    `json:"pension_status,omitempty" validate:"omitempty,max=50"`
	EffectiveDate           *time.Time `json:"effective_date,omitempty"`
	EndDate                 *time.Time `json:"end_date,omitempty"`
	ProjectedRetirementDate *time.Time `json:"projected_retirement_date,omitempty"`
}

type SaveBeneficiaryRequest struct {
	ID                      *uuid.UUID `json:"id,omitempty"`
	TenantID                uuid.UUID  `json:"-"`
	ProductID           uuid.UUID  `json:"-"`
	ParticipantID           uuid.UUID  `json:"-"`
	FamilyMemberID          uuid.UUID  `json:"family_member_id" validate:"required"`
	IdentityPhotoFilePath   *string    `json:"identity_photo_file_path,omitempty" validate:"omitempty,max=500"`
	FamilyCardPhotoFilePath *string    `json:"family_card_photo_file_path,omitempty" validate:"omitempty,max=500"`
	BankBookPhotoFilePath   *string    `json:"bank_book_photo_file_path,omitempty" validate:"omitempty,max=500"`
	AccountNumber           *string    `json:"account_number,omitempty" validate:"omitempty,max=50"`
}

type UploadFileRequest struct {
	TenantID      uuid.UUID `json:"-"`
	ProductID uuid.UUID `json:"-"`
	ParticipantID uuid.UUID `json:"-"`
	FieldName     string    `json:"-"`
}

type SubmitParticipantRequest struct {
	TenantID      uuid.UUID `json:"-"`
	ProductID uuid.UUID `json:"-"`
	ParticipantID uuid.UUID `json:"-"`
	UserID        uuid.UUID `json:"-"`
}

type ApproveParticipantRequest struct {
	TenantID      uuid.UUID `json:"-"`
	ProductID uuid.UUID `json:"-"`
	ParticipantID uuid.UUID `json:"-"`
	UserID        uuid.UUID `json:"-"`
}

type RejectParticipantRequest struct {
	TenantID      uuid.UUID `json:"-"`
	ProductID uuid.UUID `json:"-"`
	ParticipantID uuid.UUID `json:"-"`
	UserID        uuid.UUID `json:"-"`
	Reason        string    `json:"reason" validate:"required,min=10,max=500"`
}

type ListParticipantsRequest struct {
	TenantID      uuid.UUID `json:"-"`
	ProductID uuid.UUID `json:"-"`
	Status        *string    `json:"status,omitempty" validate:"omitempty,oneof=DRAFT PENDING_APPROVAL APPROVED REJECTED"`
	Search        string     `json:"search,omitempty"`
	Page          int        `json:"page" validate:"min=1"`
	PerPage       int        `json:"per_page" validate:"min=1,max=100"`
	SortBy        string     `json:"sort_by,omitempty" validate:"omitempty,oneof=created_at updated_at full_name status"`
	SortOrder     string     `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
}

type GetParticipantRequest struct {
	ParticipantID uuid.UUID `json:"-"`
	TenantID      uuid.UUID `json:"-"`
	ProductID uuid.UUID `json:"-"`
}

type DeleteParticipantRequest struct {
	ParticipantID uuid.UUID `json:"-"`
	TenantID      uuid.UUID `json:"-"`
	ProductID uuid.UUID `json:"-"`
	UserID        uuid.UUID `json:"-"`
}

type DeleteChildEntityRequest struct {
	ChildID       uuid.UUID `json:"-"`
	ParticipantID uuid.UUID `json:"-"`
	TenantID      uuid.UUID `json:"-"`
	ProductID uuid.UUID `json:"-"`
}

type SelfRegisterRequest struct {
	UserID            uuid.UUID `json:"-"`
	Organization      string    `json:"organization"       validate:"required,min=6,max=50"`
	IdentityNumber    string    `json:"identity_number"    validate:"required,len=16,numeric"`
	ParticipantNumber string    `json:"participant_number" validate:"required,min=5,max=20,alphanum"`
	PhoneNumber       string    `json:"phone_number"       validate:"required,min=9,max=16"`
}
