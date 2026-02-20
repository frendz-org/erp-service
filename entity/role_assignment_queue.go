package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type RoleAssignmentQueueStatus string

const (
	RoleAssignmentQueueStatusPending    RoleAssignmentQueueStatus = "pending"
	RoleAssignmentQueueStatusProcessing RoleAssignmentQueueStatus = "processing"
	RoleAssignmentQueueStatusCompleted  RoleAssignmentQueueStatus = "completed"
	RoleAssignmentQueueStatusFailed     RoleAssignmentQueueStatus = "failed"
	RoleAssignmentQueueStatusCancelled  RoleAssignmentQueueStatus = "cancelled"
)

type RoleAssignmentQueue struct {
	QueueID uuid.UUID `json:"queue_id" gorm:"column:queue_id;primaryKey" db:"queue_id"`

	UserID   uuid.UUID `json:"user_id" gorm:"column:user_id;not null" db:"user_id"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"column:tenant_id;not null" db:"tenant_id"`

	RoleID    uuid.UUID  `json:"role_id" gorm:"column:role_id;not null" db:"role_id"`
	ProductID *uuid.UUID `json:"product_id,omitempty" gorm:"column:product_id" db:"product_id"`
	BranchID  *uuid.UUID `json:"branch_id,omitempty" gorm:"column:branch_id" db:"branch_id"`

	EffectiveFrom time.Time  `json:"effective_from" gorm:"column:effective_from;default:CURRENT_TIMESTAMP" db:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty" gorm:"column:effective_to" db:"effective_to"`

	Status RoleAssignmentQueueStatus `json:"status" gorm:"column:status;not null;default:'pending'" db:"status"`

	AssignedBy uuid.UUID `json:"assigned_by" gorm:"column:assigned_by;not null" db:"assigned_by"`
	AssignedAt time.Time `json:"assigned_at" gorm:"column:assigned_at;not null;default:CURRENT_TIMESTAMP" db:"assigned_at"`

	BatchID       *uuid.UUID `json:"batch_id,omitempty" gorm:"column:batch_id" db:"batch_id"`
	BatchTotal    *int       `json:"batch_total,omitempty" gorm:"column:batch_total" db:"batch_total"`
	BatchSequence *int       `json:"batch_sequence,omitempty" gorm:"column:batch_sequence" db:"batch_sequence"`

	ProcessedAt         *time.Time `json:"processed_at,omitempty" gorm:"column:processed_at" db:"processed_at"`
	ProcessingStartedAt *time.Time `json:"processing_started_at,omitempty" gorm:"column:processing_started_at" db:"processing_started_at"`
	FailureReason       *string    `json:"failure_reason,omitempty" gorm:"column:failure_reason" db:"failure_reason"`
	RetryCount          int        `json:"retry_count" gorm:"column:retry_count;default:0" db:"retry_count"`

	UserRoleID *uuid.UUID `json:"user_role_id,omitempty" gorm:"column:user_role_id" db:"user_role_id"`

	Metadata json.RawMessage `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb;default:'{}'" db:"metadata"`
}

func (RoleAssignmentQueue) TableName() string {
	return "role_assignments_queue"
}

func (raq *RoleAssignmentQueue) IsPending() bool {
	return raq.Status == RoleAssignmentQueueStatusPending
}

func (raq *RoleAssignmentQueue) IsProcessing() bool {
	return raq.Status == RoleAssignmentQueueStatusProcessing
}

func (raq *RoleAssignmentQueue) IsCompleted() bool {
	return raq.Status == RoleAssignmentQueueStatusCompleted
}

func (raq *RoleAssignmentQueue) IsFailed() bool {
	return raq.Status == RoleAssignmentQueueStatusFailed
}

func (raq *RoleAssignmentQueue) IsCancelled() bool {
	return raq.Status == RoleAssignmentQueueStatusCancelled
}

func (raq *RoleAssignmentQueue) CanBeProcessed() bool {
	return raq.Status == RoleAssignmentQueueStatusPending || raq.Status == RoleAssignmentQueueStatusFailed
}

func (raq *RoleAssignmentQueue) CanBeRetried() bool {
	return raq.Status == RoleAssignmentQueueStatusFailed
}

func (raq *RoleAssignmentQueue) IsPartOfBatch() bool {
	return raq.BatchID != nil
}

func (raq *RoleAssignmentQueue) GetBatchProgress() *float64 {
	if !raq.IsPartOfBatch() || raq.BatchTotal == nil || raq.BatchSequence == nil {
		return nil
	}
	if *raq.BatchTotal == 0 {
		return nil
	}
	progress := float64(*raq.BatchSequence) / float64(*raq.BatchTotal) * 100
	return &progress
}
