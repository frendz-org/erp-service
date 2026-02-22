package files_test

import (
	"context"
	"testing"
	"time"

	"erp-service/entity"
	"erp-service/files"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type mockFileRepo struct{ mock.Mock }

var _ files.FileRepository = (*mockFileRepo)(nil)

func (m *mockFileRepo) ClaimExpired(ctx context.Context, limit int) ([]*entity.File, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.File), args.Error(1)
}

func (m *mockFileRepo) ReleaseStaleClaimsOlderThan(ctx context.Context, age time.Duration) error {
	return m.Called(ctx, age).Error(0)
}

func (m *mockFileRepo) IncrementFailedAttempts(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockFileRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type mockFileStorage struct{ mock.Mock }

var _ files.FileStorageAdapter = (*mockFileStorage)(nil)

func (m *mockFileStorage) DeleteFile(ctx context.Context, bucket, objectKey string) error {
	return m.Called(ctx, bucket, objectKey).Error(0)
}

type mockTxManager struct{ mock.Mock }

var _ files.TransactionManager = (*mockTxManager)(nil)

func (m *mockTxManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) == nil && fn != nil {
		return fn(ctx)
	}
	return args.Error(0)
}

func newUC(repo *mockFileRepo, storage *mockFileStorage, tx *mockTxManager) files.Usecase {
	return files.NewUsecase(repo, storage, tx, zap.NewNop(), files.DefaultConfig())
}

func makeFile(bucket, storageKey string) *entity.File {
	exp := time.Now().Add(-1 * time.Hour)
	return &entity.File{
		ID:         uuid.New(),
		Bucket:     bucket,
		StorageKey: storageKey,
		ExpiresAt:  &exp,
		DeletedAt:  gorm.DeletedAt{},
	}
}

const (
	testBucket     = "participants"
	testStorageKey = "participants/path/photo.jpg"
)

func TestProcessFile_Success(t *testing.T) {
	repo, storage, tx := new(mockFileRepo), new(mockFileStorage), new(mockTxManager)
	file := makeFile(testBucket, testStorageKey)

	tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	storage.On("DeleteFile", mock.Anything, testBucket, testStorageKey).Return(nil)
	repo.On("SoftDelete", mock.Anything, file.ID).Return(nil)

	err := newUC(repo, storage, tx).ProcessFile(context.Background(), file)

	require.NoError(t, err)
	storage.AssertCalled(t, "DeleteFile", mock.Anything, testBucket, testStorageKey)
	repo.AssertCalled(t, "SoftDelete", mock.Anything, file.ID)
	repo.AssertNotCalled(t, "IncrementFailedAttempts", mock.Anything, mock.Anything)
}

func TestProcessFile_StorageDeleteFails_IncrementsCounterAndReturnsError(t *testing.T) {
	repo, storage, tx := new(mockFileRepo), new(mockFileStorage), new(mockTxManager)
	file := makeFile(testBucket, testStorageKey)

	tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	storage.On("DeleteFile", mock.Anything, testBucket, testStorageKey).Return(assert.AnError)
	repo.On("IncrementFailedAttempts", mock.Anything, file.ID).Return(nil)

	err := newUC(repo, storage, tx).ProcessFile(context.Background(), file)

	require.Error(t, err)
	repo.AssertCalled(t, "IncrementFailedAttempts", mock.Anything, file.ID)
	repo.AssertNotCalled(t, "SoftDelete", mock.Anything, mock.Anything)
}

func TestProcessFile_StorageDeleteFails_IncrementAlsoFails_StillReturnsStorageError(t *testing.T) {
	repo, storage, tx := new(mockFileRepo), new(mockFileStorage), new(mockTxManager)
	file := makeFile(testBucket, testStorageKey)

	tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	storage.On("DeleteFile", mock.Anything, testBucket, testStorageKey).Return(assert.AnError)
	repo.On("IncrementFailedAttempts", mock.Anything, file.ID).Return(assert.AnError)

	err := newUC(repo, storage, tx).ProcessFile(context.Background(), file)

	require.Error(t, err)
	assert.ErrorContains(t, err, "delete from storage")
}

func TestProcessFile_SoftDeleteFails_ReturnsError(t *testing.T) {
	repo, storage, tx := new(mockFileRepo), new(mockFileStorage), new(mockTxManager)
	file := makeFile(testBucket, testStorageKey)

	tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	storage.On("DeleteFile", mock.Anything, testBucket, testStorageKey).Return(nil)
	repo.On("SoftDelete", mock.Anything, file.ID).Return(assert.AnError)

	err := newUC(repo, storage, tx).ProcessFile(context.Background(), file)

	require.Error(t, err)
	assert.ErrorContains(t, err, "soft-delete file record")
}

