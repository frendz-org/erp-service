package worker

import (
	"context"
	"io"
	"testing"
	"time"

	"erp-service/entity"
	"erp-service/saving/participant/contract"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type mockFileRepo struct{ mock.Mock }

var _ contract.FileRepository = (*mockFileRepo)(nil)

func (m *mockFileRepo) Create(ctx context.Context, f *entity.File) error {
	return m.Called(ctx, f).Error(0)
}
func (m *mockFileRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.File, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.File), args.Error(1)
}
func (m *mockFileRepo) SetPermanent(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockFileRepo) SetExpiring(ctx context.Context, id uuid.UUID, expiry time.Time) error {
	return m.Called(ctx, id, expiry).Error(0)
}
func (m *mockFileRepo) ListExpired(ctx context.Context, limit int) ([]*entity.File, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.File), args.Error(1)
}
func (m *mockFileRepo) ListExpiredForUpdate(ctx context.Context, limit int) ([]*entity.File, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.File), args.Error(1)
}
func (m *mockFileRepo) IncrementFailedAttempts(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockFileRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type mockFileStorage struct{ mock.Mock }

var _ contract.FileStorageAdapter = (*mockFileStorage)(nil)

func (m *mockFileStorage) UploadFile(ctx context.Context, bucket, key string, _ io.Reader, size int64, contentType string) (string, error) {
	return "", nil
}
func (m *mockFileStorage) DeleteFile(ctx context.Context, bucket, key string) error {
	return m.Called(ctx, bucket, key).Error(0)
}
func (m *mockFileStorage) GetPresignedURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error) {
	return "", nil
}

type mockTxManager struct{ mock.Mock }

var _ contract.TransactionManager = (*mockTxManager)(nil)

func (m *mockTxManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)

	if args.Get(0) == nil && fn != nil {
		return fn(ctx)
	}
	return args.Error(0)
}

func newTestWorker(fileRepo contract.FileRepository, fileStorage contract.FileStorageAdapter, txMgr contract.TransactionManager) *Worker {
	w := NewWorker(fileRepo, fileStorage, txMgr, zap.NewNop())

	w.interval = 10 * time.Millisecond
	return w
}

func makeExpiredFile(bucket, storageKey string) *entity.File {
	now := time.Now().Add(-1 * time.Hour)
	return &entity.File{
		ID:         uuid.New(),
		Bucket:     bucket,
		StorageKey: storageKey,
		ExpiresAt:  &now,
		DeletedAt:  gorm.DeletedAt{},
	}
}

func TestWorker_RunCleanupBatch_Success(t *testing.T) {
	fileRepo := new(mockFileRepo)
	fileStorage := new(mockFileStorage)
	txMgr := new(mockTxManager)

	file := makeExpiredFile("participants", "participants/path/photo.jpg")

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	fileRepo.On("ListExpiredForUpdate", mock.Anything, defaultBatchSize).Return([]*entity.File{file}, nil)
	fileStorage.On("DeleteFile", mock.Anything, "participants", "participants/path/photo.jpg").Return(nil)
	fileRepo.On("SoftDelete", mock.Anything, file.ID).Return(nil)

	w := newTestWorker(fileRepo, fileStorage, txMgr)
	w.runCleanupBatch(context.Background())

	fileStorage.AssertCalled(t, "DeleteFile", mock.Anything, "participants", "participants/path/photo.jpg")
	fileRepo.AssertCalled(t, "SoftDelete", mock.Anything, file.ID)

	fileRepo.AssertCalled(t, "ListExpiredForUpdate", mock.Anything, defaultBatchSize)
}

