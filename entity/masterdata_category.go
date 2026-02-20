package entity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type MasterdataCategoryStatus string

const (
	MasterdataCategoryStatusActive   MasterdataCategoryStatus = "ACTIVE"
	MasterdataCategoryStatusInactive MasterdataCategoryStatus = "INACTIVE"
)

type MasterdataCategory struct {
	ID                 uuid.UUID                `json:"id" gorm:"column:id;type:uuid;primaryKey;default:uuidv7()"`
	Code               string                   `json:"code" gorm:"column:code;type:varchar(50);not null;uniqueIndex"`
	Name               string                   `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Description        *string                  `json:"description,omitempty" gorm:"column:description;type:text"`
	ParentCategoryID   *uuid.UUID               `json:"parent_category_id,omitempty" gorm:"column:parent_category_id;type:uuid"`
	IsSystem           bool                     `json:"is_system" gorm:"column:is_system;not null;default:false"`
	IsTenantExtensible bool                     `json:"is_tenant_extensible" gorm:"column:is_tenant_extensible;not null;default:false"`
	SortOrder          int                      `json:"sort_order" gorm:"column:sort_order;not null;default:0"`
	Status             MasterdataCategoryStatus `json:"status" gorm:"column:status;type:varchar(20);not null;default:ACTIVE"`
	Metadata           json.RawMessage          `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb;not null;default:'{}'"`
	Version            int                      `json:"version" gorm:"column:version;not null;default:1"`

	Timestamps

	ParentCategory  *MasterdataCategory  `json:"parent_category,omitempty" gorm:"foreignKey:ParentCategoryID"`
	ChildCategories []MasterdataCategory `json:"child_categories,omitempty" gorm:"foreignKey:ParentCategoryID"`
}

func (MasterdataCategory) TableName() string {
	return "masterdata_categories"
}

func (c *MasterdataCategory) Validate() error {
	if c.Code == "" {
		return errors.New("code is required")
	}
	if len(c.Code) > 50 {
		return errors.New("code must not exceed 50 characters")
	}
	if c.Name == "" {
		return errors.New("name is required")
	}
	if len(c.Name) > 255 {
		return errors.New("name must not exceed 255 characters")
	}
	if c.Status != MasterdataCategoryStatusActive && c.Status != MasterdataCategoryStatusInactive {
		return errors.New("status must be ACTIVE or INACTIVE")
	}
	if c.ParentCategoryID != nil && *c.ParentCategoryID == c.ID {
		return errors.New("category cannot be its own parent")
	}
	return nil
}

func (c *MasterdataCategory) Deactivate() {
	c.Status = MasterdataCategoryStatusInactive
	c.UpdatedAt = time.Now()
}

func (c *MasterdataCategory) Activate() {
	c.Status = MasterdataCategoryStatusActive
	c.UpdatedAt = time.Now()
}

func (c *MasterdataCategory) IsActive() bool {
	return c.Status == MasterdataCategoryStatusActive
}

func (c *MasterdataCategory) IsHierarchical() bool {
	return c.ParentCategoryID != nil
}

func (c *MasterdataCategory) IsFlat() bool {
	return c.ParentCategoryID == nil
}

func (c *MasterdataCategory) IncrementVersion() {
	c.Version++
	c.UpdatedAt = time.Now()
}
