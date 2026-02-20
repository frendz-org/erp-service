package masterdatadto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ListCategoriesRequest struct {
	Status    string     `query:"status"`
	IsSystem  *bool      `query:"is_system"`
	ParentID  *uuid.UUID `query:"parent_id"`
	Page      int        `query:"page"`
	PerPage   int        `query:"per_page"`
	SortBy    string     `query:"sort_by" validate:"omitempty,oneof=name code sort_order created_at updated_at"`
	SortOrder string     `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type CreateCategoryRequest struct {
	Code               string          `json:"code" validate:"required,max=50"`
	Name               string          `json:"name" validate:"required,max=100"`
	Description        *string         `json:"description,omitempty"`
	ParentCategoryID   *uuid.UUID      `json:"parent_category_id,omitempty"`
	IsSystem           bool            `json:"is_system"`
	IsTenantExtensible bool            `json:"is_tenant_extensible"`
	SortOrder          int             `json:"sort_order"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
}

type UpdateCategoryRequest struct {
	Name               *string         `json:"name,omitempty" validate:"omitempty,max=100"`
	Description        *string         `json:"description,omitempty"`
	ParentCategoryID   *uuid.UUID      `json:"parent_category_id,omitempty"`
	IsTenantExtensible *bool           `json:"is_tenant_extensible,omitempty"`
	SortOrder          *int            `json:"sort_order,omitempty"`
	Status             *string         `json:"status,omitempty" validate:"omitempty,oneof=ACTIVE INACTIVE"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	Version            int             `json:"version" validate:"required,min=1"`
}

type ListItemsRequest struct {
	CategoryID   *uuid.UUID `query:"category_id"`
	CategoryCode string     `query:"category_code"`
	TenantID     *uuid.UUID `query:"tenant_id"`
	ParentID     *uuid.UUID `query:"parent_id"`
	ParentCode   string     `query:"parent_code"`
	Status       string     `query:"status"`
	Search       string     `query:"search"`
	IsDefault    *bool      `query:"is_default"`
	IsSystem     *bool      `query:"is_system"`
	Page         int        `query:"page"`
	PerPage      int        `query:"per_page"`
	SortBy       string     `query:"sort_by" validate:"omitempty,oneof=name code sort_order created_at updated_at"`
	SortOrder    string     `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type CreateItemRequest struct {
	CategoryID     uuid.UUID       `json:"category_id" validate:"required"`
	TenantID       *uuid.UUID      `json:"tenant_id,omitempty"`
	ParentItemID   *uuid.UUID      `json:"parent_item_id,omitempty"`
	Code           string          `json:"code" validate:"required,max=50"`
	Name           string          `json:"name" validate:"required,max=200"`
	AltName        *string         `json:"alt_name,omitempty" validate:"omitempty,max=200"`
	Description    *string         `json:"description,omitempty"`
	SortOrder      int             `json:"sort_order"`
	IsSystem       bool            `json:"is_system"`
	IsDefault      bool            `json:"is_default"`
	EffectiveFrom  *string         `json:"effective_from,omitempty"`
	EffectiveUntil *string         `json:"effective_until,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	CreatedBy      *uuid.UUID      `json:"-"`
}

type UpdateItemRequest struct {
	ParentItemID   *uuid.UUID      `json:"parent_item_id,omitempty"`
	Name           *string         `json:"name,omitempty" validate:"omitempty,max=200"`
	AltName        *string         `json:"alt_name,omitempty" validate:"omitempty,max=200"`
	Description    *string         `json:"description,omitempty"`
	SortOrder      *int            `json:"sort_order,omitempty"`
	IsDefault      *bool           `json:"is_default,omitempty"`
	Status         *string         `json:"status,omitempty" validate:"omitempty,oneof=ACTIVE INACTIVE"`
	EffectiveFrom  *string         `json:"effective_from,omitempty"`
	EffectiveUntil *string         `json:"effective_until,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	Version        int             `json:"version" validate:"required,min=1"`
}

type ValidateCodeRequest struct {
	CategoryCode   string     `json:"category_code" validate:"required"`
	ItemCode       string     `json:"item_code" validate:"required"`
	TenantID       *uuid.UUID `json:"tenant_id,omitempty"`
	ParentItemCode string     `json:"parent_item_code,omitempty"`
	RequireActive  bool       `json:"require_active,omitempty"`
}

type ValidationItem struct {
	CategoryCode string     `json:"category_code" validate:"required"`
	ItemCode     string     `json:"item_code" validate:"required"`
	TenantID     *uuid.UUID `json:"tenant_id,omitempty"`
}

type ValidateCodesRequest struct {
	Validations []ValidationItem `json:"validations" validate:"required,min=1,max=100,dive"`
}
