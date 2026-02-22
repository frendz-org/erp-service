package worker

import (
	"context"
	"sync"
	"time"

	"erp-service/files"

	"go.uber.org/zap"
)

const defaultCleanupInterval = 5 * time.Minute

type Worker struct {
	uc        files.Usecase
	logger    *zap.Logger
	interval  time.Duration
	done      chan struct{}
	startOnce sync.Once
}

func NewWorker(uc files.Usecase, logger *zap.Logger) *Worker {
	return &Worker{
		uc:       uc,
		logger:   logger,
		interval: defaultCleanupInterval,
		done:     make(chan struct{}),
	}
}

func (w *Worker) SetInterval(d time.Duration) { w.interval = d }

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
					w.runOnce(ctx)
				}
			}
		}()
	})
}

func (w *Worker) Stop() {
	<-w.done
}

func (w *Worker) runOnce(ctx context.Context) {
	result, err := w.uc.CleanupBatch(ctx)
	if err != nil {
		w.logger.Error("cleanup batch failed", zap.Error(err))
		return
	}
	if result.Processed > 0 || result.Failed > 0 {
		w.logger.Info("cleanup batch completed",
			zap.Int("processed", result.Processed),
			zap.Int("failed", result.Failed),
		)
	}
}
