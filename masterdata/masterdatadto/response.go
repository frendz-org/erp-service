package masterdatadto

import (
	"encoding/json"
	"time"

	"erp-service/entity"

	"github.com/google/uuid"
)

type CategoryResponse struct {
	ID                 uuid.UUID       `json:"id"`
	Code               string          `json:"code"`
	Name               string          `json:"name"`
	Description        *string         `json:"description,omitempty"`
	ParentCategoryID   *uuid.UUID      `json:"parent_category_id,omitempty"`
	IsSystem           bool            `json:"is_system"`
	IsTenantExtensible bool            `json:"is_tenant_extensible"`
	SortOrder          int             `json:"sort_order"`
	Status             string          `json:"status"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	Version            int             `json:"version"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type Pagination struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_pages"`
}

type ListCategoriesResponse struct {
	Categories []*CategoryResponse `json:"categories"`
	Pagination Pagination          `json:"pagination"`
}

type ItemResponse struct {
	ID             uuid.UUID       `json:"id"`
	CategoryID     uuid.UUID       `json:"category_id"`
	TenantID       *uuid.UUID      `json:"tenant_id,omitempty"`
	ParentItemID   *uuid.UUID      `json:"parent_item_id,omitempty"`
	Code           string          `json:"code"`
	Name           string          `json:"name"`
	AltName        *string         `json:"alt_name,omitempty"`
	Description    *string         `json:"description,omitempty"`
	SortOrder      int             `json:"sort_order"`
	IsSystem       bool            `json:"is_system"`
	IsDefault      bool            `json:"is_default"`
	Status         string          `json:"status"`
	EffectiveFrom  *time.Time      `json:"effective_from,omitempty"`
	EffectiveUntil *time.Time      `json:"effective_until,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	CreatedBy      *uuid.UUID      `json:"created_by,omitempty"`
	Version        int             `json:"version"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type ListItemsResponse struct {
	Items      []*ItemResponse `json:"items"`
	Pagination Pagination      `json:"pagination"`
}

type ValidateCodeResponse struct {
	Valid        bool   `json:"valid"`
	CategoryCode string `json:"category_code"`
	ItemCode     string `json:"item_code"`
	Message      string `json:"message,omitempty"`
}

type ValidationResult struct {
	CategoryCode string `json:"category_code"`
	ItemCode     string `json:"item_code"`
	Valid        bool   `json:"valid"`
	Message      string `json:"message,omitempty"`
}

type ValidateCodesResponse struct {
	AllValid bool               `json:"all_valid"`
	Results  []ValidationResult `json:"results"`
}

func MapCategoryToResponse(c *entity.MasterdataCategory) *CategoryResponse {
	if c == nil {
		return nil
	}
	return &CategoryResponse{
		ID:                 c.ID,
		Code:               c.Code,
		Name:               c.Name,
		Description:        c.Description,
		ParentCategoryID:   c.ParentCategoryID,
		IsSystem:           c.IsSystem,
		IsTenantExtensible: c.IsTenantExtensible,
		SortOrder:          c.SortOrder,
		Status:             string(c.Status),
		Metadata:           c.Metadata,
		Version:            c.Version,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}
}

func MapCategoriesToResponse(categories []*entity.MasterdataCategory) []*CategoryResponse {
	result := make([]*CategoryResponse, len(categories))
	for i, c := range categories {
		result[i] = MapCategoryToResponse(c)
	}
	return result
}

func MapItemToResponse(i *entity.MasterdataItem) *ItemResponse {
	if i == nil {
		return nil
	}
	return &ItemResponse{
		ID:             i.ID,
		CategoryID:     i.CategoryID,
		TenantID:       i.TenantID,
		ParentItemID:   i.ParentItemID,
		Code:           i.Code,
		Name:           i.Name,
		AltName:        i.AltName,
		Description:    i.Description,
		SortOrder:      i.SortOrder,
		IsSystem:       i.IsSystem,
		IsDefault:      i.IsDefault,
		Status:         string(i.Status),
		EffectiveFrom:  i.EffectiveFrom,
		EffectiveUntil: i.EffectiveUntil,
		Metadata:       i.Metadata,
		CreatedBy:      i.CreatedBy,
		Version:        i.Version,
		CreatedAt:      i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
	}
}

func MapItemsToResponse(items []*entity.MasterdataItem) []*ItemResponse {
	result := make([]*ItemResponse, len(items))
	for i, item := range items {
		result[i] = MapItemToResponse(item)
	}
	return result
}

func CalculateTotalPages(total int64, perPage int) int {
	if perPage <= 0 {
		return 0
	}
	pages := int(total) / perPage
	if int(total)%perPage > 0 {
		pages++
	}
	return pages
}
