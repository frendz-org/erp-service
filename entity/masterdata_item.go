package entity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type MasterdataItemStatus string

const (
	MasterdataItemStatusActive   MasterdataItemStatus = "ACTIVE"
	MasterdataItemStatusInactive MasterdataItemStatus = "INACTIVE"
)

type MasterdataItem struct {
	ID             uuid.UUID            `json:"id" gorm:"column:id;type:uuid;primaryKey;default:uuidv7()"`
	CategoryID     uuid.UUID            `json:"category_id" gorm:"column:category_id;type:uuid;not null;index"`
	TenantID       *uuid.UUID           `json:"tenant_id,omitempty" gorm:"column:tenant_id;type:uuid;index"`
	ParentItemID   *uuid.UUID           `json:"parent_item_id,omitempty" gorm:"column:parent_item_id;type:uuid;index"`
	Code           string               `json:"code" gorm:"column:code;type:varchar(100);not null"`
	Name           string               `json:"name" gorm:"column:name;type:varchar(255);not null"`
	AltName        *string              `json:"alt_name,omitempty" gorm:"column:alt_name;type:varchar(255)"`
	Description    *string              `json:"description,omitempty" gorm:"column:description;type:text"`
	SortOrder      int                  `json:"sort_order" gorm:"column:sort_order;not null;default:0"`
	IsSystem       bool                 `json:"is_system" gorm:"column:is_system;not null;default:false"`
	IsDefault      bool                 `json:"is_default" gorm:"column:is_default;not null;default:false"`
	Status         MasterdataItemStatus `json:"status" gorm:"column:status;type:varchar(20);not null;default:ACTIVE"`
	EffectiveFrom  *time.Time           `json:"effective_from,omitempty" gorm:"column:effective_from;type:date"`
	EffectiveUntil *time.Time           `json:"effective_until,omitempty" gorm:"column:effective_until;type:date"`
	Metadata       json.RawMessage      `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb;not null;default:'{}'"`
	CreatedBy      *uuid.UUID           `json:"created_by,omitempty" gorm:"column:created_by;type:uuid"`
	Version        int                  `json:"version" gorm:"column:version;not null;default:1"`

	Timestamps

	Category   *MasterdataCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	ParentItem *MasterdataItem     `json:"parent_item,omitempty" gorm:"foreignKey:ParentItemID"`
	ChildItems []MasterdataItem    `json:"child_items,omitempty" gorm:"foreignKey:ParentItemID"`
}

func (MasterdataItem) TableName() string {
	return "masterdata_items"
}

func (i *MasterdataItem) Validate() error {
	if i.CategoryID == uuid.Nil {
		return errors.New("category_id is required")
	}
	if i.Code == "" {
		return errors.New("code is required")
	}
	if len(i.Code) > 100 {
		return errors.New("code must not exceed 100 characters")
	}
	if i.Name == "" {
		return errors.New("name is required")
	}
	if len(i.Name) > 255 {
		return errors.New("name must not exceed 255 characters")
	}
	if i.AltName != nil && len(*i.AltName) > 255 {
		return errors.New("alt_name must not exceed 255 characters")
	}
	if i.Status != MasterdataItemStatusActive && i.Status != MasterdataItemStatusInactive {
		return errors.New("status must be ACTIVE or INACTIVE")
	}
	if i.ParentItemID != nil && *i.ParentItemID == i.ID {
		return errors.New("item cannot be its own parent")
	}
	if err := i.validateEffectiveDates(); err != nil {
		return err
	}
	return nil
}

func (i *MasterdataItem) validateEffectiveDates() error {
	if i.EffectiveFrom != nil && i.EffectiveUntil != nil {
		if i.EffectiveFrom.After(*i.EffectiveUntil) {
			return errors.New("effective_from must be before or equal to effective_until")
		}
	}
	return nil
}

func (i *MasterdataItem) IsEffective(at time.Time) bool {
	if i.Status != MasterdataItemStatusActive {
		return false
	}
	if i.EffectiveFrom != nil && at.Before(*i.EffectiveFrom) {
		return false
	}
	if i.EffectiveUntil != nil && at.After(*i.EffectiveUntil) {
		return false
	}
	return true
}

func (i *MasterdataItem) IsCurrentlyEffective() bool {
	return i.IsEffective(time.Now())
}

func (i *MasterdataItem) Deactivate() {
	i.Status = MasterdataItemStatusInactive
	i.UpdatedAt = time.Now()
}

func (i *MasterdataItem) Activate() {
	i.Status = MasterdataItemStatusActive
	i.UpdatedAt = time.Now()
}

func (i *MasterdataItem) IsActive() bool {
	return i.Status == MasterdataItemStatusActive
}

func (i *MasterdataItem) IsGlobal() bool {
	return i.TenantID == nil
}

func (i *MasterdataItem) IsTenantSpecific() bool {
	return i.TenantID != nil
}

func (i *MasterdataItem) HasParent() bool {
	return i.ParentItemID != nil
}

func (i *MasterdataItem) IsRootItem() bool {
	return i.ParentItemID == nil
}

func (i *MasterdataItem) IncrementVersion() {
	i.Version++
	i.UpdatedAt = time.Now()
}
