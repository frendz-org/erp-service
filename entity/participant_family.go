package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ParticipantFamilyMember struct {
	ID                   uuid.UUID    `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	ParticipantID        uuid.UUID    `json:"participant_id" gorm:"column:participant_id;not null" db:"participant_id"`
	FullName             string       `json:"full_name" gorm:"column:full_name;not null" db:"full_name"`
	RelationshipType     string       `json:"relationship_type" gorm:"column:relationship_type;not null" db:"relationship_type"`
	IsDependent          bool         `json:"is_dependent" gorm:"column:is_dependent;not null;default:false" db:"is_dependent"`
	SupportingDocFilePath  *string      `json:"supporting_doc_file_path,omitempty" gorm:"column:supporting_doc_file_path" db:"supporting_doc_file_path"`
	SupportingDocFileID   *uuid.UUID   `json:"supporting_doc_file_id,omitempty" gorm:"column:supporting_doc_file_id" db:"supporting_doc_file_id"`
	Version              int          `json:"version" gorm:"column:version;not null;default:1" db:"version"`
	CreatedAt            time.Time    `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
	DeletedAt            sql.NullTime `json:"deleted_at,omitempty" gorm:"column:deleted_at" db:"deleted_at"`
}

func (ParticipantFamilyMember) TableName() string {
	return "participant_family_members"
}

type ParticipantBeneficiary struct {
	ID                     uuid.UUID    `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	ParticipantID          uuid.UUID    `json:"participant_id" gorm:"column:participant_id;not null" db:"participant_id"`
	FamilyMemberID         uuid.UUID    `json:"family_member_id" gorm:"column:family_member_id;not null" db:"family_member_id"`
	IdentityPhotoFilePath   *string      `json:"identity_photo_file_path,omitempty" gorm:"column:identity_photo_file_path" db:"identity_photo_file_path"`
	IdentityPhotoFileID     *uuid.UUID   `json:"identity_photo_file_id,omitempty" gorm:"column:identity_photo_file_id" db:"identity_photo_file_id"`
	FamilyCardPhotoFilePath *string      `json:"family_card_photo_file_path,omitempty" gorm:"column:family_card_photo_file_path" db:"family_card_photo_file_path"`
	FamilyCardPhotoFileID   *uuid.UUID   `json:"family_card_photo_file_id,omitempty" gorm:"column:family_card_photo_file_id" db:"family_card_photo_file_id"`
	BankBookPhotoFilePath   *string      `json:"bank_book_photo_file_path,omitempty" gorm:"column:bank_book_photo_file_path" db:"bank_book_photo_file_path"`
	BankBookPhotoFileID     *uuid.UUID   `json:"bank_book_photo_file_id,omitempty" gorm:"column:bank_book_photo_file_id" db:"bank_book_photo_file_id"`
	AccountNumber           *string      `json:"account_number,omitempty" gorm:"column:account_number" db:"account_number"`
	Version                int          `json:"version" gorm:"column:version;not null;default:1" db:"version"`
	CreatedAt              time.Time    `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt              time.Time    `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
	DeletedAt              sql.NullTime `json:"deleted_at,omitempty" gorm:"column:deleted_at" db:"deleted_at"`
}

func (ParticipantBeneficiary) TableName() string {
	return "participant_beneficiaries"
}
