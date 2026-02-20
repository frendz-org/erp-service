package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ParticipantBankAccount struct {
	ID                uuid.UUID    `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	ParticipantID     uuid.UUID    `json:"participant_id" gorm:"column:participant_id;not null" db:"participant_id"`
	BankCode          string       `json:"bank_code" gorm:"column:bank_code;not null" db:"bank_code"`
	AccountNumber     string       `json:"account_number" gorm:"column:account_number;not null" db:"account_number"`
	AccountHolderName string       `json:"account_holder_name" gorm:"column:account_holder_name;not null" db:"account_holder_name"`
	AccountType       *string      `json:"account_type,omitempty" gorm:"column:account_type" db:"account_type"`
	CurrencyCode      string       `json:"currency_code" gorm:"column:currency_code;not null;default:IDR" db:"currency_code"`
	IsPrimary         bool         `json:"is_primary" gorm:"column:is_primary;not null;default:false" db:"is_primary"`
	IssueDate         *time.Time   `json:"issue_date,omitempty" gorm:"column:issue_date" db:"issue_date"`
	ExpiryDate        *time.Time   `json:"expiry_date,omitempty" gorm:"column:expiry_date" db:"expiry_date"`
	Version           int          `json:"version" gorm:"column:version;not null;default:1" db:"version"`
	CreatedAt         time.Time    `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
	DeletedAt         sql.NullTime `json:"deleted_at,omitempty" gorm:"column:deleted_at" db:"deleted_at"`
}

func (ParticipantBankAccount) TableName() string {
	return "participant_bank_accounts"
}
