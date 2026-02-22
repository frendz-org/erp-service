package integration

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"erp-service/config"
	"erp-service/delivery/http/controller"
	"erp-service/delivery/http/dto/response"
	"erp-service/entity"
	"erp-service/masterdata"
	"erp-service/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

type integrationTestContext struct {
	app          *fiber.App
	cfg          *config.Config
	categoryRepo *MockCategoryRepository
	itemRepo     *MockItemRepository
	cache        *MockMasterdataCache
}

func setupIntegrationTest() *integrationTestContext {
	cfg := &config.Config{
		Masterdata: config.MasterdataConfig{
			CacheTTLCategories: 5 * time.Minute,
			CacheTTLItems:      5 * time.Minute,
			CacheTTLTree:       10 * time.Minute,
		},
	}

	categoryRepo := new(MockCategoryRepository)
	itemRepo := new(MockItemRepository)
	mockCache := new(MockMasterdataCache)

	usecase := masterdata.NewUsecase(cfg, categoryRepo, itemRepo, mockCache)

	masterdataController := controller.NewMasterdataController(cfg, usecase)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var appErr *errors.AppError
			if errors.As(err, &appErr) {
				return c.Status(appErr.HTTPStatus).JSON(response.APIResponse{
					Success: false,
					Message: appErr.Message,
					Error:   appErr.Code,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response.APIResponse{
				Success: false,
				Message: "Internal server error",
				Error:   "INTERNAL_ERROR",
			})
		},
	})

	masterdataGroup := app.Group("/api/v1/masterdata")

	masterdataGroup.Get("/categories", masterdataController.ListCategories)
	masterdataGroup.Get("/categories/code/:code", masterdataController.GetCategoryByCode)
	masterdataGroup.Get("/categories/:id", masterdataController.GetCategoryByID)
	masterdataGroup.Get("/categories/:id/children", masterdataController.GetCategoryChildren)
	masterdataGroup.Post("/categories", masterdataController.CreateCategory)
	masterdataGroup.Put("/categories/:id", masterdataController.UpdateCategory)
	masterdataGroup.Delete("/categories/:id", masterdataController.DeleteCategory)

	masterdataGroup.Get("/items", masterdataController.ListItems)
	masterdataGroup.Get("/items/tree/:categoryCode", masterdataController.GetItemTree)
	masterdataGroup.Get("/items/by-parent/:categoryCode/:parentCode", masterdataController.ListItemsByParent)
	masterdataGroup.Get("/items/default/:categoryCode", masterdataController.GetDefaultItem)
	masterdataGroup.Get("/items/code/:categoryCode/:itemCode", masterdataController.GetItemByCode)
	masterdataGroup.Get("/items/:id", masterdataController.GetItemByID)
	masterdataGroup.Get("/items/:id/children", masterdataController.GetItemChildren)
	masterdataGroup.Post("/items", masterdataController.CreateItem)
	masterdataGroup.Put("/items/:id", masterdataController.UpdateItem)
	masterdataGroup.Delete("/items/:id", masterdataController.DeleteItem)

	masterdataGroup.Post("/validate", masterdataController.ValidateCode)
	masterdataGroup.Post("/validate/batch", masterdataController.ValidateCodes)

	return &integrationTestContext{
		app:          app,
		cfg:          cfg,
		categoryRepo: categoryRepo,
		itemRepo:     itemRepo,
		cache:        mockCache,
	}
}

