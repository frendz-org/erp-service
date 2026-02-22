package postgres

import (
	"context"
	"time"

	"erp-service/entity"
	apperrors "erp-service/pkg/errors"
	"erp-service/saving/participant"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fileRepository struct {
	baseRepository
}

func NewFileRepository(db *gorm.DB) participant.FileRepository {
	return &fileRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *fileRepository) Create(ctx context.Context, file *entity.File) error {
	if err := r.getDB(ctx).Create(file).Error; err != nil {
		return translateError(err, "file")
	}
	return nil
}

func (r *fileRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.File, error) {
	var file entity.File

	err := r.getDB(ctx).Where("id = ?", id).First(&file).Error
	if err != nil {
		return nil, translateError(err, "file")
	}
	return &file, nil
}

func (r *fileRepository) SetPermanent(ctx context.Context, id uuid.UUID) error {
	result := r.getDB(ctx).Model(&entity.File{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"expires_at": nil,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return translateError(result.Error, "file")
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound("file not found")
	}
	return nil
}

func (r *fileRepository) SetExpiring(ctx context.Context, id uuid.UUID, expiry time.Time) error {
	result := r.getDB(ctx).Model(&entity.File{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"expires_at": expiry,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return translateError(result.Error, "file")
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound("file not found")
	}
	return nil
}

func (r *fileRepository) ListExpired(ctx context.Context, limit int) ([]*entity.File, error) {
	var files []*entity.File
	err := r.getDB(ctx).
		Where("expires_at <= ? AND deleted_at IS NULL AND failed_delete_attempts < 5", time.Now()).
		Order("expires_at ASC").
		Limit(limit).
		Find(&files).Error
	if err != nil {
		return nil, translateError(err, "file")
	}
	return files, nil
}

func (r *fileRepository) ClaimExpired(ctx context.Context, limit int) ([]*entity.File, error) {
	var files []*entity.File
	err := r.getDB(ctx).Raw(`
		UPDATE files
		SET claimed_at = NOW(), updated_at = NOW()
		WHERE id IN (
			SELECT id FROM files
			WHERE expires_at <= NOW()
			  AND deleted_at IS NULL
			  AND claimed_at IS NULL
			  AND failed_delete_attempts < 5
			ORDER BY expires_at ASC
			LIMIT ?
			FOR UPDATE SKIP LOCKED
		)
		RETURNING *
	`, limit).Scan(&files).Error
	if err != nil {
		return nil, translateError(err, "file")
	}
	return files, nil
}

func (r *fileRepository) ReleaseStaleClaimsOlderThan(ctx context.Context, age time.Duration) error {
	cutoff := time.Now().Add(-age)
	return r.getDB(ctx).Model(&entity.File{}).
		Where("claimed_at < ? AND deleted_at IS NULL", cutoff).
		Updates(map[string]interface{}{
			"claimed_at": nil,
			"updated_at": time.Now(),
		}).Error
}

func (r *fileRepository) IncrementFailedAttempts(ctx context.Context, id uuid.UUID) error {
	err := r.getDB(ctx).Model(&entity.File{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"failed_delete_attempts": gorm.Expr("failed_delete_attempts + 1"),
			"updated_at":             time.Now(),
		}).Error
	if err != nil {
		return translateError(err, "file")
	}
	return nil
}

func (r *fileRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {

	err := r.getDB(ctx).Where("id = ?", id).Delete(&entity.File{}).Error
	if err != nil {
		return translateError(err, "file")
	}
	return nil
}
