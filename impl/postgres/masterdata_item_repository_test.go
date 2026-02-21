package postgres

import (
	"context"
	"regexp"
	"testing"
	"time"

	"erp-service/entity"
	"erp-service/masterdata"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMasterdataItemRepository_GetByID(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewMasterdataItemRepository(gormDB)
	ctx := context.Background()

	itemID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	t.Run("existing item", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "category_id", "tenant_id", "parent_item_id", "code", "name",
			"alt_name", "description", "sort_order", "is_system", "is_default",
			"status", "effective_from", "effective_until", "metadata",
			"created_by", "version", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			itemID, categoryID, nil, nil, "MALE", "Male",
			nil, nil, 1, true, false,
			"ACTIVE", nil, nil, []byte("{}"),
			nil, 1, now, now, nil,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_items" WHERE id = $1 AND deleted_at IS NULL ORDER BY "masterdata_items"."id" LIMIT $2`)).
			WithArgs(itemID, 1).
			WillReturnRows(rows)

		item, err := repo.GetByID(ctx, itemID)

		require.NoError(t, err)
		assert.Equal(t, itemID, item.ID)
		assert.Equal(t, "MALE", item.Code)
		assert.Equal(t, "Male", item.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("item not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_items" WHERE id = $1 AND deleted_at IS NULL ORDER BY "masterdata_items"."id" LIMIT $2`)).
			WithArgs(itemID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		item, err := repo.GetByID(ctx, itemID)

		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMasterdataItemRepository_GetByCode(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewMasterdataItemRepository(gormDB)
	ctx := context.Background()

	itemID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	t.Run("existing item - global", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "category_id", "tenant_id", "parent_item_id", "code", "name",
			"alt_name", "description", "sort_order", "is_system", "is_default",
			"status", "effective_from", "effective_until", "metadata",
			"created_by", "version", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			itemID, categoryID, nil, nil, "MALE", "Male",
			nil, nil, 1, true, false,
			"ACTIVE", nil, nil, []byte("{}"),
			nil, 1, now, now, nil,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_items" WHERE (category_id = $1 AND code = $2 AND deleted_at IS NULL) AND tenant_id IS NULL ORDER BY "masterdata_items"."id" LIMIT $3`)).
			WithArgs(categoryID, "MALE", 1).
			WillReturnRows(rows)

		item, err := repo.GetByCode(ctx, categoryID, nil, "MALE")

		require.NoError(t, err)
		assert.Equal(t, "MALE", item.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("existing item - with tenant", func(t *testing.T) {
		tenantID := uuid.New()
		rows := sqlmock.NewRows([]string{
			"id", "category_id", "tenant_id", "parent_item_id", "code", "name",
			"alt_name", "description", "sort_order", "is_system", "is_default",
			"status", "effective_from", "effective_until", "metadata",
			"created_by", "version", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			itemID, categoryID, tenantID, nil, "DEPT_IT", "IT Department",
			nil, nil, 1, false, false,
			"ACTIVE", nil, nil, []byte("{}"),
			nil, 1, now, now, nil,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_items" WHERE (category_id = $1 AND code = $2 AND deleted_at IS NULL) AND ((tenant_id IS NULL OR tenant_id = $3)) ORDER BY "masterdata_items"."id" LIMIT $4`)).
			WithArgs(categoryID, "DEPT_IT", tenantID, 1).
			WillReturnRows(rows)

		item, err := repo.GetByCode(ctx, categoryID, &tenantID, "DEPT_IT")

		require.NoError(t, err)
		assert.Equal(t, "DEPT_IT", item.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMasterdataItemRepository_ValidateCode(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewMasterdataItemRepository(gormDB)
	ctx := context.Background()

	t.Run("valid code exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "masterdata_items"`)).
			WillReturnRows(rows)

		valid, err := repo.ValidateCode(ctx, "GENDER", "MALE", nil)

		require.NoError(t, err)
		assert.True(t, valid)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid code", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "masterdata_items"`)).
			WillReturnRows(rows)

		valid, err := repo.ValidateCode(ctx, "GENDER", "INVALID", nil)

		require.NoError(t, err)
		assert.False(t, valid)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMasterdataItemRepository_Create(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewMasterdataItemRepository(gormDB)
	ctx := context.Background()

	categoryID := uuid.New()

	t.Run("successful create", func(t *testing.T) {
		item := &entity.MasterdataItem{
			CategoryID: categoryID,
			Code:       "NEW_ITEM",
			Name:       "New Item",
			Status:     entity.MasterdataItemStatusActive,
			Metadata:   []byte("{}"),
			Version:    1,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "masterdata_items"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectCommit()

		err := repo.Create(ctx, item)

		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMasterdataItemRepository_List(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewMasterdataItemRepository(gormDB)
	ctx := context.Background()

	categoryID := uuid.New()
	now := time.Now()

	t.Run("list by category", func(t *testing.T) {
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(3)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "masterdata_items"`)).
			WillReturnRows(countRows)

		dataRows := sqlmock.NewRows([]string{
			"id", "category_id", "tenant_id", "parent_item_id", "code", "name",
			"alt_name", "description", "sort_order", "is_system", "is_default",
			"status", "effective_from", "effective_until", "metadata",
			"created_by", "version", "created_at", "updated_at", "deleted_at",
		}).
			AddRow(uuid.New(), categoryID, nil, nil, "MALE", "Male", nil, nil, 1, true, false, "ACTIVE", nil, nil, []byte("{}"), nil, 1, now, now, nil).
			AddRow(uuid.New(), categoryID, nil, nil, "FEMALE", "Female", nil, nil, 2, true, false, "ACTIVE", nil, nil, []byte("{}"), nil, 1, now, now, nil)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_items"`)).
			WillReturnRows(dataRows)

		filter := &masterdata.ItemFilter{
			CategoryID: &categoryID,
			Status:     "ACTIVE",
		}
		items, total, err := repo.List(ctx, filter)

		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, items, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMasterdataItemRepository_Delete(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewMasterdataItemRepository(gormDB)
	ctx := context.Background()

	itemID := uuid.New()

	t.Run("successful delete", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(`UPDATE "masterdata_items" SET "status"=\$1,"updated_at"=\$2 WHERE id = \$3 AND deleted_at IS NULL`).
			WithArgs(entity.MasterdataItemStatusInactive, sqlmock.AnyArg(), itemID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete(ctx, itemID)

		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("item not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "masterdata_items" SET "status"=\$1,"updated_at"=\$2 WHERE id = \$3 AND deleted_at IS NULL`).
			WithArgs(entity.MasterdataItemStatusInactive, sqlmock.AnyArg(), itemID).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.Delete(ctx, itemID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMasterdataItemRepository_ExistsByCode(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewMasterdataItemRepository(gormDB)
	ctx := context.Background()

	categoryID := uuid.New()

	t.Run("item exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "masterdata_items" WHERE (category_id = $1 AND code = $2 AND deleted_at IS NULL) AND tenant_id IS NULL`)).
			WithArgs(categoryID, "MALE").
			WillReturnRows(rows)

		exists, err := repo.ExistsByCode(ctx, categoryID, nil, "MALE")

		require.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("item does not exist", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "masterdata_items" WHERE (category_id = $1 AND code = $2 AND deleted_at IS NULL) AND tenant_id IS NULL`)).
			WithArgs(categoryID, "NONEXISTENT").
			WillReturnRows(rows)

		exists, err := repo.ExistsByCode(ctx, categoryID, nil, "NONEXISTENT")

		require.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMasterdataItemRepository_GetChildren(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewMasterdataItemRepository(gormDB)
	ctx := context.Background()

	parentID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	t.Run("get children", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "category_id", "tenant_id", "parent_item_id", "code", "name",
			"alt_name", "description", "sort_order", "is_system", "is_default",
			"status", "effective_from", "effective_until", "metadata",
			"created_by", "version", "created_at", "updated_at", "deleted_at",
		}).
			AddRow(uuid.New(), categoryID, nil, parentID, "CHILD1", "Child 1", nil, nil, 1, true, false, "ACTIVE", nil, nil, []byte("{}"), nil, 1, now, now, nil).
			AddRow(uuid.New(), categoryID, nil, parentID, "CHILD2", "Child 2", nil, nil, 2, true, false, "ACTIVE", nil, nil, []byte("{}"), nil, 1, now, now, nil)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_items" WHERE parent_item_id = $1 AND deleted_at IS NULL AND status = $2 ORDER BY sort_order ASC`)).
			WithArgs(parentID, entity.MasterdataItemStatusActive).
			WillReturnRows(rows)

		children, err := repo.GetChildren(ctx, parentID)

		require.NoError(t, err)
		assert.Len(t, children, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMasterdataItemRepository_GetDefaultItem(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewMasterdataItemRepository(gormDB)
	ctx := context.Background()

	categoryID := uuid.New()
	itemID := uuid.New()
	now := time.Now()

	t.Run("default item exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "category_id", "tenant_id", "parent_item_id", "code", "name",
			"alt_name", "description", "sort_order", "is_system", "is_default",
			"status", "effective_from", "effective_until", "metadata",
			"created_by", "version", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			itemID, categoryID, nil, nil, "MALE", "Male",
			nil, nil, 1, true, true,
			"ACTIVE", nil, nil, []byte("{}"),
			nil, 1, now, now, nil,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_items" WHERE (category_id = $1 AND is_default = true AND status = $2 AND deleted_at IS NULL) AND tenant_id IS NULL ORDER BY "masterdata_items"."id" LIMIT $3`)).
			WithArgs(categoryID, entity.MasterdataItemStatusActive, 1).
			WillReturnRows(rows)

		item, err := repo.GetDefaultItem(ctx, categoryID, nil)

		require.NoError(t, err)
		assert.Equal(t, "MALE", item.Code)
		assert.True(t, item.IsDefault)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no default item", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_items" WHERE (category_id = $1 AND is_default = true AND status = $2 AND deleted_at IS NULL) AND tenant_id IS NULL ORDER BY "masterdata_items"."id" LIMIT $3`)).
			WithArgs(categoryID, entity.MasterdataItemStatusActive, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		item, err := repo.GetDefaultItem(ctx, categoryID, nil)

		assert.Error(t, err)
		assert.Nil(t, item)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
