package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ParticipantIdentity struct {
	ID                uuid.UUID    `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	ParticipantID     uuid.UUID    `json:"participant_id" gorm:"column:participant_id;not null" db:"participant_id"`
	IdentityType      string       `json:"identity_type" gorm:"column:identity_type;not null" db:"identity_type"`
	IdentityNumber    string       `json:"identity_number" gorm:"column:identity_number;not null" db:"identity_number"`
	IdentityAuthority *string      `json:"identity_authority,omitempty" gorm:"column:identity_authority" db:"identity_authority"`
	IssueDate         *time.Time   `json:"issue_date,omitempty" gorm:"column:issue_date" db:"issue_date"`
	ExpiryDate        *time.Time   `json:"expiry_date,omitempty" gorm:"column:expiry_date" db:"expiry_date"`
	PhotoFilePath     *string      `json:"photo_file_path,omitempty" gorm:"column:photo_file_path" db:"photo_file_path"`
	Version           int          `json:"version" gorm:"column:version;not null;default:1" db:"version"`
	CreatedAt         time.Time    `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
	DeletedAt         sql.NullTime `json:"deleted_at,omitempty" gorm:"column:deleted_at" db:"deleted_at"`
}

func (ParticipantIdentity) TableName() string {
	return "participant_identities"
}
