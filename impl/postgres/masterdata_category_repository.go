package postgres

import (
	"context"

	"erp-service/entity"
	"erp-service/masterdata"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type masterdataCategoryRepository struct {
	baseRepository
}

func NewMasterdataCategoryRepository(db *gorm.DB) masterdata.CategoryRepository {
	return &masterdataCategoryRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *masterdataCategoryRepository) List(ctx context.Context, filter *masterdata.CategoryFilter) ([]*entity.MasterdataCategory, int64, error) {
	var categories []*entity.MasterdataCategory
	var total int64

	query := r.getDB(ctx).Model(&entity.MasterdataCategory{}).Where("deleted_at IS NULL")

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.IsSystem != nil {
		query = query.Where("is_system = ?", *filter.IsSystem)
	}
	if filter.ParentID != nil {
		query = query.Where("parent_category_id = ?", *filter.ParentID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, translateError(err, "masterdata_category")
	}

	validSortColumns := map[string]string{
		"name":       "name",
		"code":       "code",
		"sort_order": "sort_order",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
	sortColumn := "sort_order"
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

	if err := query.Find(&categories).Error; err != nil {
		return nil, 0, translateError(err, "masterdata_category")
	}

	return categories, total, nil
}

func (r *masterdataCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.MasterdataCategory, error) {
	var category entity.MasterdataCategory
	err := r.getDB(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&category).Error
	if err != nil {
		return nil, translateError(err, "masterdata_category")
	}
	return &category, nil
}

func (r *masterdataCategoryRepository) GetByCode(ctx context.Context, code string) (*entity.MasterdataCategory, error) {
	var category entity.MasterdataCategory
	err := r.getDB(ctx).
		Where("code = ? AND deleted_at IS NULL", code).
		First(&category).Error
	if err != nil {
		return nil, translateError(err, "masterdata_category")
	}
	return &category, nil
}

func (r *masterdataCategoryRepository) Create(ctx context.Context, category *entity.MasterdataCategory) error {
	if err := r.getDB(ctx).Create(category).Error; err != nil {
		return translateError(err, "masterdata_category")
	}
	return nil
}

func (r *masterdataCategoryRepository) Update(ctx context.Context, category *entity.MasterdataCategory) error {
	result := r.getDB(ctx).
		Where("id = ? AND deleted_at IS NULL", category.ID).
		Save(category)
	if result.Error != nil {
		return translateError(result.Error, "masterdata_category")
	}
	if result.RowsAffected == 0 {
		return translateError(gorm.ErrRecordNotFound, "masterdata_category")
	}
	return nil
}

func (r *masterdataCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.getDB(ctx).
		Model(&entity.MasterdataCategory{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("status", entity.MasterdataCategoryStatusInactive)
	if result.Error != nil {
		return translateError(result.Error, "masterdata_category")
	}
	if result.RowsAffected == 0 {
		return translateError(gorm.ErrRecordNotFound, "masterdata_category")
	}
	return nil
}

func (r *masterdataCategoryRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.MasterdataCategory, error) {
	var categories []*entity.MasterdataCategory
	err := r.getDB(ctx).
		Where("parent_category_id = ? AND deleted_at IS NULL AND status = ?", parentID, entity.MasterdataCategoryStatusActive).
		Order("sort_order ASC").
		Find(&categories).Error
	if err != nil {
		return nil, translateError(err, "masterdata_category")
	}
	return categories, nil
}

func (r *masterdataCategoryRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.getDB(ctx).
		Model(&entity.MasterdataCategory{}).
		Where("code = ? AND deleted_at IS NULL", code).
		Count(&count).Error
	if err != nil {
		return false, translateError(err, "masterdata_category")
	}
	return count > 0, nil
}
