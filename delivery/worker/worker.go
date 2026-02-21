package worker

import (
	"context"
	"sync"
	"time"

	"erp-service/entity"
	"erp-service/saving/participant/contract"

	"go.uber.org/zap"
)

const (
	defaultCleanupInterval = 5 * time.Minute
	defaultBatchSize       = 50
)

type Worker struct {
	fileRepo    contract.FileRepository
	fileStorage contract.FileStorageAdapter
	txManager   contract.TransactionManager
	logger      *zap.Logger
	interval    time.Duration
	batchSize   int
	done        chan struct{}
	startOnce   sync.Once
}

func NewWorker(
	fileRepo contract.FileRepository,
	fileStorage contract.FileStorageAdapter,
	txManager contract.TransactionManager,
	log *zap.Logger,
) *Worker {
	return &Worker{
		fileRepo:    fileRepo,
		fileStorage: fileStorage,
		txManager:   txManager,
		logger:      log,
		interval:    defaultCleanupInterval,
		batchSize:   defaultBatchSize,
		done:        make(chan struct{}),
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.startOnce.Do(func() {
		go func() {
			defer close(w.done)
			ticker := time.NewTicker(w.interval)
			defer ticker.Stop()
			w.logger.Info("file cleanup worker started", zap.Duration("interval", w.interval))
			for {
				select {
				case <-ctx.Done():
					w.logger.Info("file cleanup worker stopping")
					return
				case <-ticker.C:
					if ctx.Err() != nil {
						return
					}
					w.runCleanupBatch(ctx)
				}
			}
		}()
	})
}

func (w *Worker) Stop() {
	<-w.done
}

func (w *Worker) runCleanupBatch(ctx context.Context) {

	var files []*entity.File
	err := w.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var listErr error
		files, listErr = w.fileRepo.ListExpiredForUpdate(txCtx, w.batchSize)
		return listErr
	})
	if err != nil {
		w.logger.Error("failed to list expired files", zap.Error(err))
		return
	}

	for _, file := range files {
		w.processOneFile(ctx, file)
	}
}

func (w *Worker) processOneFile(ctx context.Context, file *entity.File) {

	if err := w.fileStorage.DeleteFile(ctx, file.Bucket, file.StorageKey); err != nil {
		w.logger.Warn("failed to delete from storage",
			zap.String("file_id", file.ID.String()),
			zap.Error(err),
		)

		if incrErr := w.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
			return w.fileRepo.IncrementFailedAttempts(txCtx, file.ID)
		}); incrErr != nil {
			w.logger.Error("failed to increment failed attempts",
				zap.String("file_id", file.ID.String()),
				zap.Error(incrErr),
			)
		}
		return
	}

	if err := w.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return w.fileRepo.SoftDelete(txCtx, file.ID)
	}); err != nil {
		w.logger.Error("failed to soft-delete file",
			zap.String("file_id", file.ID.String()),
			zap.Error(err),
		)
	}
}
