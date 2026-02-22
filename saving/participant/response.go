package participant

import (
	"time"

	"github.com/google/uuid"
)

type StepsCompleted struct {
	PersonalData  bool `json:"personal_data"`
	Address       bool `json:"address"`
	BankAccount   bool `json:"bank_account"`
	FamilyMembers bool `json:"family_members"`
	Employment    bool `json:"employment"`
	Beneficiaries bool `json:"beneficiaries"`
	Pension       bool `json:"pension"`
}

type ParticipantResponse struct {
	ID              uuid.UUID              `json:"id"`
	TenantID        uuid.UUID              `json:"tenant_id"`
	ProductID       uuid.UUID              `json:"product_id"`
	UserID          *uuid.UUID             `json:"user_id,omitempty"`
	FullName        string                 `json:"full_name"`
	Gender          *string                `json:"gender,omitempty"`
	PlaceOfBirth    *string                `json:"place_of_birth,omitempty"`
	DateOfBirth     *time.Time             `json:"date_of_birth,omitempty"`
	MaritalStatus   *string                `json:"marital_status,omitempty"`
	Citizenship     *string                `json:"citizenship,omitempty"`
	Religion        *string                `json:"religion,omitempty"`
	KTPNumber       *string                `json:"ktp_number,omitempty"`
	EmployeeNumber  *string                `json:"employee_number,omitempty"`
	PhoneNumber     *string                `json:"phone_number,omitempty"`
	Status          string                 `json:"status"`
	StepsCompleted  StepsCompleted         `json:"steps_completed"`
	CreatedBy       uuid.UUID              `json:"created_by"`
	SubmittedBy     *uuid.UUID             `json:"submitted_by,omitempty"`
	SubmittedAt     *time.Time             `json:"submitted_at,omitempty"`
	ApprovedBy      *uuid.UUID             `json:"approved_by,omitempty"`
	ApprovedAt      *time.Time             `json:"approved_at,omitempty"`
	RejectedBy      *uuid.UUID             `json:"rejected_by,omitempty"`
	RejectedAt      *time.Time             `json:"rejected_at,omitempty"`
	RejectionReason *string                `json:"rejection_reason,omitempty"`
	Version         int                    `json:"version"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Identities      []IdentityResponse     `json:"identities,omitempty"`
	Addresses       []AddressResponse      `json:"addresses,omitempty"`
	BankAccounts    []BankAccountResponse  `json:"bank_accounts,omitempty"`
	FamilyMembers   []FamilyMemberResponse `json:"family_members,omitempty"`
	Employment      *EmploymentResponse    `json:"employment,omitempty"`
	Pension         *PensionResponse       `json:"pension,omitempty"`
	Beneficiaries   []BeneficiaryResponse  `json:"beneficiaries,omitempty"`
}

type ParticipantSummaryResponse struct {
	ID             uuid.UUID  `json:"id"`
	FullName       string     `json:"full_name"`
	KTPNumber      *string    `json:"ktp_number,omitempty"`
	EmployeeNumber *string    `json:"employee_number,omitempty"`
	PhoneNumber    *string    `json:"phone_number,omitempty"`
	Status         string     `json:"status"`
	SubmittedAt    *time.Time `json:"submitted_at,omitempty"`
	ApprovedAt     *time.Time `json:"approved_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type IdentityResponse struct {
	ID                uuid.UUID  `json:"id"`
	IdentityType      string     `json:"identity_type"`
	IdentityNumber    string     `json:"identity_number"`
	IdentityAuthority *string    `json:"identity_authority,omitempty"`
	IssueDate         *time.Time `json:"issue_date,omitempty"`
	ExpiryDate        *time.Time `json:"expiry_date,omitempty"`
	PhotoFilePath     *string    `json:"photo_file_path,omitempty"`
	PhotoFileID       *uuid.UUID `json:"photo_file_id,omitempty"`
	PhotoURL          *string    `json:"photo_url,omitempty"`
	Version           int        `json:"version"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type AddressResponse struct {
	ID              uuid.UUID `json:"id"`
	AddressType     string    `json:"address_type"`
	CountryCode     *string   `json:"country_code,omitempty"`
	ProvinceCode    *string   `json:"province_code,omitempty"`
	CityCode        *string   `json:"city_code,omitempty"`
	DistrictCode    *string   `json:"district_code,omitempty"`
	SubdistrictCode *string   `json:"subdistrict_code,omitempty"`
	PostalCode      *string   `json:"postal_code,omitempty"`
	RT              *string   `json:"rt,omitempty"`
	RW              *string   `json:"rw,omitempty"`
	AddressLine     *string   `json:"address_line,omitempty"`
	IsPrimary       bool      `json:"is_primary"`
	Version         int       `json:"version"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type BankAccountResponse struct {
	ID                uuid.UUID  `json:"id"`
	BankCode          string     `json:"bank_code"`
	AccountNumber     string     `json:"account_number"`
	AccountHolderName string     `json:"account_holder_name"`
	AccountType       *string    `json:"account_type,omitempty"`
	CurrencyCode      string     `json:"currency_code"`
	IsPrimary         bool       `json:"is_primary"`
	IssueDate         *time.Time `json:"issue_date,omitempty"`
	ExpiryDate        *time.Time `json:"expiry_date,omitempty"`
	Version           int        `json:"version"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type FamilyMemberResponse struct {
	ID                    uuid.UUID  `json:"id"`
	FullName              string     `json:"full_name"`
	RelationshipType      string     `json:"relationship_type"`
	IsDependent           bool       `json:"is_dependent"`
	SupportingDocFilePath *string    `json:"supporting_doc_file_path,omitempty"`
	SupportingDocFileID   *uuid.UUID `json:"supporting_doc_file_id,omitempty"`
	SupportingDocURL      *string    `json:"supporting_doc_url,omitempty"`
	Version               int        `json:"version"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type EmploymentResponse struct {
	ID                 uuid.UUID  `json:"id"`
	PersonnelNumber    *string    `json:"personnel_number,omitempty"`
	DateOfHire         *time.Time `json:"date_of_hire,omitempty"`
	CorporateGroupName *string    `json:"corporate_group_name,omitempty"`
	LegalEntityCode    *string    `json:"legal_entity_code,omitempty"`
	LegalEntityName    *string    `json:"legal_entity_name,omitempty"`
	BusinessUnitCode   *string    `json:"business_unit_code,omitempty"`
	BusinessUnitName   *string    `json:"business_unit_name,omitempty"`
	TenantName         *string    `json:"tenant_name,omitempty"`
	EmploymentStatus   *string    `json:"employment_status,omitempty"`
	PositionName       *string    `json:"position_name,omitempty"`
	JobLevel           *string    `json:"job_level,omitempty"`
	LocationCode       *string    `json:"location_code,omitempty"`
	LocationName       *string    `json:"location_name,omitempty"`
	SubLocationName    *string    `json:"sub_location_name,omitempty"`
	RetirementDate     *time.Time `json:"retirement_date,omitempty"`
	RetirementTypeCode *string    `json:"retirement_type_code,omitempty"`
	Version            int        `json:"version"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type PensionResponse struct {
	ID                      uuid.UUID  `json:"id"`
	ParticipantNumber       *string    `json:"participant_number,omitempty"`
	PensionCategory         *string    `json:"pension_category,omitempty"`
	PensionStatus           *string    `json:"pension_status,omitempty"`
	EffectiveDate           *time.Time `json:"effective_date,omitempty"`
	EndDate                 *time.Time `json:"end_date,omitempty"`
	ProjectedRetirementDate *time.Time `json:"projected_retirement_date,omitempty"`
	Version                 int        `json:"version"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

type BeneficiaryResponse struct {
	ID                      uuid.UUID  `json:"id"`
	FamilyMemberID          uuid.UUID  `json:"family_member_id"`
	IdentityPhotoFilePath   *string    `json:"identity_photo_file_path,omitempty"`
	IdentityPhotoFileID     *uuid.UUID `json:"identity_photo_file_id,omitempty"`
	IdentityPhotoURL        *string    `json:"identity_photo_url,omitempty"`
	FamilyCardPhotoFilePath *string    `json:"family_card_photo_file_path,omitempty"`
	FamilyCardPhotoFileID   *uuid.UUID `json:"family_card_photo_file_id,omitempty"`
	FamilyCardPhotoURL      *string    `json:"family_card_photo_url,omitempty"`
	BankBookPhotoFilePath   *string    `json:"bank_book_photo_file_path,omitempty"`
	BankBookPhotoFileID     *uuid.UUID `json:"bank_book_photo_file_id,omitempty"`
	BankBookPhotoURL        *string    `json:"bank_book_photo_url,omitempty"`
	AccountNumber           *string    `json:"account_number,omitempty"`
	Version                 int        `json:"version"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

type StatusHistoryResponse struct {
	ID         uuid.UUID `json:"id"`
	FromStatus *string   `json:"from_status,omitempty"`
	ToStatus   string    `json:"to_status"`
	ChangedBy  uuid.UUID `json:"changed_by"`
	Reason     *string   `json:"reason,omitempty"`
	ChangedAt  time.Time `json:"changed_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type ListParticipantsResponse struct {
	Participants []ParticipantSummaryResponse `json:"participants"`
	Pagination   PaginationMeta               `json:"pagination"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type FileUploadResponse struct {
	FileID uuid.UUID `json:"file_id"`
}

type SelfRegisterParticipantData struct {
	ParticipantNumber string    `json:"participant_number"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

type SelfRegisterResponse struct {
	IsLinked           bool                         `json:"is_linked"`
	RegistrationStatus string                       `json:"registration_status"`
	Data               *SelfRegisterParticipantData `json:"data"`
}
