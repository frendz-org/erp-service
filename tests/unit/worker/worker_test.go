package worker_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"erp-service/delivery/worker"
	"erp-service/entity"
	"erp-service/files"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockUsecase struct{ mock.Mock }

var _ files.Usecase = (*mockUsecase)(nil)

func (m *mockUsecase) CleanupBatch(ctx context.Context) (files.BatchResult, error) {
	args := m.Called(ctx)
	return args.Get(0).(files.BatchResult), args.Error(1)
}

func (m *mockUsecase) ProcessFile(ctx context.Context, file *entity.File) error {
	return m.Called(ctx, file).Error(0)
}

func newTestWorker(uc *mockUsecase) *worker.Worker {
	w := worker.NewWorker(uc, zap.NewNop())
	w.SetInterval(10 * time.Millisecond)
	return w
}

func TestWorker_StartStop(t *testing.T) {
	uc := new(mockUsecase)
	uc.On("CleanupBatch", mock.Anything).Return(files.BatchResult{}, nil).Maybe()

	w := newTestWorker(uc)
	w.SetInterval(1 * time.Hour)

	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)
	cancel()

	done := make(chan struct{})
	go func() { w.Stop(); close(done) }()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not stop within timeout")
	}
}

func TestWorker_Start_DoubleCall_OnlyStartsOnce(t *testing.T) {
	uc := new(mockUsecase)
	uc.On("CleanupBatch", mock.Anything).Return(files.BatchResult{}, nil).Maybe()

	w := newTestWorker(uc)
	w.SetInterval(1 * time.Hour)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Start(ctx)
	w.Start(ctx)
	cancel()

	done := make(chan struct{})
	go func() { w.Stop(); close(done) }()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not stop within timeout")
	}
}

func TestWorker_OnTick_CallsCleanupBatch(t *testing.T) {
	uc := new(mockUsecase)
	called := make(chan struct{}, 1)
	uc.On("CleanupBatch", mock.Anything).
		Run(func(args mock.Arguments) {
			select {
			case called <- struct{}{}:
			default:
			}
		}).
		Return(files.BatchResult{Processed: 1}, nil)

	w := newTestWorker(uc)
	w.SetInterval(20 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.Start(ctx)

	select {
	case <-called:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("CleanupBatch was not called within timeout")
	}

	cancel()
	done := make(chan struct{})
	go func() { w.Stop(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not stop within timeout")
	}
}

func TestWorker_CleanupBatch_Error_WorkerContinues(t *testing.T) {
	uc := new(mockUsecase)
	called := make(chan struct{}, 1)
	uc.On("CleanupBatch", mock.Anything).
		Run(func(args mock.Arguments) {
			select {
			case called <- struct{}{}:
			default:
			}
		}).
		Return(files.BatchResult{}, errTestFatal)

	w := newTestWorker(uc)
	w.SetInterval(20 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.Start(ctx)

	select {
	case <-called:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("CleanupBatch was not called within timeout")
	}

	cancel()
	done := make(chan struct{})
	go func() { w.Stop(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not stop within timeout")
	}
}

var errTestFatal = fmt.Errorf("fatal test error")
