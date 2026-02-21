package masterdata_test

import (
	"context"
	"time"

	"erp-service/entity"
	"erp-service/masterdata"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) List(ctx context.Context, filter *masterdata.CategoryFilter) ([]*entity.MasterdataCategory, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.MasterdataCategory), args.Get(1).(int64), args.Error(2)
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.MasterdataCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MasterdataCategory), args.Error(1)
}

func (m *MockCategoryRepository) GetByCode(ctx context.Context, code string) (*entity.MasterdataCategory, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MasterdataCategory), args.Error(1)
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *entity.MasterdataCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *entity.MasterdataCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.MasterdataCategory, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.MasterdataCategory), args.Error(1)
}

func (m *MockCategoryRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) List(ctx context.Context, filter *masterdata.ItemFilter) ([]*entity.MasterdataItem, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.MasterdataItem), args.Get(1).(int64), args.Error(2)
}

func (m *MockItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.MasterdataItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MasterdataItem), args.Error(1)
}

func (m *MockItemRepository) GetByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string) (*entity.MasterdataItem, error) {
	args := m.Called(ctx, categoryID, tenantID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MasterdataItem), args.Error(1)
}

func (m *MockItemRepository) ValidateCode(ctx context.Context, categoryCode string, itemCode string, tenantID *uuid.UUID) (bool, error) {
	args := m.Called(ctx, categoryCode, itemCode, tenantID)
	return args.Bool(0), args.Error(1)
}

func (m *MockItemRepository) Create(ctx context.Context, item *entity.MasterdataItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemRepository) Update(ctx context.Context, item *entity.MasterdataItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.MasterdataItem, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.MasterdataItem), args.Error(1)
}

func (m *MockItemRepository) GetTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID) ([]*entity.MasterdataItem, error) {
	args := m.Called(ctx, categoryCode, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.MasterdataItem), args.Error(1)
}

func (m *MockItemRepository) ListByParent(ctx context.Context, categoryCode string, parentCode string, tenantID *uuid.UUID) ([]*entity.MasterdataItem, error) {
	args := m.Called(ctx, categoryCode, parentCode, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.MasterdataItem), args.Error(1)
}

func (m *MockItemRepository) ExistsByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string) (bool, error) {
	args := m.Called(ctx, categoryID, tenantID, code)
	return args.Bool(0), args.Error(1)
}

func (m *MockItemRepository) GetDefaultItem(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID) (*entity.MasterdataItem, error) {
	args := m.Called(ctx, categoryID, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MasterdataItem), args.Error(1)
}

type MockMasterdataCache struct {
	mock.Mock
}

func (m *MockMasterdataCache) GetCategoryByID(ctx context.Context, id uuid.UUID) (*masterdata.CategoryResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*masterdata.CategoryResponse), args.Error(1)
}

func (m *MockMasterdataCache) SetCategoryByID(ctx context.Context, id uuid.UUID, category *masterdata.CategoryResponse, ttl time.Duration) error {
	args := m.Called(ctx, id, category, ttl)
	return args.Error(0)
}

func (m *MockMasterdataCache) GetCategoryByCode(ctx context.Context, code string) (*masterdata.CategoryResponse, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*masterdata.CategoryResponse), args.Error(1)
}

func (m *MockMasterdataCache) SetCategoryByCode(ctx context.Context, code string, category *masterdata.CategoryResponse, ttl time.Duration) error {
	args := m.Called(ctx, code, category, ttl)
	return args.Error(0)
}

func (m *MockMasterdataCache) GetCategoriesList(ctx context.Context, filterHash string) (*masterdata.ListCategoriesResponse, error) {
	args := m.Called(ctx, filterHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*masterdata.ListCategoriesResponse), args.Error(1)
}

func (m *MockMasterdataCache) SetCategoriesList(ctx context.Context, filterHash string, response *masterdata.ListCategoriesResponse, ttl time.Duration) error {
	args := m.Called(ctx, filterHash, response, ttl)
	return args.Error(0)
}

func (m *MockMasterdataCache) GetCategoryChildren(ctx context.Context, parentID uuid.UUID) ([]*masterdata.CategoryResponse, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*masterdata.CategoryResponse), args.Error(1)
}

func (m *MockMasterdataCache) SetCategoryChildren(ctx context.Context, parentID uuid.UUID, categories []*masterdata.CategoryResponse, ttl time.Duration) error {
	args := m.Called(ctx, parentID, categories, ttl)
	return args.Error(0)
}

func (m *MockMasterdataCache) InvalidateCategories(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMasterdataCache) GetItemByID(ctx context.Context, id uuid.UUID) (*masterdata.ItemResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*masterdata.ItemResponse), args.Error(1)
}

func (m *MockMasterdataCache) SetItemByID(ctx context.Context, id uuid.UUID, item *masterdata.ItemResponse, ttl time.Duration) error {
	args := m.Called(ctx, id, item, ttl)
	return args.Error(0)
}

func (m *MockMasterdataCache) GetItemByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string) (*masterdata.ItemResponse, error) {
	args := m.Called(ctx, categoryID, tenantID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*masterdata.ItemResponse), args.Error(1)
}

func (m *MockMasterdataCache) SetItemByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string, item *masterdata.ItemResponse, ttl time.Duration) error {
	args := m.Called(ctx, categoryID, tenantID, code, item, ttl)
	return args.Error(0)
}

func (m *MockMasterdataCache) GetItemsList(ctx context.Context, filterHash string) (*masterdata.ListItemsResponse, error) {
	args := m.Called(ctx, filterHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*masterdata.ListItemsResponse), args.Error(1)
}

func (m *MockMasterdataCache) SetItemsList(ctx context.Context, filterHash string, response *masterdata.ListItemsResponse, ttl time.Duration) error {
	args := m.Called(ctx, filterHash, response, ttl)
	return args.Error(0)
}

func (m *MockMasterdataCache) GetItemChildren(ctx context.Context, parentID uuid.UUID) ([]*masterdata.ItemResponse, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*masterdata.ItemResponse), args.Error(1)
}

func (m *MockMasterdataCache) SetItemChildren(ctx context.Context, parentID uuid.UUID, items []*masterdata.ItemResponse, ttl time.Duration) error {
	args := m.Called(ctx, parentID, items, ttl)
	return args.Error(0)
}

func (m *MockMasterdataCache) GetItemTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID) ([]*masterdata.ItemResponse, error) {
	args := m.Called(ctx, categoryCode, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*masterdata.ItemResponse), args.Error(1)
}

func (m *MockMasterdataCache) SetItemTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID, items []*masterdata.ItemResponse, ttl time.Duration) error {
	args := m.Called(ctx, categoryCode, tenantID, items, ttl)
	return args.Error(0)
}

func (m *MockMasterdataCache) GetItemDefault(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID) (*masterdata.ItemResponse, error) {
	args := m.Called(ctx, categoryID, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*masterdata.ItemResponse), args.Error(1)
}

func (m *MockMasterdataCache) SetItemDefault(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, item *masterdata.ItemResponse, ttl time.Duration) error {
	args := m.Called(ctx, categoryID, tenantID, item, ttl)
	return args.Error(0)
}

func (m *MockMasterdataCache) InvalidateItems(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
