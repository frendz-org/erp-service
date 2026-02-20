package postgres

import (
	"context"
	"strings"

	"iam-service/entity"
	"iam-service/masterdata/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type masterdataItemRepository struct {
	baseRepository
}

func NewMasterdataItemRepository(db *gorm.DB) contract.ItemRepository {
	return &masterdataItemRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *masterdataItemRepository) List(ctx context.Context, filter *contract.ItemFilter) ([]*entity.MasterdataItem, int64, error) {
	var items []*entity.MasterdataItem
	var total int64

	query := r.getDB(ctx).Model(&entity.MasterdataItem{}).Where("masterdata_items.deleted_at IS NULL")

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	} else if filter.CategoryCode != "" {
		query = query.Joins("JOIN masterdata_categories ON masterdata_categories.id = masterdata_items.category_id").
			Where("masterdata_categories.code = ? AND masterdata_categories.deleted_at IS NULL", filter.CategoryCode)
	}

	if filter.TenantID != nil {
		query = query.Where("(tenant_id IS NULL OR tenant_id = ?)", *filter.TenantID)
	} else {

		query = query.Where("tenant_id IS NULL")
	}

	if filter.ParentID != nil {
		query = query.Where("parent_item_id = ?", *filter.ParentID)
	} else if filter.ParentCode != "" {
		subquery := r.getDB(ctx).
			Model(&entity.MasterdataItem{}).
			Select("id").
			Where("code = ? AND deleted_at IS NULL", filter.ParentCode)
		query = query.Where("parent_item_id IN (?)", subquery)
	}

	if filter.Status != "" {
		query = query.Where("masterdata_items.status = ?", filter.Status)
	}

	if filter.Search != "" {
		searchPattern := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where(
			"(LOWER(masterdata_items.name) LIKE ? OR LOWER(masterdata_items.alt_name) LIKE ? OR LOWER(masterdata_items.code) LIKE ?)",
			searchPattern, searchPattern, searchPattern,
		)
	}

	if filter.IsDefault != nil {
		query = query.Where("is_default = ?", *filter.IsDefault)
	}

	if filter.IsSystem != nil {
		query = query.Where("is_system = ?", *filter.IsSystem)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, translateError(err, "masterdata_item")
	}

	validSortColumns := map[string]string{
		"name":       "masterdata_items.name",
		"code":       "masterdata_items.code",
		"sort_order": "masterdata_items.sort_order",
		"created_at": "masterdata_items.created_at",
		"updated_at": "masterdata_items.updated_at",
	}
	sortColumn := "masterdata_items.sort_order"
	if col, ok := validSortColumns[filter.SortBy]; ok {
		sortColumn = col
	}
	sortOrder := "ASC"
	if filter.SortOrder == "desc" {
		sortOrder = "DESC"
	}
	query = query.Order(sortColumn + " " + sortOrder)

	if filter.Page > 0 && filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		query = query.Offset(offset).Limit(filter.PerPage)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, 0, translateError(err, "masterdata_item")
	}

	return items, total, nil
}

func (r *masterdataItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.MasterdataItem, error) {
	var item entity.MasterdataItem
	err := r.getDB(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&item).Error
	if err != nil {
		return nil, translateError(err, "masterdata_item")
	}
	return &item, nil
}