func TestWorker_RunCleanupBatch_StorageDeleteFails_IncrementsFailedAttempts(t *testing.T) {
	fileRepo := new(mockFileRepo)
	fileStorage := new(mockFileStorage)
	txMgr := new(mockTxManager)

	file := makeExpiredFile("participants", "participants/path/photo.jpg")

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	fileRepo.On("ListExpiredForUpdate", mock.Anything, defaultBatchSize).Return([]*entity.File{file}, nil)
	fileStorage.On("DeleteFile", mock.Anything, "participants", "participants/path/photo.jpg").
		Return(assert.AnError)
	fileRepo.On("IncrementFailedAttempts", mock.Anything, file.ID).Return(nil)

	w := newTestWorker(fileRepo, fileStorage, txMgr)
	w.runCleanupBatch(context.Background())

	fileRepo.AssertCalled(t, "IncrementFailedAttempts", mock.Anything, file.ID)

	fileRepo.AssertNotCalled(t, "SoftDelete", mock.Anything, file.ID)
}

func TestWorker_RunCleanupBatch_NoExpiredFiles(t *testing.T) {
	fileRepo := new(mockFileRepo)
	fileStorage := new(mockFileStorage)
	txMgr := new(mockTxManager)

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	fileRepo.On("ListExpiredForUpdate", mock.Anything, defaultBatchSize).Return([]*entity.File{}, nil)

	w := newTestWorker(fileRepo, fileStorage, txMgr)
	w.runCleanupBatch(context.Background())

	fileStorage.AssertNotCalled(t, "DeleteFile", mock.Anything, mock.Anything, mock.Anything)
	fileRepo.AssertNotCalled(t, "SoftDelete", mock.Anything, mock.Anything)
}

func TestWorker_StartStop(t *testing.T) {
	fileRepo := new(mockFileRepo)
	fileStorage := new(mockFileStorage)
	txMgr := new(mockTxManager)

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	fileRepo.On("ListExpiredForUpdate", mock.Anything, mock.Anything).Return([]*entity.File{}, nil)

	w := newTestWorker(fileRepo, fileStorage, txMgr)
	w.interval = 1 * time.Hour

	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)

	cancel()

	done := make(chan struct{})
	go func() {
		w.Stop()
		close(done)
	}()

	select {
	case <-done:

	case <-time.After(2 * time.Second):
		t.Fatal("worker did not stop within timeout")
	}
}

func TestWorker_MultipleFiles_BatchProcessed(t *testing.T) {
	fileRepo := new(mockFileRepo)
	fileStorage := new(mockFileStorage)
	txMgr := new(mockTxManager)

	file1 := makeExpiredFile("participants", "participants/path/a.jpg")
	file2 := makeExpiredFile("participants", "participants/path/b.jpg")

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	fileRepo.On("ListExpiredForUpdate", mock.Anything, defaultBatchSize).Return([]*entity.File{file1, file2}, nil)
	fileStorage.On("DeleteFile", mock.Anything, "participants", "participants/path/a.jpg").Return(nil)
	fileStorage.On("DeleteFile", mock.Anything, "participants", "participants/path/b.jpg").Return(nil)
	fileRepo.On("SoftDelete", mock.Anything, file1.ID).Return(nil)
	fileRepo.On("SoftDelete", mock.Anything, file2.ID).Return(nil)

	w := newTestWorker(fileRepo, fileStorage, txMgr)
	w.runCleanupBatch(context.Background())

	fileRepo.AssertCalled(t, "SoftDelete", mock.Anything, file1.ID)
	fileRepo.AssertCalled(t, "SoftDelete", mock.Anything, file2.ID)
	assert.Equal(t, 2, len(fileStorage.Calls))
}

func TestWorker_Start_DoubleCall_OnlyStartsOnce(t *testing.T) {
	fileRepo := new(mockFileRepo)
	fileStorage := new(mockFileStorage)
	txMgr := new(mockTxManager)

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	fileRepo.On("ListExpiredForUpdate", mock.Anything, mock.Anything).Return([]*entity.File{}, nil)

	w := newTestWorker(fileRepo, fileStorage, txMgr)
	w.interval = 1 * time.Hour

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Start(ctx)
	w.Start(ctx)

	cancel()

	done := make(chan struct{})
	go func() {
		w.Stop()
		close(done)
	}()

	select {
	case <-done:

	case <-time.After(2 * time.Second):
		t.Fatal("worker did not stop within timeout")
	}
}
