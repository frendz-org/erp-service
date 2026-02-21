package masterdata

import (
	"context"

	"erp-service/config"

	"github.com/google/uuid"
)

type CategoryUseCase interface {
	ListCategories(ctx context.Context, req *ListCategoriesRequest) (*ListCategoriesResponse, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*CategoryResponse, error)
	GetCategoryByCode(ctx context.Context, code string) (*CategoryResponse, error)
	CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, req *UpdateCategoryRequest) (*CategoryResponse, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	GetCategoryChildren(ctx context.Context, parentID uuid.UUID) ([]*CategoryResponse, error)
}

type ItemUseCase interface {
	ListItems(ctx context.Context, req *ListItemsRequest) (*ListItemsResponse, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (*ItemResponse, error)
	GetItemByCode(ctx context.Context, categoryCode string, tenantID *uuid.UUID, itemCode string) (*ItemResponse, error)
	CreateItem(ctx context.Context, req *CreateItemRequest) (*ItemResponse, error)
	UpdateItem(ctx context.Context, id uuid.UUID, req *UpdateItemRequest) (*ItemResponse, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
	GetItemChildren(ctx context.Context, parentID uuid.UUID) ([]*ItemResponse, error)
	GetItemTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID) ([]*ItemResponse, error)
	ListItemsByParent(ctx context.Context, categoryCode string, parentCode string, tenantID *uuid.UUID) ([]*ItemResponse, error)
	GetItemDefault(ctx context.Context, categoryCode string, tenantID *uuid.UUID) (*ItemResponse, error)
	ValidateItemCode(ctx context.Context, req *ValidateCodeRequest) (*ValidateCodeResponse, error)
	ValidateItemCodes(ctx context.Context, req *ValidateCodesRequest) (*ValidateCodesResponse, error)
}

type Usecase interface {
	CategoryUseCase
	ItemUseCase
}

func NewUsecase(
	cfg *config.Config,
	categoryRepo CategoryRepository,
	itemRepo ItemRepository,
	cache MasterdataCache,
) Usecase {
	return newUsecase(cfg, categoryRepo, itemRepo, cache)
}
