package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ParticipantPension struct {
	ID                      uuid.UUID    `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	ParticipantID           uuid.UUID    `json:"participant_id" gorm:"column:participant_id;not null" db:"participant_id"`
	ParticipantNumber       *string      `json:"participant_number,omitempty" gorm:"column:participant_number" db:"participant_number"`
	PensionCategory         *string      `json:"pension_category,omitempty" gorm:"column:pension_category" db:"pension_category"`
	PensionStatus           *string      `json:"pension_status,omitempty" gorm:"column:pension_status" db:"pension_status"`
	EffectiveDate           *time.Time   `json:"effective_date,omitempty" gorm:"column:effective_date;type:date" db:"effective_date"`
	EndDate                 *time.Time   `json:"end_date,omitempty" gorm:"column:end_date;type:date" db:"end_date"`
	ProjectedRetirementDate *time.Time   `json:"projected_retirement_date,omitempty" gorm:"column:projected_retirement_date;type:date" db:"projected_retirement_date"`
	Version                 int          `json:"version" gorm:"column:version;not null;default:1" db:"version"`
	CreatedAt               time.Time    `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt               time.Time    `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
	DeletedAt               sql.NullTime `json:"deleted_at,omitempty" gorm:"column:deleted_at" db:"deleted_at"`
}

func (ParticipantPension) TableName() string {
	return "participant_pensions"
}
