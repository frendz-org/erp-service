package masterdata

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MasterdataCache interface {
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*CategoryResponse, error)
	SetCategoryByID(ctx context.Context, id uuid.UUID, category *CategoryResponse, ttl time.Duration) error
	GetCategoryByCode(ctx context.Context, code string) (*CategoryResponse, error)
	SetCategoryByCode(ctx context.Context, code string, category *CategoryResponse, ttl time.Duration) error
	GetCategoriesList(ctx context.Context, filterHash string) (*ListCategoriesResponse, error)
	SetCategoriesList(ctx context.Context, filterHash string, response *ListCategoriesResponse, ttl time.Duration) error
	GetCategoryChildren(ctx context.Context, parentID uuid.UUID) ([]*CategoryResponse, error)
	SetCategoryChildren(ctx context.Context, parentID uuid.UUID, categories []*CategoryResponse, ttl time.Duration) error
	InvalidateCategories(ctx context.Context) error

	GetItemByID(ctx context.Context, id uuid.UUID) (*ItemResponse, error)
	SetItemByID(ctx context.Context, id uuid.UUID, item *ItemResponse, ttl time.Duration) error
	GetItemByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string) (*ItemResponse, error)
	SetItemByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string, item *ItemResponse, ttl time.Duration) error
	GetItemsList(ctx context.Context, filterHash string) (*ListItemsResponse, error)
	SetItemsList(ctx context.Context, filterHash string, response *ListItemsResponse, ttl time.Duration) error
	GetItemChildren(ctx context.Context, parentID uuid.UUID) ([]*ItemResponse, error)
	SetItemChildren(ctx context.Context, parentID uuid.UUID, items []*ItemResponse, ttl time.Duration) error
	GetItemTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID) ([]*ItemResponse, error)
	SetItemTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID, items []*ItemResponse, ttl time.Duration) error
	GetItemDefault(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID) (*ItemResponse, error)
	SetItemDefault(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, item *ItemResponse, ttl time.Duration) error
	InvalidateItems(ctx context.Context) error
}
