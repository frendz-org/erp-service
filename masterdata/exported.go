package masterdata

import (
	"context"
	"erp-service/config"
	"erp-service/masterdata/contract"
	"erp-service/masterdata/internal"
	"erp-service/masterdata/masterdatadto"

	"github.com/google/uuid"
)

type CategoryUseCase interface {
	ListCategories(ctx context.Context, req *masterdatadto.ListCategoriesRequest) (*masterdatadto.ListCategoriesResponse, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*masterdatadto.CategoryResponse, error)
	GetCategoryByCode(ctx context.Context, code string) (*masterdatadto.CategoryResponse, error)
	CreateCategory(ctx context.Context, req *masterdatadto.CreateCategoryRequest) (*masterdatadto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, req *masterdatadto.UpdateCategoryRequest) (*masterdatadto.CategoryResponse, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	GetCategoryChildren(ctx context.Context, parentID uuid.UUID) ([]*masterdatadto.CategoryResponse, error)
}

type ItemUseCase interface {
	ListItems(ctx context.Context, req *masterdatadto.ListItemsRequest) (*masterdatadto.ListItemsResponse, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (*masterdatadto.ItemResponse, error)
	GetItemByCode(ctx context.Context, categoryCode string, tenantID *uuid.UUID, itemCode string) (*masterdatadto.ItemResponse, error)
	CreateItem(ctx context.Context, req *masterdatadto.CreateItemRequest) (*masterdatadto.ItemResponse, error)
	UpdateItem(ctx context.Context, id uuid.UUID, req *masterdatadto.UpdateItemRequest) (*masterdatadto.ItemResponse, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
	GetItemChildren(ctx context.Context, parentID uuid.UUID) ([]*masterdatadto.ItemResponse, error)
	GetItemTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID) ([]*masterdatadto.ItemResponse, error)
	ListItemsByParent(ctx context.Context, categoryCode string, parentCode string, tenantID *uuid.UUID) ([]*masterdatadto.ItemResponse, error)
	GetItemDefault(ctx context.Context, categoryCode string, tenantID *uuid.UUID) (*masterdatadto.ItemResponse, error)
	ValidateItemCode(ctx context.Context, req *masterdatadto.ValidateCodeRequest) (*masterdatadto.ValidateCodeResponse, error)
	ValidateItemCodes(ctx context.Context, req *masterdatadto.ValidateCodesRequest) (*masterdatadto.ValidateCodesResponse, error)
}

type Usecase interface {
	CategoryUseCase
	ItemUseCase
}

func NewUsecase(
	cfg *config.Config,
	categoryRepo contract.CategoryRepository,
	itemRepo contract.ItemRepository,
	cache contract.MasterdataCache,
) Usecase {
	return internal.NewUsecase(cfg, categoryRepo, itemRepo, cache)
}