func (r *masterdataItemRepository) GetByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string) (*entity.MasterdataItem, error) {
	var item entity.MasterdataItem
	query := r.getDB(ctx).
		Where("category_id = ? AND code = ? AND deleted_at IS NULL", categoryID, code)

	if tenantID != nil {
		query = query.Where("(tenant_id IS NULL OR tenant_id = ?)", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}

	err := query.First(&item).Error
	if err != nil {
		return nil, translateError(err, "masterdata_item")
	}
	return &item, nil
}

func (r *masterdataItemRepository) ValidateCode(ctx context.Context, categoryCode string, itemCode string, tenantID *uuid.UUID) (bool, error) {
	var count int64
	query := r.getDB(ctx).
		Model(&entity.MasterdataItem{}).
		Joins("JOIN masterdata_categories ON masterdata_categories.id = masterdata_items.category_id").
		Where("masterdata_categories.code = ?", categoryCode).
		Where("masterdata_items.code = ?", itemCode).
		Where("masterdata_items.status = ?", entity.MasterdataItemStatusActive).
		Where("masterdata_items.deleted_at IS NULL").
		Where("masterdata_categories.deleted_at IS NULL")

	if tenantID != nil {
		query = query.Where("(masterdata_items.tenant_id IS NULL OR masterdata_items.tenant_id = ?)", *tenantID)
	} else {
		query = query.Where("masterdata_items.tenant_id IS NULL")
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, translateError(err, "masterdata_item")
	}
	return count > 0, nil
}

func (r *masterdataItemRepository) Create(ctx context.Context, item *entity.MasterdataItem) error {
	if err := r.getDB(ctx).Create(item).Error; err != nil {
		return translateError(err, "masterdata_item")
	}
	return nil
}

func (r *masterdataItemRepository) Update(ctx context.Context, item *entity.MasterdataItem) error {
	result := r.getDB(ctx).
		Where("id = ? AND deleted_at IS NULL", item.ID).
		Save(item)
	if result.Error != nil {
		return translateError(result.Error, "masterdata_item")
	}
	if result.RowsAffected == 0 {
		return translateError(gorm.ErrRecordNotFound, "masterdata_item")
	}
	return nil
}

func (r *masterdataItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.getDB(ctx).
		Model(&entity.MasterdataItem{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("status", entity.MasterdataItemStatusInactive)
	if result.Error != nil {
		return translateError(result.Error, "masterdata_item")
	}
	if result.RowsAffected == 0 {
		return translateError(gorm.ErrRecordNotFound, "masterdata_item")
	}
	return nil
}

func (r *masterdataItemRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.MasterdataItem, error) {
	var items []*entity.MasterdataItem
	err := r.getDB(ctx).
		Where("parent_item_id = ? AND deleted_at IS NULL AND status = ?", parentID, entity.MasterdataItemStatusActive).
		Order("sort_order ASC").
		Find(&items).Error
	if err != nil {
		return nil, translateError(err, "masterdata_item")
	}
	return items, nil
}

func (r *masterdataItemRepository) GetTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID) ([]*entity.MasterdataItem, error) {
	var items []*entity.MasterdataItem

	query := r.getDB(ctx).
		Joins("JOIN masterdata_categories ON masterdata_categories.id = masterdata_items.category_id").
		Where("masterdata_categories.code = ?", categoryCode).
		Where("masterdata_items.status = ?", entity.MasterdataItemStatusActive).
		Where("masterdata_items.deleted_at IS NULL").
		Where("masterdata_categories.deleted_at IS NULL")

	if tenantID != nil {
		query = query.Where("(masterdata_items.tenant_id IS NULL OR masterdata_items.tenant_id = ?)", *tenantID)
	} else {
		query = query.Where("masterdata_items.tenant_id IS NULL")
	}

	err := query.Order("masterdata_items.sort_order ASC").Find(&items).Error
	if err != nil {
		return nil, translateError(err, "masterdata_item")
	}
	return items, nil
}

func (r *masterdataItemRepository) ListByParent(ctx context.Context, categoryCode string, parentCode string, tenantID *uuid.UUID) ([]*entity.MasterdataItem, error) {
	var items []*entity.MasterdataItem

	parentSubquery := r.getDB(ctx).
		Model(&entity.MasterdataItem{}).
		Select("id").
		Where("code = ? AND deleted_at IS NULL AND status = ?", parentCode, entity.MasterdataItemStatusActive)

	query := r.getDB(ctx).
		Joins("JOIN masterdata_categories ON masterdata_categories.id = masterdata_items.category_id").
		Where("masterdata_categories.code = ?", categoryCode).
		Where("masterdata_items.parent_item_id IN (?)", parentSubquery).
		Where("masterdata_items.status = ?", entity.MasterdataItemStatusActive).
		Where("masterdata_items.deleted_at IS NULL").
		Where("masterdata_categories.deleted_at IS NULL")

	if tenantID != nil {
		query = query.Where("(masterdata_items.tenant_id IS NULL OR masterdata_items.tenant_id = ?)", *tenantID)
	} else {
		query = query.Where("masterdata_items.tenant_id IS NULL")
	}

	err := query.Order("masterdata_items.sort_order ASC").Find(&items).Error
	if err != nil {
		return nil, translateError(err, "masterdata_item")
	}
	return items, nil
}

func (r *masterdataItemRepository) ExistsByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string) (bool, error) {
	var count int64
	query := r.getDB(ctx).
		Model(&entity.MasterdataItem{}).
		Where("category_id = ? AND code = ? AND deleted_at IS NULL", categoryID, code)

	if tenantID != nil {
		query = query.Where("(tenant_id IS NULL OR tenant_id = ?)", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, translateError(err, "masterdata_item")
	}
	return count > 0, nil
}

func (r *masterdataItemRepository) GetDefaultItem(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID) (*entity.MasterdataItem, error) {
	var item entity.MasterdataItem
	query := r.getDB(ctx).
		Where("category_id = ? AND is_default = true AND status = ? AND deleted_at IS NULL", categoryID, entity.MasterdataItemStatusActive)

	if tenantID != nil {
		query = query.Where("(tenant_id IS NULL OR tenant_id = ?)", *tenantID)
	} else {
		query = query.Where("tenant_id IS NULL")
	}

	err := query.First(&item).Error
	if err != nil {
		return nil, translateError(err, "masterdata_item")
	}
	return &item, nil
}
