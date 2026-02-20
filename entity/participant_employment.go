package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ParticipantEmployment struct {
	ID                 uuid.UUID    `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	ParticipantID      uuid.UUID    `json:"participant_id" gorm:"column:participant_id;not null" db:"participant_id"`
	PersonnelNumber    *string      `json:"personnel_number,omitempty" gorm:"column:personnel_number" db:"personnel_number"`
	DateOfHire         *time.Time   `json:"date_of_hire,omitempty" gorm:"column:date_of_hire" db:"date_of_hire"`
	CorporateGroupName *string      `json:"corporate_group_name,omitempty" gorm:"column:corporate_group_name" db:"corporate_group_name"`
	LegalEntityCode    *string      `json:"legal_entity_code,omitempty" gorm:"column:legal_entity_code" db:"legal_entity_code"`
	LegalEntityName    *string      `json:"legal_entity_name,omitempty" gorm:"column:legal_entity_name" db:"legal_entity_name"`
	BusinessUnitCode   *string      `json:"business_unit_code,omitempty" gorm:"column:business_unit_code" db:"business_unit_code"`
	BusinessUnitName   *string      `json:"business_unit_name,omitempty" gorm:"column:business_unit_name" db:"business_unit_name"`
	TenantName         *string      `json:"tenant_name,omitempty" gorm:"column:tenant_name" db:"tenant_name"`
	EmploymentStatus   *string      `json:"employment_status,omitempty" gorm:"column:employment_status" db:"employment_status"`
	PositionName       *string      `json:"position_name,omitempty" gorm:"column:position_name" db:"position_name"`
	JobLevel           *string      `json:"job_level,omitempty" gorm:"column:job_level" db:"job_level"`
	LocationCode       *string      `json:"location_code,omitempty" gorm:"column:location_code" db:"location_code"`
	LocationName       *string      `json:"location_name,omitempty" gorm:"column:location_name" db:"location_name"`
	SubLocationName    *string      `json:"sub_location_name,omitempty" gorm:"column:sub_location_name" db:"sub_location_name"`
	RetirementDate     *time.Time   `json:"retirement_date,omitempty" gorm:"column:retirement_date" db:"retirement_date"`
	RetirementTypeCode *string      `json:"retirement_type_code,omitempty" gorm:"column:retirement_type_code" db:"retirement_type_code"`
	Version            int          `json:"version" gorm:"column:version;not null;default:1" db:"version"`
	CreatedAt          time.Time    `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
	DeletedAt          sql.NullTime `json:"deleted_at,omitempty" gorm:"column:deleted_at" db:"deleted_at"`
}

func (ParticipantEmployment) TableName() string {
	return "participant_employments"
}
