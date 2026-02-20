package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type MFADeviceType string

const (
	MFADeviceTypeTOTP        MFADeviceType = "totp"
	MFADeviceTypeBackupCodes MFADeviceType = "backup_codes"
)

type MFADevice struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	UserID          uuid.UUID     `json:"user_id" db:"user_id"`
	DeviceType      MFADeviceType `json:"device_type" db:"device_type"`
	DeviceName      *string       `json:"device_name,omitempty" db:"device_name"`
	SecretEncrypted string        `json:"-" db:"secret_encrypted"`
	IsVerified      bool          `json:"is_verified" db:"is_verified"`
	VerifiedAt      *time.Time    `json:"verified_at,omitempty" db:"verified_at"`
	LastUsedAt      *time.Time    `json:"last_used_at,omitempty" db:"last_used_at"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	DeletedAt       sql.NullTime  `json:"deleted_at,omitempty" db:"deleted_at"`
}

func (m *MFADevice) IsTOTP() bool {
	return m.DeviceType == MFADeviceTypeTOTP
}

func (m *MFADevice) IsBackupCodes() bool {
	return m.DeviceType == MFADeviceTypeBackupCodes
}

func (m *MFADevice) IsDeleted() bool {
	return m.DeletedAt.Valid
}

func (m *MFADevice) CanBeUsed() bool {
	return m.IsVerified && !m.IsDeleted()
}

func NewMFADevice(userID uuid.UUID, deviceType MFADeviceType, deviceName *string) *MFADevice {
	return &MFADevice{
		ID:         uuid.New(),
		UserID:     userID,
		DeviceType: deviceType,
		DeviceName: deviceName,
		IsVerified: false,
		CreatedAt:  time.Now(),
	}
}

func NewTOTPDevice(userID uuid.UUID, deviceName string) *MFADevice {
	name := deviceName
	return NewMFADevice(userID, MFADeviceTypeTOTP, &name)
}

func NewBackupCodesDevice(userID uuid.UUID) *MFADevice {
	return NewMFADevice(userID, MFADeviceTypeBackupCodes, nil)
}
