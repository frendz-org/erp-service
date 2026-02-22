package files

import (
	"context"

	"erp-service/entity"

	"go.uber.org/zap"
)

type BatchResult struct {
	Processed int
	Failed    int
}

type Usecase interface {
	CleanupBatch(ctx context.Context) (BatchResult, error)
	ProcessFile(ctx context.Context, file *entity.File) error
}

func NewUsecase(
	fileRepo FileRepository,
	fileStorage FileStorageAdapter,
	txManager TransactionManager,
	logger *zap.Logger,
	cfg Config,
) Usecase {
	return &usecase{
		fileRepo:    fileRepo,
		fileStorage: fileStorage,
		txManager:   txManager,
		logger:      logger,
		cfg:         cfg,
	}
}
