package files

import (
	"context"

	"erp-service/entity"

	"go.uber.org/zap"
)

func (uc *usecase) CleanupBatch(ctx context.Context) (BatchResult, error) {
	if err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return uc.fileRepo.ReleaseStaleClaimsOlderThan(txCtx, uc.cfg.StaleClaimAge)
	}); err != nil {
		uc.logger.Warn("failed to release stale file claims", zap.Error(err))
	}

	var files []*entity.File
	if err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		files, err = uc.fileRepo.ClaimExpired(txCtx, uc.cfg.BatchSize)
		return err
	}); err != nil {
		uc.logger.Error("failed to claim expired files", zap.Error(err))
		return BatchResult{}, err
	}

	var result BatchResult
	for _, file := range files {
		if err := uc.ProcessFile(ctx, file); err != nil {
			uc.logger.Warn("failed to process file",
				zap.String("file_id", file.ID.String()),
				zap.Error(err),
			)
			result.Failed++
			continue
		}
		result.Processed++
	}

	return result, nil
}
