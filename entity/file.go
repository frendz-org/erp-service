package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID                   uuid.UUID      `gorm:"column:id;primaryKey;type:uuid;default:uuidv7()"`
	TenantID             uuid.UUID      `gorm:"column:tenant_id;not null"`
	ProductID            uuid.UUID      `gorm:"column:product_id;not null"`
	UploadedBy           uuid.UUID      `gorm:"column:uploaded_by;not null"`
	Bucket               string         `gorm:"column:bucket;not null"       json:"-"`
	StorageKey           string         `gorm:"column:storage_key;not null"  json:"-"`
	OriginalName         string         `gorm:"column:original_name;not null"`
	ContentType          string         `gorm:"column:content_type;not null"`
	SizeBytes            int64          `gorm:"column:size_bytes;not null;default:0"`
	ExpiresAt            *time.Time     `gorm:"column:expires_at"`
	ClaimedAt            *time.Time     `gorm:"column:claimed_at"`
	FailedDeleteAttempts int            `gorm:"column:failed_delete_attempts;not null;default:0"`
	Version              int            `gorm:"column:version;not null;default:1"`
	CreatedAt            time.Time      `gorm:"column:created_at"`
	UpdatedAt            time.Time      `gorm:"column:updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (File) TableName() string {
	return "files"
}
