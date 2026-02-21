package internal

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"erp-service/config"
	"erp-service/entity"
	"erp-service/masterdata/masterdatadto"
	pkgerrors "erp-service/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupItemTest() (*usecase, *MockCategoryRepository, *MockItemRepository, *MockMasterdataCache) {
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

	uc := NewUsecase(cfg, categoryRepo, itemRepo, cache)

	return uc, categoryRepo, itemRepo, cache
}

func TestItemList_Success(t *testing.T) {
	uc, _, itemRepo, cache := setupItemTest()
	ctx := context.Background()

	categoryID := uuid.New()
	req := &masterdatadto.ListItemsRequest{
		CategoryID: &categoryID,
		Page:       1,
		PerPage:    10,
	}

	expectedItems := []*entity.MasterdataItem{
		{
			ID:         uuid.New(),
			CategoryID: categoryID,
			Code:       "MALE",
			Name:       "Male",
			Status:     entity.MasterdataItemStatusActive,
		},
		{
			ID:         uuid.New(),
			CategoryID: categoryID,
			Code:       "FEMALE",
			Name:       "Female",
			Status:     entity.MasterdataItemStatusActive,
		},
	}

	cache.On("GetItemsList", ctx, mock.AnythingOfType("string")).Return(nil, stderrors.New("cache miss"))

	itemRepo.On("List", ctx, mock.AnythingOfType("*contract.ItemFilter")).
		Return(expectedItems, int64(2), nil)

	cache.On("SetItemsList", ctx, mock.AnythingOfType("string"), mock.Anything, 5*time.Minute).Return(nil)

	result, err := uc.ListItems(ctx, req)

	require.NoError(t, err)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, int64(2), result.Pagination.Total)

	itemRepo.AssertExpectations(t)
}

func TestItemGetByID_Success(t *testing.T) {
	uc, _, itemRepo, cache := setupItemTest()
	ctx := context.Background()

	itemID := uuid.New()
	categoryID := uuid.New()
	expectedItem := &entity.MasterdataItem{
		ID:         itemID,
		CategoryID: categoryID,
		Code:       "MALE",
		Name:       "Male",
		Status:     entity.MasterdataItemStatusActive,
	}

	cache.On("GetItemByID", ctx, itemID).Return(nil, stderrors.New("cache miss"))

	itemRepo.On("GetByID", ctx, itemID).Return(expectedItem, nil)

	cache.On("SetItemByID", ctx, itemID, mock.Anything, 5*time.Minute).Return(nil)

	result, err := uc.GetItemByID(ctx, itemID)

	require.NoError(t, err)
	assert.Equal(t, itemID, result.ID)
	assert.Equal(t, "MALE", result.Code)
}

func TestItemGetByID_NotFound(t *testing.T) {
	uc, _, itemRepo, cache := setupItemTest()
	ctx := context.Background()

	itemID := uuid.New()

	cache.On("GetItemByID", ctx, itemID).Return(nil, stderrors.New("cache miss"))

	itemRepo.On("GetByID", ctx, itemID).Return(nil, pkgerrors.ErrNotFound("not found"))

	result, err := uc.GetItemByID(ctx, itemID)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.True(t, pkgerrors.IsNotFound(err))
}

func TestItemGetByCode_Success(t *testing.T) {
	uc, categoryRepo, itemRepo, cache := setupItemTest()
	ctx := context.Background()

	categoryID := uuid.New()
	itemID := uuid.New()

	expectedCategory := &entity.MasterdataCategory{
		ID:   categoryID,
		Code: "GENDER",
		Name: "Gender",
	}

	expectedItem := &entity.MasterdataItem{
		ID:         itemID,
		CategoryID: categoryID,
		Code:       "MALE",
		Name:       "Male",
		Status:     entity.MasterdataItemStatusActive,
	}

	categoryRepo.On("GetByCode", ctx, "GENDER").Return(expectedCategory, nil)

	cache.On("GetItemByCode", ctx, categoryID, (*uuid.UUID)(nil), "MALE").Return(nil, stderrors.New("cache miss"))

	itemRepo.On("GetByCode", ctx, categoryID, (*uuid.UUID)(nil), "MALE").Return(expectedItem, nil)

	cache.On("SetItemByCode", ctx, categoryID, (*uuid.UUID)(nil), "MALE", mock.Anything, 5*time.Minute).Return(nil)

	result, err := uc.GetItemByCode(ctx, "GENDER", nil, "MALE")

	require.NoError(t, err)
	assert.Equal(t, "MALE", result.Code)
}

