package postgres_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"erp-service/entity"
	implpg "erp-service/impl/postgres"
	"erp-service/masterdata"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestMasterdataCategoryRepository_GetByID(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := implpg.NewMasterdataCategoryRepository(gormDB)
	ctx := context.Background()

	categoryID := uuid.New()
	now := time.Now()

	t.Run("existing category", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "code", "name", "description", "parent_category_id",
			"is_system", "is_tenant_extensible", "sort_order", "status",
			"metadata", "version", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			categoryID, "GENDER", "Gender", nil, nil,
			true, false, 1, "ACTIVE",
			[]byte("{}"), 1, now, now, nil,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_categories" WHERE id = $1 AND deleted_at IS NULL ORDER BY "masterdata_categories"."id" LIMIT $2`)).
			WithArgs(categoryID, 1).
			WillReturnRows(rows)

		category, err := repo.GetByID(ctx, categoryID)

		require.NoError(t, err)
		assert.Equal(t, categoryID, category.ID)
		assert.Equal(t, "GENDER", category.Code)
		assert.Equal(t, "Gender", category.Name)
		assert.True(t, category.IsSystem)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("category not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_categories" WHERE id = $1 AND deleted_at IS NULL ORDER BY "masterdata_categories"."id" LIMIT $2`)).
			WithArgs(categoryID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		category, err := repo.GetByID(ctx, categoryID)

		assert.Error(t, err)
		assert.Nil(t, category)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMasterdataCategoryRepository_GetByCode(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := implpg.NewMasterdataCategoryRepository(gormDB)
	ctx := context.Background()

	categoryID := uuid.New()
	now := time.Now()

	t.Run("existing category", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "code", "name", "description", "parent_category_id",
			"is_system", "is_tenant_extensible", "sort_order", "status",
			"metadata", "version", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			categoryID, "GENDER", "Gender", nil, nil,
			true, false, 1, "ACTIVE",
			[]byte("{}"), 1, now, now, nil,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_categories" WHERE code = $1 AND deleted_at IS NULL ORDER BY "masterdata_categories"."id" LIMIT $2`)).
			WithArgs("GENDER", 1).
			WillReturnRows(rows)

		category, err := repo.GetByCode(ctx, "GENDER")

		require.NoError(t, err)
		assert.Equal(t, "GENDER", category.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("category not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_categories" WHERE code = $1 AND deleted_at IS NULL ORDER BY "masterdata_categories"."id" LIMIT $2`)).
			WithArgs("NONEXISTENT", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		category, err := repo.GetByCode(ctx, "NONEXISTENT")

		assert.Error(t, err)
		assert.Nil(t, category)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMasterdataCategoryRepository_Create(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := implpg.NewMasterdataCategoryRepository(gormDB)
	ctx := context.Background()

	t.Run("successful create", func(t *testing.T) {
		category := &entity.MasterdataCategory{
			Code:     "TEST_CATEGORY",
			Name:     "Test Category",
			Status:   entity.MasterdataCategoryStatusActive,
			Metadata: []byte("{}"),
			Version:  1,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "masterdata_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectCommit()

		err := repo.Create(ctx, category)

		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMasterdataCategoryRepository_List(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := implpg.NewMasterdataCategoryRepository(gormDB)
	ctx := context.Background()

	now := time.Now()

	t.Run("list all active categories", func(t *testing.T) {
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "masterdata_categories" WHERE deleted_at IS NULL AND status = $1`)).
			WithArgs("ACTIVE").
			WillReturnRows(countRows)

		dataRows := sqlmock.NewRows([]string{
			"id", "code", "name", "description", "parent_category_id",
			"is_system", "is_tenant_extensible", "sort_order", "status",
			"metadata", "version", "created_at", "updated_at", "deleted_at",
		}).
			AddRow(uuid.New(), "GENDER", "Gender", nil, nil, true, false, 1, "ACTIVE", []byte("{}"), 1, now, now, nil).
			AddRow(uuid.New(), "MARITAL_STATUS", "Marital Status", nil, nil, true, false, 2, "ACTIVE", []byte("{}"), 1, now, now, nil)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_categories" WHERE deleted_at IS NULL AND status = $1 ORDER BY sort_order ASC`)).
			WithArgs("ACTIVE").
			WillReturnRows(dataRows)

		filter := &masterdata.CategoryFilter{
			Status: "ACTIVE",
		}
		categories, total, err := repo.List(ctx, filter)

		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, categories, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("list with pagination", func(t *testing.T) {
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(10)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "masterdata_categories" WHERE deleted_at IS NULL`)).
			WillReturnRows(countRows)

		dataRows := sqlmock.NewRows([]string{
			"id", "code", "name", "description", "parent_category_id",
			"is_system", "is_tenant_extensible", "sort_order", "status",
			"metadata", "version", "created_at", "updated_at", "deleted_at",
		}).AddRow(uuid.New(), "GENDER", "Gender", nil, nil, true, false, 1, "ACTIVE", []byte("{}"), 1, now, now, nil)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_categories" WHERE deleted_at IS NULL ORDER BY sort_order ASC LIMIT $1`)).
			WithArgs(5).
			WillReturnRows(dataRows)

		filter := &masterdata.CategoryFilter{
			Page:    1,
			PerPage: 5,
		}
		categories, total, err := repo.List(ctx, filter)

		require.NoError(t, err)
		assert.Equal(t, int64(10), total)
		assert.Len(t, categories, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMasterdataCategoryRepository_ExistsByCode(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := implpg.NewMasterdataCategoryRepository(gormDB)
	ctx := context.Background()

	t.Run("category exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "masterdata_categories" WHERE code = $1 AND deleted_at IS NULL`)).
			WithArgs("GENDER").
			WillReturnRows(rows)

		exists, err := repo.ExistsByCode(ctx, "GENDER")

		require.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("category does not exist", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "masterdata_categories" WHERE code = $1 AND deleted_at IS NULL`)).
			WithArgs("NONEXISTENT").
			WillReturnRows(rows)

		exists, err := repo.ExistsByCode(ctx, "NONEXISTENT")

		require.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMasterdataCategoryRepository_Delete(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := implpg.NewMasterdataCategoryRepository(gormDB)
	ctx := context.Background()

	categoryID := uuid.New()

	t.Run("successful delete", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(`UPDATE "masterdata_categories" SET "status"=\$1,"updated_at"=\$2 WHERE id = \$3 AND deleted_at IS NULL`).
			WithArgs(entity.MasterdataCategoryStatusInactive, sqlmock.AnyArg(), categoryID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete(ctx, categoryID)

		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("category not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "masterdata_categories" SET "status"=\$1,"updated_at"=\$2 WHERE id = \$3 AND deleted_at IS NULL`).
			WithArgs(entity.MasterdataCategoryStatusInactive, sqlmock.AnyArg(), categoryID).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.Delete(ctx, categoryID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMasterdataCategoryRepository_GetChildren(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := implpg.NewMasterdataCategoryRepository(gormDB)
	ctx := context.Background()

	parentID := uuid.New()
	now := time.Now()

	t.Run("get children", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "code", "name", "description", "parent_category_id",
			"is_system", "is_tenant_extensible", "sort_order", "status",
			"metadata", "version", "created_at", "updated_at", "deleted_at",
		}).
			AddRow(uuid.New(), "PROVINCE", "Province", nil, parentID, true, false, 1, "ACTIVE", []byte("{}"), 1, now, now, sql.NullTime{})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "masterdata_categories" WHERE parent_category_id = $1 AND deleted_at IS NULL AND status = $2 ORDER BY sort_order ASC`)).
			WithArgs(parentID, entity.MasterdataCategoryStatusActive).
			WillReturnRows(rows)

		children, err := repo.GetChildren(ctx, parentID)

		require.NoError(t, err)
		assert.Len(t, children, 1)
		assert.Equal(t, "PROVINCE", children[0].Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
