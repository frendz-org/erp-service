package files

import (
	"context"
	"fmt"

	"erp-service/entity"

	"go.uber.org/zap"
)

func (uc *usecase) ProcessFile(ctx context.Context, file *entity.File) error {
	if err := uc.fileStorage.DeleteFile(ctx, file.Bucket, file.StorageKey); err != nil {
		uc.logger.Warn("failed to delete file from storage",
			zap.String("file_id", file.ID.String()),
			zap.String("bucket", file.Bucket),
			zap.String("storage_key", file.StorageKey),
			zap.Error(err),
		)

		if incrErr := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
			return uc.fileRepo.IncrementFailedAttempts(txCtx, file.ID)
		}); incrErr != nil {
			uc.logger.Error("failed to increment failed attempts",
				zap.String("file_id", file.ID.String()),
				zap.Error(incrErr),
			)
		}

		return fmt.Errorf("delete from storage: %w", err)
	}

	if err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return uc.fileRepo.SoftDelete(txCtx, file.ID)
	}); err != nil {
		uc.logger.Error("failed to soft-delete file record",
			zap.String("file_id", file.ID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("soft-delete file record: %w", err)
	}

	return nil
}
