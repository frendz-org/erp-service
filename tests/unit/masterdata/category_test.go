package masterdata_test

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"erp-service/config"
	"erp-service/entity"
	"erp-service/masterdata"
	pkgerrors "erp-service/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupCategoryTest() (masterdata.Usecase, *MockCategoryRepository, *MockItemRepository, *MockMasterdataCache) {
	categoryRepo := new(MockCategoryRepository)
	itemRepo := new(MockItemRepository)
	cache := new(MockMasterdataCache)

	cfg := &config.Config{
		Masterdata: config.MasterdataConfig{
			CacheTTLCategories: 5 * time.Minute,
			CacheTTLItems:      5 * time.Minute,
			CacheTTLTree:       10 * time.Minute,
		},
	}

	uc := masterdata.NewUsecase(cfg, categoryRepo, itemRepo, cache)

	return uc, categoryRepo, itemRepo, cache
}

func TestCategoryList_Success(t *testing.T) {
	uc, categoryRepo, _, cache := setupCategoryTest()
	ctx := context.Background()

	req := &masterdata.ListCategoriesRequest{
		Page:    1,
		PerPage: 10,
	}

	expectedCategories := []*entity.MasterdataCategory{
		{
			ID:     uuid.New(),
			Code:   "GENDER",
			Name:   "Gender",
			Status: entity.MasterdataCategoryStatusActive,
		},
		{
			ID:     uuid.New(),
			Code:   "COUNTRY",
			Name:   "Country",
			Status: entity.MasterdataCategoryStatusActive,
		},
	}

	cache.On("GetCategoriesList", ctx, mock.AnythingOfType("string")).Return(nil, stderrors.New("cache miss"))

	categoryRepo.On("List", ctx, mock.AnythingOfType("*masterdata.CategoryFilter")).
		Return(expectedCategories, int64(2), nil)

	cache.On("SetCategoriesList", ctx, mock.AnythingOfType("string"), mock.Anything, 5*time.Minute).Return(nil)

	result, err := uc.ListCategories(ctx, req)

	require.NoError(t, err)
	assert.Len(t, result.Categories, 2)
	assert.Equal(t, int64(2), result.Pagination.Total)
	assert.Equal(t, 1, result.Pagination.Page)
	assert.Equal(t, 10, result.Pagination.PerPage)

	categoryRepo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestCategoryList_CacheHit(t *testing.T) {
	uc, _, _, cache := setupCategoryTest()
	ctx := context.Background()

	req := &masterdata.ListCategoriesRequest{
		Page:    1,
		PerPage: 10,
	}

	cachedResponse := &masterdata.ListCategoriesResponse{
		Categories: []*masterdata.CategoryResponse{
			{ID: uuid.New(), Code: "CACHED", Name: "Cached Category"},
		},
		Pagination: masterdata.Pagination{Total: 1, Page: 1, PerPage: 10},
	}
	cache.On("GetCategoriesList", ctx, mock.AnythingOfType("string")).Return(cachedResponse, nil)

	result, err := uc.ListCategories(ctx, req)

	require.NoError(t, err)
	assert.Len(t, result.Categories, 1)
	assert.Equal(t, "CACHED", result.Categories[0].Code)
}

func TestCategoryGetByID_Success(t *testing.T) {
	uc, categoryRepo, _, cache := setupCategoryTest()
	ctx := context.Background()

	categoryID := uuid.New()
	expectedCategory := &entity.MasterdataCategory{
		ID:     categoryID,
		Code:   "GENDER",
		Name:   "Gender",
		Status: entity.MasterdataCategoryStatusActive,
	}

	cache.On("GetCategoryByID", ctx, categoryID).Return(nil, stderrors.New("cache miss"))

	categoryRepo.On("GetByID", ctx, categoryID).Return(expectedCategory, nil)

	cache.On("SetCategoryByID", ctx, categoryID, mock.Anything, 5*time.Minute).Return(nil)

	result, err := uc.GetCategoryByID(ctx, categoryID)

	require.NoError(t, err)
	assert.Equal(t, categoryID, result.ID)
	assert.Equal(t, "GENDER", result.Code)

	categoryRepo.AssertExpectations(t)
}

func TestCategoryGetByID_NotFound(t *testing.T) {
	uc, categoryRepo, _, cache := setupCategoryTest()
	ctx := context.Background()

	categoryID := uuid.New()

	cache.On("GetCategoryByID", ctx, categoryID).Return(nil, stderrors.New("cache miss"))

	categoryRepo.On("GetByID", ctx, categoryID).Return(nil, pkgerrors.ErrNotFound("not found"))

	result, err := uc.GetCategoryByID(ctx, categoryID)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.True(t, pkgerrors.IsNotFound(err))
}

func TestCategoryGetByCode_Success(t *testing.T) {
	uc, categoryRepo, _, cache := setupCategoryTest()
	ctx := context.Background()

	categoryID := uuid.New()
	expectedCategory := &entity.MasterdataCategory{
		ID:     categoryID,
		Code:   "GENDER",
		Name:   "Gender",
		Status: entity.MasterdataCategoryStatusActive,
	}

	cache.On("GetCategoryByCode", ctx, "GENDER").Return(nil, stderrors.New("cache miss"))

	categoryRepo.On("GetByCode", ctx, "GENDER").Return(expectedCategory, nil)

	cache.On("SetCategoryByCode", ctx, "GENDER", mock.Anything, 5*time.Minute).Return(nil)

	result, err := uc.GetCategoryByCode(ctx, "GENDER")

	require.NoError(t, err)
	assert.Equal(t, "GENDER", result.Code)
}

func TestCategoryCreate_Success(t *testing.T) {
	uc, categoryRepo, _, cache := setupCategoryTest()
	ctx := context.Background()

	req := &masterdata.CreateCategoryRequest{
		Code:        "NATIONALITY",
		Name:        "Nationality",
		Description: ptrString("User nationality"),
	}

	categoryRepo.On("ExistsByCode", ctx, "NATIONALITY").Return(false, nil)

	categoryRepo.On("Create", ctx, mock.AnythingOfType("*entity.MasterdataCategory")).Return(nil)

	cache.On("InvalidateCategories", ctx).Return(nil)

	result, err := uc.CreateCategory(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, "NATIONALITY", result.Code)
	assert.Equal(t, "Nationality", result.Name)
	assert.Equal(t, string(entity.MasterdataCategoryStatusActive), result.Status)

	categoryRepo.AssertExpectations(t)
}

func TestCategoryCreate_CodeExists(t *testing.T) {
	uc, categoryRepo, _, _ := setupCategoryTest()
	ctx := context.Background()

	req := &masterdata.CreateCategoryRequest{
		Code: "GENDER",
		Name: "Gender",
	}

	categoryRepo.On("ExistsByCode", ctx, "GENDER").Return(true, nil)

	result, err := uc.CreateCategory(ctx, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.True(t, pkgerrors.IsConflict(err))
}

func TestCategoryCreate_WithParent(t *testing.T) {
	uc, categoryRepo, _, cache := setupCategoryTest()
	ctx := context.Background()

	parentID := uuid.New()
	parentCategory := &entity.MasterdataCategory{
		ID:   parentID,
		Code: "LOCATION",
		Name: "Location",
	}

	req := &masterdata.CreateCategoryRequest{
		Code:             "COUNTRY",
		Name:             "Country",
		ParentCategoryID: &parentID,
	}

	categoryRepo.On("ExistsByCode", ctx, "COUNTRY").Return(false, nil)

	categoryRepo.On("GetByID", ctx, parentID).Return(parentCategory, nil)

	categoryRepo.On("Create", ctx, mock.AnythingOfType("*entity.MasterdataCategory")).Return(nil)

	cache.On("InvalidateCategories", ctx).Return(nil)

	result, err := uc.CreateCategory(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, "COUNTRY", result.Code)
}

func TestCategoryUpdate_Success(t *testing.T) {
	uc, categoryRepo, _, cache := setupCategoryTest()
	ctx := context.Background()

	categoryID := uuid.New()
	existingCategory := &entity.MasterdataCategory{
		ID:      categoryID,
		Code:    "GENDER",
		Name:    "Gender",
		Status:  entity.MasterdataCategoryStatusActive,
		Version: 1,
	}

	newName := "Gender Type"
	req := &masterdata.UpdateCategoryRequest{
		Name:    &newName,
		Version: 1,
	}

	categoryRepo.On("GetByID", ctx, categoryID).Return(existingCategory, nil)

	categoryRepo.On("Update", ctx, mock.AnythingOfType("*entity.MasterdataCategory")).Return(nil)

	cache.On("InvalidateCategories", ctx).Return(nil)

	result, err := uc.UpdateCategory(ctx, categoryID, req)

	require.NoError(t, err)
	assert.Equal(t, "Gender Type", result.Name)
	assert.Equal(t, 2, result.Version)
}

func TestCategoryUpdate_VersionConflict(t *testing.T) {
	uc, categoryRepo, _, _ := setupCategoryTest()
	ctx := context.Background()

	categoryID := uuid.New()
	existingCategory := &entity.MasterdataCategory{
		ID:      categoryID,
		Code:    "GENDER",
		Name:    "Gender",
		Version: 2,
	}

	newName := "Gender Type"
	req := &masterdata.UpdateCategoryRequest{
		Name:    &newName,
		Version: 1,
	}

	categoryRepo.On("GetByID", ctx, categoryID).Return(existingCategory, nil)

	result, err := uc.UpdateCategory(ctx, categoryID, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.True(t, pkgerrors.IsConflict(err))
}

func TestCategoryUpdate_SystemCategoryDeactivation(t *testing.T) {
	uc, categoryRepo, _, _ := setupCategoryTest()
	ctx := context.Background()

	categoryID := uuid.New()
	existingCategory := &entity.MasterdataCategory{
		ID:       categoryID,
		Code:     "GENDER",
		Name:     "Gender",
		IsSystem: true,
		Version:  1,
	}

	status := string(entity.MasterdataCategoryStatusInactive)
	req := &masterdata.UpdateCategoryRequest{
		Status:  &status,
		Version: 1,
	}

	categoryRepo.On("GetByID", ctx, categoryID).Return(existingCategory, nil)

	result, err := uc.UpdateCategory(ctx, categoryID, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.True(t, pkgerrors.IsValidation(err))
}

func TestCategoryDelete_Success(t *testing.T) {
	uc, categoryRepo, _, cache := setupCategoryTest()
	ctx := context.Background()

	categoryID := uuid.New()
	existingCategory := &entity.MasterdataCategory{
		ID:       categoryID,
		Code:     "CUSTOM",
		Name:     "Custom Category",
		IsSystem: false,
	}

	categoryRepo.On("GetByID", ctx, categoryID).Return(existingCategory, nil)

	categoryRepo.On("Delete", ctx, categoryID).Return(nil)

	cache.On("InvalidateCategories", ctx).Return(nil)

	err := uc.DeleteCategory(ctx, categoryID)

	require.NoError(t, err)
	categoryRepo.AssertExpectations(t)
}

func TestCategoryDelete_SystemCategory(t *testing.T) {
	uc, categoryRepo, _, _ := setupCategoryTest()
	ctx := context.Background()

	categoryID := uuid.New()
	existingCategory := &entity.MasterdataCategory{
		ID:       categoryID,
		Code:     "GENDER",
		Name:     "Gender",
		IsSystem: true,
	}

	categoryRepo.On("GetByID", ctx, categoryID).Return(existingCategory, nil)

	err := uc.DeleteCategory(ctx, categoryID)

	assert.Error(t, err)
	assert.True(t, pkgerrors.IsValidation(err))
}

func TestCategoryGetChildren_Success(t *testing.T) {
	uc, categoryRepo, _, cache := setupCategoryTest()
	ctx := context.Background()

	parentID := uuid.New()
	expectedChildren := []*entity.MasterdataCategory{
		{ID: uuid.New(), Code: "PROVINCE", Name: "Province"},
		{ID: uuid.New(), Code: "CITY", Name: "City"},
	}

	cache.On("GetCategoryChildren", ctx, parentID).Return(nil, stderrors.New("cache miss"))

	categoryRepo.On("GetChildren", ctx, parentID).Return(expectedChildren, nil)

	cache.On("SetCategoryChildren", ctx, parentID, mock.Anything, 5*time.Minute).Return(nil)

	result, err := uc.GetCategoryChildren(ctx, parentID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func ptrString(s string) *string {
	return &s
}