func TestItemCreate_Success(t *testing.T) {
	uc, categoryRepo, itemRepo, cache := setupItemTest()
	ctx := context.Background()

	categoryID := uuid.New()
	expectedCategory := &entity.MasterdataCategory{
		ID:                 categoryID,
		Code:               "GENDER",
		Name:               "Gender",
		IsTenantExtensible: false,
	}

	req := &masterdatadto.CreateItemRequest{
		CategoryID: categoryID,
		Code:       "MALE",
		Name:       "Male",
	}

	categoryRepo.On("GetByID", ctx, categoryID).Return(expectedCategory, nil)

	itemRepo.On("ExistsByCode", ctx, categoryID, (*uuid.UUID)(nil), "MALE").Return(false, nil)

	itemRepo.On("Create", ctx, mock.AnythingOfType("*entity.MasterdataItem")).Return(nil)

	cache.On("InvalidateItems", ctx).Return(nil)

	result, err := uc.CreateItem(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, "MALE", result.Code)
	assert.Equal(t, "Male", result.Name)
	assert.Equal(t, string(entity.MasterdataItemStatusActive), result.Status)

	itemRepo.AssertExpectations(t)
}

func TestItemCreate_CodeExists(t *testing.T) {
	uc, categoryRepo, itemRepo, _ := setupItemTest()
	ctx := context.Background()

	categoryID := uuid.New()
	expectedCategory := &entity.MasterdataCategory{
		ID:   categoryID,
		Code: "GENDER",
		Name: "Gender",
	}

	req := &masterdatadto.CreateItemRequest{
		CategoryID: categoryID,
		Code:       "MALE",
		Name:       "Male",
	}

	categoryRepo.On("GetByID", ctx, categoryID).Return(expectedCategory, nil)

	itemRepo.On("ExistsByCode", ctx, categoryID, (*uuid.UUID)(nil), "MALE").Return(true, nil)

	result, err := uc.CreateItem(ctx, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.True(t, pkgerrors.IsConflict(err))
}

func TestItemCreate_CategoryNotExtensible(t *testing.T) {
	uc, categoryRepo, _, _ := setupItemTest()
	ctx := context.Background()

	categoryID := uuid.New()
	tenantID := uuid.New()
	expectedCategory := &entity.MasterdataCategory{
		ID:                 categoryID,
		Code:               "GENDER",
		Name:               "Gender",
		IsTenantExtensible: false,
	}

	req := &masterdatadto.CreateItemRequest{
		CategoryID: categoryID,
		TenantID:   &tenantID,
		Code:       "OTHER",
		Name:       "Other",
	}

	categoryRepo.On("GetByID", ctx, categoryID).Return(expectedCategory, nil)

	result, err := uc.CreateItem(ctx, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.True(t, pkgerrors.IsValidation(err))
}

func TestItemUpdate_Success(t *testing.T) {
	uc, _, itemRepo, cache := setupItemTest()
	ctx := context.Background()

	itemID := uuid.New()
	categoryID := uuid.New()
	existingItem := &entity.MasterdataItem{
		ID:         itemID,
		CategoryID: categoryID,
		Code:       "MALE",
		Name:       "Male",
		Status:     entity.MasterdataItemStatusActive,
		Version:    1,
	}

	newName := "Male Gender"
	req := &masterdatadto.UpdateItemRequest{
		Name:    &newName,
		Version: 1,
	}

	itemRepo.On("GetByID", ctx, itemID).Return(existingItem, nil)

	itemRepo.On("Update", ctx, mock.AnythingOfType("*entity.MasterdataItem")).Return(nil)

	cache.On("InvalidateItems", ctx).Return(nil)

	result, err := uc.UpdateItem(ctx, itemID, req)

	require.NoError(t, err)
	assert.Equal(t, "Male Gender", result.Name)
	assert.Equal(t, 2, result.Version)
}

func TestItemUpdate_VersionConflict(t *testing.T) {
	uc, _, itemRepo, _ := setupItemTest()
	ctx := context.Background()

	itemID := uuid.New()
	categoryID := uuid.New()
	existingItem := &entity.MasterdataItem{
		ID:         itemID,
		CategoryID: categoryID,
		Code:       "MALE",
		Name:       "Male",
		Version:    2,
	}

	newName := "Male Gender"
	req := &masterdatadto.UpdateItemRequest{
		Name:    &newName,
		Version: 1,
	}

	itemRepo.On("GetByID", ctx, itemID).Return(existingItem, nil)

	result, err := uc.UpdateItem(ctx, itemID, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.True(t, pkgerrors.IsConflict(err))
}

func TestItemUpdate_SystemItemDeactivation(t *testing.T) {
	uc, _, itemRepo, _ := setupItemTest()
	ctx := context.Background()

	itemID := uuid.New()
	categoryID := uuid.New()
	existingItem := &entity.MasterdataItem{
		ID:         itemID,
		CategoryID: categoryID,
		Code:       "MALE",
		Name:       "Male",
		IsSystem:   true,
		Version:    1,
	}

	status := string(entity.MasterdataItemStatusInactive)
	req := &masterdatadto.UpdateItemRequest{
		Status:  &status,
		Version: 1,
	}

	itemRepo.On("GetByID", ctx, itemID).Return(existingItem, nil)

	result, err := uc.UpdateItem(ctx, itemID, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.True(t, pkgerrors.IsValidation(err))
}

func TestItemDelete_Success(t *testing.T) {
	uc, _, itemRepo, cache := setupItemTest()
	ctx := context.Background()

	itemID := uuid.New()
	categoryID := uuid.New()
	existingItem := &entity.MasterdataItem{
		ID:         itemID,
		CategoryID: categoryID,
		Code:       "CUSTOM",
		Name:       "Custom Item",
		IsSystem:   false,
	}

	itemRepo.On("GetByID", ctx, itemID).Return(existingItem, nil)

	itemRepo.On("Delete", ctx, itemID).Return(nil)

	cache.On("InvalidateItems", ctx).Return(nil)

	err := uc.DeleteItem(ctx, itemID)

	require.NoError(t, err)
	itemRepo.AssertExpectations(t)
}

func TestItemDelete_SystemItem(t *testing.T) {
	uc, _, itemRepo, _ := setupItemTest()
	ctx := context.Background()

	itemID := uuid.New()
	categoryID := uuid.New()
	existingItem := &entity.MasterdataItem{
		ID:         itemID,
		CategoryID: categoryID,
		Code:       "MALE",
		Name:       "Male",
		IsSystem:   true,
	}

	itemRepo.On("GetByID", ctx, itemID).Return(existingItem, nil)

	err := uc.DeleteItem(ctx, itemID)

	assert.Error(t, err)
	assert.True(t, pkgerrors.IsValidation(err))
}

func TestItemGetChildren_Success(t *testing.T) {
	uc, _, itemRepo, cache := setupItemTest()
	ctx := context.Background()

	parentID := uuid.New()
	categoryID := uuid.New()
	expectedChildren := []*entity.MasterdataItem{
		{ID: uuid.New(), CategoryID: categoryID, Code: "JAKARTA", Name: "Jakarta"},
		{ID: uuid.New(), CategoryID: categoryID, Code: "BANDUNG", Name: "Bandung"},
	}

	cache.On("GetItemChildren", ctx, parentID).Return(nil, stderrors.New("cache miss"))

	itemRepo.On("GetChildren", ctx, parentID).Return(expectedChildren, nil)

	cache.On("SetItemChildren", ctx, parentID, mock.Anything, 5*time.Minute).Return(nil)

	result, err := uc.GetItemChildren(ctx, parentID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestItemGetTree_Success(t *testing.T) {
	uc, _, itemRepo, cache := setupItemTest()
	ctx := context.Background()

	categoryCode := "LOCATION"
	categoryID := uuid.New()
	expectedItems := []*entity.MasterdataItem{
		{ID: uuid.New(), CategoryID: categoryID, Code: "ID", Name: "Indonesia"},
		{ID: uuid.New(), CategoryID: categoryID, Code: "WJ", Name: "West Java"},
		{ID: uuid.New(), CategoryID: categoryID, Code: "BDG", Name: "Bandung"},
	}

	cache.On("GetItemTree", ctx, categoryCode, (*uuid.UUID)(nil)).Return(nil, stderrors.New("cache miss"))

	itemRepo.On("GetTree", ctx, categoryCode, (*uuid.UUID)(nil)).Return(expectedItems, nil)

	cache.On("SetItemTree", ctx, categoryCode, (*uuid.UUID)(nil), mock.Anything, 10*time.Minute).Return(nil)

	result, err := uc.GetItemTree(ctx, categoryCode, nil)

	require.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestItemGetDefault_Success(t *testing.T) {
	uc, categoryRepo, itemRepo, cache := setupItemTest()
	ctx := context.Background()

	categoryID := uuid.New()
	itemID := uuid.New()

	expectedCategory := &entity.MasterdataCategory{
		ID:   categoryID,
		Code: "GENDER",
		Name: "Gender",
	}

	expectedItem := &entity.MasterdataItem{
		ID:         itemID,
		CategoryID: categoryID,
		Code:       "MALE",
		Name:       "Male",
		IsDefault:  true,
		Status:     entity.MasterdataItemStatusActive,
	}

	categoryRepo.On("GetByCode", ctx, "GENDER").Return(expectedCategory, nil)

	cache.On("GetItemDefault", ctx, categoryID, (*uuid.UUID)(nil)).Return(nil, stderrors.New("cache miss"))

	itemRepo.On("GetDefaultItem", ctx, categoryID, (*uuid.UUID)(nil)).Return(expectedItem, nil)

	cache.On("SetItemDefault", ctx, categoryID, (*uuid.UUID)(nil), mock.Anything, 5*time.Minute).Return(nil)

	result, err := uc.GetItemDefault(ctx, "GENDER", nil)

	require.NoError(t, err)
	assert.Equal(t, "MALE", result.Code)
	assert.True(t, result.IsDefault)
}

func TestItemGetDefault_NotFound(t *testing.T) {
	uc, categoryRepo, itemRepo, cache := setupItemTest()
	ctx := context.Background()

	categoryID := uuid.New()

	expectedCategory := &entity.MasterdataCategory{
		ID:   categoryID,
		Code: "GENDER",
		Name: "Gender",
	}

	categoryRepo.On("GetByCode", ctx, "GENDER").Return(expectedCategory, nil)

	cache.On("GetItemDefault", ctx, categoryID, (*uuid.UUID)(nil)).Return(nil, stderrors.New("cache miss"))

	itemRepo.On("GetDefaultItem", ctx, categoryID, (*uuid.UUID)(nil)).Return(nil, pkgerrors.ErrNotFound("not found"))

	result, err := uc.GetItemDefault(ctx, "GENDER", nil)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.True(t, pkgerrors.IsNotFound(err))
}

func TestItemValidateCode_Success(t *testing.T) {
	uc, _, itemRepo, _ := setupItemTest()
	ctx := context.Background()

	req := &masterdatadto.ValidateCodeRequest{
		CategoryCode: "GENDER",
		ItemCode:     "MALE",
	}

	itemRepo.On("ValidateCode", ctx, "GENDER", "MALE", (*uuid.UUID)(nil)).Return(true, nil)

	result, err := uc.ValidateItemCode(ctx, req)

	require.NoError(t, err)
	assert.True(t, result.Valid)
	assert.Equal(t, "GENDER", result.CategoryCode)
	assert.Equal(t, "MALE", result.ItemCode)
	assert.Empty(t, result.Message)
}

func TestItemValidateCode_Invalid(t *testing.T) {
	uc, _, itemRepo, _ := setupItemTest()
	ctx := context.Background()

	req := &masterdatadto.ValidateCodeRequest{
		CategoryCode: "GENDER",
		ItemCode:     "INVALID",
	}

	itemRepo.On("ValidateCode", ctx, "GENDER", "INVALID", (*uuid.UUID)(nil)).Return(false, nil)

	result, err := uc.ValidateItemCode(ctx, req)

	require.NoError(t, err)
	assert.False(t, result.Valid)
	assert.Contains(t, result.Message, "INVALID")
	assert.Contains(t, result.Message, "GENDER")
}

func TestItemValidateCodes_AllValid(t *testing.T) {
	uc, _, itemRepo, _ := setupItemTest()
	ctx := context.Background()

	req := &masterdatadto.ValidateCodesRequest{
		Validations: []masterdatadto.ValidationItem{
			{CategoryCode: "GENDER", ItemCode: "MALE"},
			{CategoryCode: "GENDER", ItemCode: "FEMALE"},
		},
	}

	itemRepo.On("ValidateCode", ctx, "GENDER", "MALE", (*uuid.UUID)(nil)).Return(true, nil)
	itemRepo.On("ValidateCode", ctx, "GENDER", "FEMALE", (*uuid.UUID)(nil)).Return(true, nil)

	result, err := uc.ValidateItemCodes(ctx, req)

	require.NoError(t, err)
	assert.True(t, result.AllValid)
	assert.Len(t, result.Results, 2)
	assert.True(t, result.Results[0].Valid)
	assert.True(t, result.Results[1].Valid)
}

func TestItemValidateCodes_SomeInvalid(t *testing.T) {
	uc, _, itemRepo, _ := setupItemTest()
	ctx := context.Background()

	req := &masterdatadto.ValidateCodesRequest{
		Validations: []masterdatadto.ValidationItem{
			{CategoryCode: "GENDER", ItemCode: "MALE"},
			{CategoryCode: "GENDER", ItemCode: "INVALID"},
		},
	}

	itemRepo.On("ValidateCode", ctx, "GENDER", "MALE", (*uuid.UUID)(nil)).Return(true, nil)
	itemRepo.On("ValidateCode", ctx, "GENDER", "INVALID", (*uuid.UUID)(nil)).Return(false, nil)

	result, err := uc.ValidateItemCodes(ctx, req)

	require.NoError(t, err)
	assert.False(t, result.AllValid)
	assert.True(t, result.Results[0].Valid)
	assert.False(t, result.Results[1].Valid)
	assert.Contains(t, result.Results[1].Message, "not found")
}

func TestItemListByParent_Success(t *testing.T) {
	uc, _, itemRepo, _ := setupItemTest()
	ctx := context.Background()

	categoryCode := "CITY"
	parentCode := "WEST_JAVA"
	categoryID := uuid.New()

	expectedItems := []*entity.MasterdataItem{
		{ID: uuid.New(), CategoryID: categoryID, Code: "BANDUNG", Name: "Bandung"},
		{ID: uuid.New(), CategoryID: categoryID, Code: "CIREBON", Name: "Cirebon"},
	}

	itemRepo.On("ListByParent", ctx, categoryCode, parentCode, (*uuid.UUID)(nil)).Return(expectedItems, nil)

	result, err := uc.ListItemsByParent(ctx, categoryCode, parentCode, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