func TestIntegration_CategoryCRUD(t *testing.T) {
	tc := setupIntegrationTest()
	categoryID := uuid.New()

	t.Run("Create category with valid data", func(t *testing.T) {
		tc.categoryRepo.On("ExistsByCode", mock.Anything, "GENDER").Return(false, nil).Once()
		tc.categoryRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.MasterdataCategory")).Return(nil).Once()
		tc.cache.On("InvalidateCategories", mock.Anything).Return(nil).Once()

		createReq := masterdata.CreateCategoryRequest{
			Code: "GENDER",
			Name: "Gender",
		}
		body, _ := json.Marshal(createReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/masterdata/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result response.APIResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result.Success)
	})

	t.Run("Create category with duplicate code fails", func(t *testing.T) {
		tc.categoryRepo.On("ExistsByCode", mock.Anything, "GENDER").Return(true, nil).Once()

		createReq := masterdata.CreateCategoryRequest{
			Code: "GENDER",
			Name: "Gender Duplicate",
		}
		body, _ := json.Marshal(createReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/masterdata/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("Get category by ID with cache miss", func(t *testing.T) {
		category := &entity.MasterdataCategory{
			ID:     categoryID,
			Code:   "GENDER",
			Name:   "Gender",
			Status: entity.MasterdataCategoryStatusActive,
		}

		tc.cache.On("GetCategoryByID", mock.Anything, categoryID).Return(nil, stderrors.New("cache miss")).Once()
		tc.categoryRepo.On("GetByID", mock.Anything, categoryID).Return(category, nil).Once()
		tc.cache.On("SetCategoryByID", mock.Anything, categoryID, mock.Anything, mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/masterdata/categories/"+categoryID.String(), nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result response.APIResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result.Success)
	})

	t.Run("Update category with optimistic locking", func(t *testing.T) {
		category := &entity.MasterdataCategory{
			ID:      categoryID,
			Code:    "GENDER",
			Name:    "Gender",
			Status:  entity.MasterdataCategoryStatusActive,
			Version: 1,
		}

		tc.categoryRepo.On("GetByID", mock.Anything, categoryID).Return(category, nil).Once()
		tc.categoryRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.MasterdataCategory")).Return(nil).Once()
		tc.cache.On("InvalidateCategories", mock.Anything).Return(nil).Once()

		newName := "Gender Type"
		updateReq := masterdata.UpdateCategoryRequest{
			Name:    &newName,
			Version: 1,
		}
		body, _ := json.Marshal(updateReq)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/masterdata/categories/"+categoryID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Update category with wrong version fails", func(t *testing.T) {
		category := &entity.MasterdataCategory{
			ID:      categoryID,
			Code:    "GENDER",
			Name:    "Gender",
			Status:  entity.MasterdataCategoryStatusActive,
			Version: 5,
		}

		tc.categoryRepo.On("GetByID", mock.Anything, categoryID).Return(category, nil).Once()

		newName := "Gender Type"
		updateReq := masterdata.UpdateCategoryRequest{
			Name:    &newName,
			Version: 1,
		}
		body, _ := json.Marshal(updateReq)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/masterdata/categories/"+categoryID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("Delete non-system category succeeds", func(t *testing.T) {
		category := &entity.MasterdataCategory{
			ID:       categoryID,
			Code:     "CUSTOM",
			Name:     "Custom Category",
			IsSystem: false,
		}

		tc.categoryRepo.On("GetByID", mock.Anything, categoryID).Return(category, nil).Once()
		tc.categoryRepo.On("Delete", mock.Anything, categoryID).Return(nil).Once()
		tc.cache.On("InvalidateCategories", mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/masterdata/categories/"+categoryID.String(), nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete system category fails", func(t *testing.T) {
		category := &entity.MasterdataCategory{
			ID:       categoryID,
			Code:     "GENDER",
			Name:     "Gender",
			IsSystem: true,
		}

		tc.categoryRepo.On("GetByID", mock.Anything, categoryID).Return(category, nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/masterdata/categories/"+categoryID.String(), nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestIntegration_CategoryList_Pagination(t *testing.T) {
	tc := setupIntegrationTest()

	t.Run("List with default pagination", func(t *testing.T) {
		categories := []*entity.MasterdataCategory{
			{ID: uuid.New(), Code: "GENDER", Name: "Gender", Status: entity.MasterdataCategoryStatusActive},
			{ID: uuid.New(), Code: "COUNTRY", Name: "Country", Status: entity.MasterdataCategoryStatusActive},
		}

		tc.cache.On("GetCategoriesList", mock.Anything, mock.AnythingOfType("string")).Return(nil, stderrors.New("cache miss")).Once()
		tc.categoryRepo.On("List", mock.Anything, mock.Anything).Return(categories, int64(2), nil).Once()
		tc.cache.On("SetCategoriesList", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/masterdata/categories", nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result response.APIResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result.Success)
		assert.NotNil(t, result.Pagination)
		assert.Equal(t, int64(2), result.Pagination.Total)
	})

	t.Run("List with custom page and per_page", func(t *testing.T) {
		categories := []*entity.MasterdataCategory{
			{ID: uuid.New(), Code: "MARITAL", Name: "Marital Status", Status: entity.MasterdataCategoryStatusActive},
		}

		tc.cache.On("GetCategoriesList", mock.Anything, mock.AnythingOfType("string")).Return(nil, stderrors.New("cache miss")).Once()
		tc.categoryRepo.On("List", mock.Anything, mock.MatchedBy(func(f *masterdata.CategoryFilter) bool {
			return f.Page == 2 && f.PerPage == 5
		})).Return(categories, int64(6), nil).Once()
		tc.cache.On("SetCategoriesList", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/masterdata/categories?page=2&per_page=5", nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result response.APIResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result.Success)
		assert.Equal(t, int64(6), result.Pagination.Total)
		assert.Equal(t, 2, result.Pagination.TotalPages)
	})

	t.Run("List with status filter", func(t *testing.T) {
		categories := []*entity.MasterdataCategory{
			{ID: uuid.New(), Code: "GENDER", Name: "Gender", Status: entity.MasterdataCategoryStatusActive},
		}

		tc.cache.On("GetCategoriesList", mock.Anything, mock.AnythingOfType("string")).Return(nil, stderrors.New("cache miss")).Once()
		tc.categoryRepo.On("List", mock.Anything, mock.MatchedBy(func(f *masterdata.CategoryFilter) bool {
			return f.Status == "ACTIVE"
		})).Return(categories, int64(1), nil).Once()
		tc.cache.On("SetCategoriesList", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/masterdata/categories?status=ACTIVE", nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestIntegration_ItemCRUD(t *testing.T) {
	tc := setupIntegrationTest()
	categoryID := uuid.New()
	itemID := uuid.New()

	t.Run("Create item with valid data", func(t *testing.T) {
		category := &entity.MasterdataCategory{
			ID:                 categoryID,
			Code:               "GENDER",
			IsTenantExtensible: false,
		}

		tc.categoryRepo.On("GetByID", mock.Anything, categoryID).Return(category, nil).Once()
		tc.itemRepo.On("ExistsByCode", mock.Anything, categoryID, (*uuid.UUID)(nil), "MALE").Return(false, nil).Once()
		tc.itemRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.MasterdataItem")).Return(nil).Once()
		tc.cache.On("InvalidateItems", mock.Anything).Return(nil).Once()

		createReq := masterdata.CreateItemRequest{
			CategoryID: categoryID,
			Code:       "MALE",
			Name:       "Male",
		}
		body, _ := json.Marshal(createReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/masterdata/items", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("Get item by ID", func(t *testing.T) {
		item := &entity.MasterdataItem{
			ID:         itemID,
			CategoryID: categoryID,
			Code:       "MALE",
			Name:       "Male",
			Status:     entity.MasterdataItemStatusActive,
		}

		tc.cache.On("GetItemByID", mock.Anything, itemID).Return(nil, stderrors.New("cache miss")).Once()
		tc.itemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil).Once()
		tc.cache.On("SetItemByID", mock.Anything, itemID, mock.Anything, mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/masterdata/items/"+itemID.String(), nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Get item by code with category code", func(t *testing.T) {
		category := &entity.MasterdataCategory{
			ID:   categoryID,
			Code: "GENDER",
		}
		item := &entity.MasterdataItem{
			ID:         itemID,
			CategoryID: categoryID,
			Code:       "MALE",
			Name:       "Male",
			Status:     entity.MasterdataItemStatusActive,
		}

		tc.categoryRepo.On("GetByCode", mock.Anything, "GENDER").Return(category, nil).Once()
		tc.cache.On("GetItemByCode", mock.Anything, categoryID, (*uuid.UUID)(nil), "MALE").Return(nil, stderrors.New("cache miss")).Once()
		tc.itemRepo.On("GetByCode", mock.Anything, categoryID, (*uuid.UUID)(nil), "MALE").Return(item, nil).Once()
		tc.cache.On("SetItemByCode", mock.Anything, categoryID, (*uuid.UUID)(nil), "MALE", mock.Anything, mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/masterdata/items/code/GENDER/MALE", nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Update item", func(t *testing.T) {
		item := &entity.MasterdataItem{
			ID:         itemID,
			CategoryID: categoryID,
			Code:       "MALE",
			Name:       "Male",
			Status:     entity.MasterdataItemStatusActive,
			Version:    1,
		}

		tc.itemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil).Once()
		tc.itemRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.MasterdataItem")).Return(nil).Once()
		tc.cache.On("InvalidateItems", mock.Anything).Return(nil).Once()

		newName := "Male Gender"
		updateReq := masterdata.UpdateItemRequest{
			Name:    &newName,
			Version: 1,
		}
		body, _ := json.Marshal(updateReq)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/masterdata/items/"+itemID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete non-system item", func(t *testing.T) {
		item := &entity.MasterdataItem{
			ID:       itemID,
			Code:     "CUSTOM",
			Name:     "Custom",
			IsSystem: false,
		}

		tc.itemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil).Once()
		tc.itemRepo.On("Delete", mock.Anything, itemID).Return(nil).Once()
		tc.cache.On("InvalidateItems", mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/masterdata/items/"+itemID.String(), nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestIntegration_ItemTree(t *testing.T) {
	tc := setupIntegrationTest()
	categoryID := uuid.New()

	t.Run("Get item tree for hierarchical category", func(t *testing.T) {
		items := []*entity.MasterdataItem{
			{ID: uuid.New(), CategoryID: categoryID, Code: "ID", Name: "Indonesia", SortOrder: 0},
			{ID: uuid.New(), CategoryID: categoryID, Code: "WJ", Name: "West Java", SortOrder: 1},
			{ID: uuid.New(), CategoryID: categoryID, Code: "BDG", Name: "Bandung", SortOrder: 2},
		}

		tc.cache.On("GetItemTree", mock.Anything, "LOCATION", (*uuid.UUID)(nil)).Return(nil, stderrors.New("cache miss")).Once()
		tc.itemRepo.On("GetTree", mock.Anything, "LOCATION", (*uuid.UUID)(nil)).Return(items, nil).Once()
		tc.cache.On("SetItemTree", mock.Anything, "LOCATION", (*uuid.UUID)(nil), mock.Anything, mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/masterdata/items/tree/LOCATION", nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result response.APIResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result.Success)
	})

	t.Run("Get item tree with tenant filter", func(t *testing.T) {
		tenantID := uuid.New()
		items := []*entity.MasterdataItem{
			{ID: uuid.New(), CategoryID: categoryID, TenantID: &tenantID, Code: "CUSTOM1", Name: "Custom 1"},
		}

		tc.cache.On("GetItemTree", mock.Anything, "CUSTOM_CAT", &tenantID).Return(nil, stderrors.New("cache miss")).Once()
		tc.itemRepo.On("GetTree", mock.Anything, "CUSTOM_CAT", &tenantID).Return(items, nil).Once()
		tc.cache.On("SetItemTree", mock.Anything, "CUSTOM_CAT", &tenantID, mock.Anything, mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/masterdata/items/tree/CUSTOM_CAT?tenant_id="+tenantID.String(), nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestIntegration_ItemListByParent(t *testing.T) {
	tc := setupIntegrationTest()
	categoryID := uuid.New()

	t.Run("List items by parent code", func(t *testing.T) {
		items := []*entity.MasterdataItem{
			{ID: uuid.New(), CategoryID: categoryID, Code: "WJ", Name: "West Java"},
			{ID: uuid.New(), CategoryID: categoryID, Code: "EJ", Name: "East Java"},
		}

		tc.itemRepo.On("ListByParent", mock.Anything, "LOCATION", "ID", (*uuid.UUID)(nil)).Return(items, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/masterdata/items/by-parent/LOCATION/ID", nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result response.APIResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result.Success)
	})
}

func TestIntegration_GetDefaultItem(t *testing.T) {
	tc := setupIntegrationTest()
	categoryID := uuid.New()
	itemID := uuid.New()

	t.Run("Get default item for category", func(t *testing.T) {
		category := &entity.MasterdataCategory{
			ID:   categoryID,
			Code: "GENDER",
		}
		item := &entity.MasterdataItem{
			ID:         itemID,
			CategoryID: categoryID,
			Code:       "UNSPECIFIED",
			Name:       "Unspecified",
			IsDefault:  true,
		}

		tc.categoryRepo.On("GetByCode", mock.Anything, "GENDER").Return(category, nil).Once()
		tc.cache.On("GetItemDefault", mock.Anything, categoryID, (*uuid.UUID)(nil)).Return(nil, stderrors.New("cache miss")).Once()
		tc.itemRepo.On("GetDefaultItem", mock.Anything, categoryID, (*uuid.UUID)(nil)).Return(item, nil).Once()
		tc.cache.On("SetItemDefault", mock.Anything, categoryID, (*uuid.UUID)(nil), mock.Anything, mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/masterdata/items/default/GENDER", nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Get default item when none exists", func(t *testing.T) {
		category := &entity.MasterdataCategory{
			ID:   categoryID,
			Code: "NODEFAULT",
		}

		tc.categoryRepo.On("GetByCode", mock.Anything, "NODEFAULT").Return(category, nil).Once()
		tc.cache.On("GetItemDefault", mock.Anything, categoryID, (*uuid.UUID)(nil)).Return(nil, stderrors.New("cache miss")).Once()
		tc.itemRepo.On("GetDefaultItem", mock.Anything, categoryID, (*uuid.UUID)(nil)).Return(nil, errors.ErrNotFound("no default item")).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/masterdata/items/default/NODEFAULT", nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestIntegration_ValidateCode(t *testing.T) {
	tc := setupIntegrationTest()

	t.Run("Validate existing code returns valid", func(t *testing.T) {
		tc.itemRepo.On("ValidateCode", mock.Anything, "GENDER", "MALE", (*uuid.UUID)(nil)).Return(true, nil).Once()

		validateReq := masterdata.ValidateCodeRequest{
			CategoryCode: "GENDER",
			ItemCode:     "MALE",
		}
		body, _ := json.Marshal(validateReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/masterdata/validate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result response.APIResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result.Success)
	})

	t.Run("Validate non-existing code returns invalid", func(t *testing.T) {
		tc.itemRepo.On("ValidateCode", mock.Anything, "GENDER", "UNKNOWN", (*uuid.UUID)(nil)).Return(false, nil).Once()

		validateReq := masterdata.ValidateCodeRequest{
			CategoryCode: "GENDER",
			ItemCode:     "UNKNOWN",
		}
		body, _ := json.Marshal(validateReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/masterdata/validate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestIntegration_ValidateCodes_Batch(t *testing.T) {
	tc := setupIntegrationTest()

	t.Run("Batch validate multiple codes", func(t *testing.T) {
		tc.itemRepo.On("ValidateCode", mock.Anything, "GENDER", "MALE", (*uuid.UUID)(nil)).Return(true, nil).Once()
		tc.itemRepo.On("ValidateCode", mock.Anything, "GENDER", "FEMALE", (*uuid.UUID)(nil)).Return(true, nil).Once()
		tc.itemRepo.On("ValidateCode", mock.Anything, "COUNTRY", "ID", (*uuid.UUID)(nil)).Return(true, nil).Once()

		validateReq := masterdata.ValidateCodesRequest{
			Validations: []masterdata.ValidationItem{
				{CategoryCode: "GENDER", ItemCode: "MALE"},
				{CategoryCode: "GENDER", ItemCode: "FEMALE"},
				{CategoryCode: "COUNTRY", ItemCode: "ID"},
			},
		}
		body, _ := json.Marshal(validateReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/masterdata/validate/batch", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result response.APIResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result.Success)
	})

	t.Run("Batch validate with some invalid codes", func(t *testing.T) {
		tc.itemRepo.On("ValidateCode", mock.Anything, "GENDER", "MALE", (*uuid.UUID)(nil)).Return(true, nil).Once()
		tc.itemRepo.On("ValidateCode", mock.Anything, "GENDER", "INVALID", (*uuid.UUID)(nil)).Return(false, nil).Once()

		validateReq := masterdata.ValidateCodesRequest{
			Validations: []masterdata.ValidationItem{
				{CategoryCode: "GENDER", ItemCode: "MALE"},
				{CategoryCode: "GENDER", ItemCode: "INVALID"},
			},
		}
		body, _ := json.Marshal(validateReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/masterdata/validate/batch", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestIntegration_ErrorHandling(t *testing.T) {
	tc := setupIntegrationTest()

	t.Run("Invalid UUID returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/masterdata/categories/not-a-uuid", nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Not found returns 404", func(t *testing.T) {
		nonExistentID := uuid.New()

		tc.cache.On("GetCategoryByID", mock.Anything, nonExistentID).Return(nil, stderrors.New("cache miss")).Once()
		tc.categoryRepo.On("GetByID", mock.Anything, nonExistentID).Return(nil, errors.ErrNotFound("category not found")).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/masterdata/categories/"+nonExistentID.String(), nil)
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Invalid request body returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/masterdata/categories", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Validation error returns 400 with error code", func(t *testing.T) {

		createReq := masterdata.CreateCategoryRequest{}
		body, _ := json.Marshal(createReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/masterdata/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := tc.app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var result response.APIResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.False(t, result.Success)
		assert.NotEmpty(t, result.Error)
	})
}
