package contract

import (
	"context"
	"time"

	"erp-service/masterdata/masterdatadto"

	"github.com/google/uuid"
)

type MasterdataCache interface {
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*masterdatadto.CategoryResponse, error)
	SetCategoryByID(ctx context.Context, id uuid.UUID, category *masterdatadto.CategoryResponse, ttl time.Duration) error
	GetCategoryByCode(ctx context.Context, code string) (*masterdatadto.CategoryResponse, error)
	SetCategoryByCode(ctx context.Context, code string, category *masterdatadto.CategoryResponse, ttl time.Duration) error
	GetCategoriesList(ctx context.Context, filterHash string) (*masterdatadto.ListCategoriesResponse, error)
	SetCategoriesList(ctx context.Context, filterHash string, response *masterdatadto.ListCategoriesResponse, ttl time.Duration) error
	GetCategoryChildren(ctx context.Context, parentID uuid.UUID) ([]*masterdatadto.CategoryResponse, error)
	SetCategoryChildren(ctx context.Context, parentID uuid.UUID, categories []*masterdatadto.CategoryResponse, ttl time.Duration) error
	InvalidateCategories(ctx context.Context) error

	GetItemByID(ctx context.Context, id uuid.UUID) (*masterdatadto.ItemResponse, error)
	SetItemByID(ctx context.Context, id uuid.UUID, item *masterdatadto.ItemResponse, ttl time.Duration) error
	GetItemByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string) (*masterdatadto.ItemResponse, error)
	SetItemByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string, item *masterdatadto.ItemResponse, ttl time.Duration) error
	GetItemsList(ctx context.Context, filterHash string) (*masterdatadto.ListItemsResponse, error)
	SetItemsList(ctx context.Context, filterHash string, response *masterdatadto.ListItemsResponse, ttl time.Duration) error
	GetItemChildren(ctx context.Context, parentID uuid.UUID) ([]*masterdatadto.ItemResponse, error)
	SetItemChildren(ctx context.Context, parentID uuid.UUID, items []*masterdatadto.ItemResponse, ttl time.Duration) error
	GetItemTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID) ([]*masterdatadto.ItemResponse, error)
	SetItemTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID, items []*masterdatadto.ItemResponse, ttl time.Duration) error
	GetItemDefault(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID) (*masterdatadto.ItemResponse, error)
	SetItemDefault(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, item *masterdatadto.ItemResponse, ttl time.Duration) error
	InvalidateItems(ctx context.Context) error
}
