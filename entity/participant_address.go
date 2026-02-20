package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ParticipantAddress struct {
	ID              uuid.UUID    `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	ParticipantID   uuid.UUID    `json:"participant_id" gorm:"column:participant_id;not null" db:"participant_id"`
	AddressType     string       `json:"address_type" gorm:"column:address_type;not null" db:"address_type"`
	CountryCode     *string      `json:"country_code,omitempty" gorm:"column:country_code" db:"country_code"`
	ProvinceCode    *string      `json:"province_code,omitempty" gorm:"column:province_code" db:"province_code"`
	CityCode        *string      `json:"city_code,omitempty" gorm:"column:city_code" db:"city_code"`
	DistrictCode    *string      `json:"district_code,omitempty" gorm:"column:district_code" db:"district_code"`
	SubdistrictCode *string      `json:"subdistrict_code,omitempty" gorm:"column:subdistrict_code" db:"subdistrict_code"`
	PostalCode      *string      `json:"postal_code,omitempty" gorm:"column:postal_code" db:"postal_code"`
	RT              *string      `json:"rt,omitempty" gorm:"column:rt" db:"rt"`
	RW              *string      `json:"rw,omitempty" gorm:"column:rw" db:"rw"`
	AddressLine     *string      `json:"address_line,omitempty" gorm:"column:address_line" db:"address_line"`
	IsPrimary       bool         `json:"is_primary" gorm:"column:is_primary;not null;default:false" db:"is_primary"`
	Version         int          `json:"version" gorm:"column:version;not null;default:1" db:"version"`
	CreatedAt       time.Time    `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
	DeletedAt       sql.NullTime `json:"deleted_at,omitempty" gorm:"column:deleted_at" db:"deleted_at"`
}

func (ParticipantAddress) TableName() string {
	return "participant_addresses"
}
