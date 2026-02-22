package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ParticipantStatus string

const (
	ParticipantStatusDraft           ParticipantStatus = "DRAFT"
	ParticipantStatusPendingApproval ParticipantStatus = "PENDING_APPROVAL"
	ParticipantStatusApproved        ParticipantStatus = "APPROVED"
	ParticipantStatusRejected        ParticipantStatus = "REJECTED"
)

type Participant struct {
	ID        uuid.UUID  `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	TenantID  uuid.UUID  `json:"tenant_id" gorm:"column:tenant_id;not null" db:"tenant_id"`
	ProductID uuid.UUID  `json:"product_id" gorm:"column:product_id;not null" db:"product_id"`
	UserID    *uuid.UUID `json:"user_id,omitempty" gorm:"column:user_id" db:"user_id"`

	FullName      string     `json:"full_name" gorm:"column:full_name;not null" db:"full_name"`
	Gender        *string    `json:"gender,omitempty" gorm:"column:gender" db:"gender"`
	PlaceOfBirth  *string    `json:"place_of_birth,omitempty" gorm:"column:place_of_birth" db:"place_of_birth"`
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty" gorm:"column:date_of_birth" db:"date_of_birth"`
	MaritalStatus *string    `json:"marital_status,omitempty" gorm:"column:marital_status" db:"marital_status"`
	Citizenship   *string    `json:"citizenship,omitempty" gorm:"column:citizenship" db:"citizenship"`
	Religion      *string    `json:"religion,omitempty" gorm:"column:religion" db:"religion"`

	KTPNumber      *string `json:"ktp_number,omitempty" gorm:"column:ktp_number" db:"ktp_number"`
	EmployeeNumber *string `json:"employee_number,omitempty" gorm:"column:employee_number" db:"employee_number"`
	PhoneNumber    *string `json:"phone_number,omitempty" gorm:"column:phone_number" db:"phone_number"`

	StepsCompleted map[string]bool `json:"steps_completed,omitempty" gorm:"column:steps_completed;type:jsonb;serializer:json" db:"steps_completed"`

	Status          ParticipantStatus `json:"status" gorm:"column:status;not null;default:DRAFT" db:"status"`
	CreatedBy       uuid.UUID         `json:"created_by" gorm:"column:created_by;not null" db:"created_by"`
	SubmittedBy     *uuid.UUID        `json:"submitted_by,omitempty" gorm:"column:submitted_by" db:"submitted_by"`
	SubmittedAt     *time.Time        `json:"submitted_at,omitempty" gorm:"column:submitted_at" db:"submitted_at"`
	ApprovedBy      *uuid.UUID        `json:"approved_by,omitempty" gorm:"column:approved_by" db:"approved_by"`
	ApprovedAt      *time.Time        `json:"approved_at,omitempty" gorm:"column:approved_at" db:"approved_at"`
	RejectedBy      *uuid.UUID        `json:"rejected_by,omitempty" gorm:"column:rejected_by" db:"rejected_by"`
	RejectedAt      *time.Time        `json:"rejected_at,omitempty" gorm:"column:rejected_at" db:"rejected_at"`
	RejectionReason *string           `json:"rejection_reason,omitempty" gorm:"column:rejection_reason" db:"rejection_reason"`

	Version   int          `json:"version" gorm:"column:version;not null;default:1" db:"version"`
	CreatedAt time.Time    `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at,omitempty" gorm:"column:deleted_at" db:"deleted_at"`
}

func (Participant) TableName() string {
	return "participants"
}

func (p *Participant) IsDraft() bool {
	return p.Status == ParticipantStatusDraft
}

func (p *Participant) IsPendingApproval() bool {
	return p.Status == ParticipantStatusPendingApproval
}

func (p *Participant) IsApproved() bool {
	return p.Status == ParticipantStatusApproved
}

func (p *Participant) IsRejected() bool {
	return p.Status == ParticipantStatusRejected
}

func (p *Participant) CanBeEdited() bool {
	return p.Status == ParticipantStatusDraft || p.Status == ParticipantStatusRejected
}

func (p *Participant) CanBeSubmitted() bool {
	return p.Status == ParticipantStatusDraft || p.Status == ParticipantStatusRejected
}

func (p *Participant) CanBeApproved() bool {
	return p.Status == ParticipantStatusPendingApproval
}

func (p *Participant) CanBeRejected() bool {
	return p.Status == ParticipantStatusPendingApproval
}

type ParticipantStatusHistory struct {
	ID            uuid.UUID `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	ParticipantID uuid.UUID `json:"participant_id" gorm:"column:participant_id;not null" db:"participant_id"`
	FromStatus    *string   `json:"from_status,omitempty" gorm:"column:from_status" db:"from_status"`
	ToStatus      string    `json:"to_status" gorm:"column:to_status;not null" db:"to_status"`
	ChangedBy     uuid.UUID `json:"changed_by" gorm:"column:changed_by;not null" db:"changed_by"`
	Reason        *string   `json:"reason,omitempty" gorm:"column:reason" db:"reason"`
	ChangedAt     time.Time `json:"changed_at" gorm:"column:changed_at;not null" db:"changed_at"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
}

func (ParticipantStatusHistory) TableName() string {
	return "participant_status_history"
}