func TestCleanupBatch_Success_SingleFile(t *testing.T) {
	repo, storage, tx := new(mockFileRepo), new(mockFileStorage), new(mockTxManager)
	file := makeFile(testBucket, testStorageKey)

	tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	repo.On("ReleaseStaleClaimsOlderThan", mock.Anything, files.DefaultConfig().StaleClaimAge).Return(nil)
	repo.On("ClaimExpired", mock.Anything, files.DefaultConfig().BatchSize).Return([]*entity.File{file}, nil)
	storage.On("DeleteFile", mock.Anything, testBucket, testStorageKey).Return(nil)
	repo.On("SoftDelete", mock.Anything, file.ID).Return(nil)

	result, err := newUC(repo, storage, tx).CleanupBatch(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 1, result.Processed)
	assert.Equal(t, 0, result.Failed)
}

func TestCleanupBatch_NoExpiredFiles_ReturnsZeroCounts(t *testing.T) {
	repo, storage, tx := new(mockFileRepo), new(mockFileStorage), new(mockTxManager)

	tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	repo.On("ReleaseStaleClaimsOlderThan", mock.Anything, mock.Anything).Return(nil)
	repo.On("ClaimExpired", mock.Anything, mock.Anything).Return([]*entity.File{}, nil)

	result, err := newUC(repo, storage, tx).CleanupBatch(context.Background())

	require.NoError(t, err)
	assert.Equal(t, files.BatchResult{}, result)
	storage.AssertNotCalled(t, "DeleteFile", mock.Anything, mock.Anything, mock.Anything)
}

func TestCleanupBatch_ClaimExpiredFails_ReturnsFatalError(t *testing.T) {
	repo, storage, tx := new(mockFileRepo), new(mockFileStorage), new(mockTxManager)

	tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	repo.On("ReleaseStaleClaimsOlderThan", mock.Anything, mock.Anything).Return(nil)
	repo.On("ClaimExpired", mock.Anything, mock.Anything).Return(nil, assert.AnError)

	result, err := newUC(repo, storage, tx).CleanupBatch(context.Background())

	require.Error(t, err)
	assert.Equal(t, files.BatchResult{}, result)
	storage.AssertNotCalled(t, "DeleteFile", mock.Anything, mock.Anything, mock.Anything)
}

func TestCleanupBatch_ReleaseStaleClaimsFails_NonFatal_BatchContinues(t *testing.T) {
	repo, storage, tx := new(mockFileRepo), new(mockFileStorage), new(mockTxManager)
	file := makeFile(testBucket, testStorageKey)

	tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)

	repo.On("ReleaseStaleClaimsOlderThan", mock.Anything, mock.Anything).Return(assert.AnError)
	repo.On("ClaimExpired", mock.Anything, mock.Anything).Return([]*entity.File{file}, nil)
	storage.On("DeleteFile", mock.Anything, testBucket, testStorageKey).Return(nil)
	repo.On("SoftDelete", mock.Anything, file.ID).Return(nil)

	result, err := newUC(repo, storage, tx).CleanupBatch(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 1, result.Processed)
}

func TestCleanupBatch_MultipleFiles_CountsCorrectly(t *testing.T) {
	repo, storage, tx := new(mockFileRepo), new(mockFileStorage), new(mockTxManager)
	good := makeFile(testBucket, "participants/a.jpg")
	bad := makeFile(testBucket, "participants/b.jpg")

	tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	repo.On("ReleaseStaleClaimsOlderThan", mock.Anything, mock.Anything).Return(nil)
	repo.On("ClaimExpired", mock.Anything, mock.Anything).Return([]*entity.File{good, bad}, nil)

	storage.On("DeleteFile", mock.Anything, testBucket, "participants/a.jpg").Return(nil)
	repo.On("SoftDelete", mock.Anything, good.ID).Return(nil)

	storage.On("DeleteFile", mock.Anything, testBucket, "participants/b.jpg").Return(assert.AnError)
	repo.On("IncrementFailedAttempts", mock.Anything, bad.ID).Return(nil)

	result, err := newUC(repo, storage, tx).CleanupBatch(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 1, result.Processed)
	assert.Equal(t, 1, result.Failed)
}

func TestCleanupBatch_UsesConfigValues(t *testing.T) {
	repo, storage, tx := new(mockFileRepo), new(mockFileStorage), new(mockTxManager)

	customCfg := files.Config{BatchSize: 10, StaleClaimAge: 15 * time.Minute}
	uc := files.NewUsecase(repo, storage, tx, zap.NewNop(), customCfg)

	tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	repo.On("ReleaseStaleClaimsOlderThan", mock.Anything, 15*time.Minute).Return(nil)
	repo.On("ClaimExpired", mock.Anything, 10).Return([]*entity.File{}, nil)

	_, err := uc.CleanupBatch(context.Background())

	require.NoError(t, err)

	repo.AssertCalled(t, "ReleaseStaleClaimsOlderThan", mock.Anything, 15*time.Minute)
	repo.AssertCalled(t, "ClaimExpired", mock.Anything, 10)
}
