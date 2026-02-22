package participant_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"erp-service/config"
	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func makeUploadUsecase() (participant.Usecase, *MockParticipantRepository, *MockFileRepository, *MockFileStorageAdapter) {
	participantRepo := new(MockParticipantRepository)
	fileRepo := new(MockFileRepository)
	fileStorage := new(MockFileStorageAdapter)

	uc := participant.NewUsecase(
		&config.Config{},
		zap.NewNop(),
		new(MockTransactionManager),
		participantRepo,
		new(MockParticipantIdentityRepository),
		new(MockParticipantAddressRepository),
		new(MockParticipantBankAccountRepository),
		new(MockParticipantFamilyMemberRepository),
		new(MockParticipantEmploymentRepository),
		new(MockParticipantPensionRepository),
		new(MockParticipantBeneficiaryRepository),
		new(MockParticipantStatusHistoryRepository),
		fileStorage,
		fileRepo,
		nil, nil, nil, nil, nil, nil,
	)
	return uc, participantRepo, fileRepo, fileStorage
}

func TestUploadFile_Success_ReturnsFileID(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()
	uploaderID := uuid.New()

	uc, participantRepo, fileRepo, fileStorage := makeUploadUsecase()

	p := &entity.Participant{
		ID:        participantID,
		TenantID:  tenantID,
		ProductID: productID,
		Status:    entity.ParticipantStatusDraft,
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	participantRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)

	fileData := []byte("fake-image-data")
	storageKey := "participants/tenant/product/participant/ktp_photo/photo.jpg"
	fileStorage.On("UploadFile", mock.Anything, "participants", mock.AnythingOfType("string"),
		mock.Anything, int64(len(fileData)), "image/jpeg").
		Return(storageKey, nil)

	var createdFile *entity.File
	fileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.File")).
		Run(func(args mock.Arguments) {
			createdFile = args.Get(1).(*entity.File)
			createdFile.ID = uuid.New()
		}).
		Return(nil)

	req := &participant.UploadFileRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UploadedBy:    uploaderID,
		FileName:      "photo.jpg",
		ContentType:   "image/jpeg",
		Reader:        bytes.NewReader(fileData),
		Size:          int64(len(fileData)),
		FieldName:     "ktp_photo",
	}

	result, err := uc.UploadFile(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.NotEqual(t, uuid.Nil, result.FileID)

	fileStorage.AssertCalled(t, "UploadFile", mock.Anything, "participants", mock.AnythingOfType("string"),
		mock.Anything, int64(len(fileData)), "image/jpeg")

	require.NotNil(t, createdFile)
	assert.Equal(t, "participants", createdFile.Bucket)
	assert.Equal(t, storageKey, createdFile.StorageKey)

	require.NotNil(t, createdFile.ExpiresAt)
	duration := createdFile.ExpiresAt.Sub(time.Now())
	assert.True(t, duration > 23*time.Hour && duration < 25*time.Hour,
		"expires_at should be ~24h from now, got %v", duration)
}

func TestUploadFile_ParticipantNotFound(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()
	uploaderID := uuid.New()

	uc, participantRepo, _, _ := makeUploadUsecase()

	participantRepo.On("GetByID", mock.Anything, participantID).
		Return(nil, errors.ErrNotFound("participant not found"))

	req := &participant.UploadFileRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UploadedBy:    uploaderID,
		FileName:      "photo.jpg",
		ContentType:   "image/jpeg",
		Reader:        bytes.NewReader([]byte("data")),
		Size:          4,
	}

	result, err := uc.UploadFile(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.IsNotFound(err))
}

func TestUploadFile_WrongTenant(t *testing.T) {
	tenantID := uuid.New()
	otherTenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()
	uploaderID := uuid.New()

	uc, participantRepo, _, _ := makeUploadUsecase()

	p := &entity.Participant{
		ID:        participantID,
		TenantID:  otherTenantID,
		ProductID: productID,
		Status:    entity.ParticipantStatusDraft,
	}

	participantRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)

	req := &participant.UploadFileRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UploadedBy:    uploaderID,
		FileName:      "photo.jpg",
		ContentType:   "image/jpeg",
		Reader:        bytes.NewReader([]byte("data")),
		Size:          4,
	}

	result, err := uc.UploadFile(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestUploadFile_NonEditableParticipant(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()
	uploaderID := uuid.New()

	uc, participantRepo, _, _ := makeUploadUsecase()

	p := &entity.Participant{
		ID:        participantID,
		TenantID:  tenantID,
		ProductID: productID,
		Status:    entity.ParticipantStatusApproved,
	}

	participantRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)

	req := &participant.UploadFileRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UploadedBy:    uploaderID,
		FileName:      "photo.jpg",
		ContentType:   "image/jpeg",
		Reader:        bytes.NewReader([]byte("data")),
		Size:          4,
	}

	result, err := uc.UploadFile(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestUploadFile_StoredWithUploaderID(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()
	uploaderID := uuid.New()

	uc, participantRepo, fileRepo, fileStorage := makeUploadUsecase()

	p := &entity.Participant{
		ID:        participantID,
		TenantID:  tenantID,
		ProductID: productID,
		Status:    entity.ParticipantStatusDraft,
		Version:   1,
	}

	participantRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)
	fileStorage.On("UploadFile", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).
		Return("participants/path/doc.pdf", nil)

	var createdFile *entity.File
	fileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.File")).
		Run(func(args mock.Arguments) {
			createdFile = args.Get(1).(*entity.File)
		}).
		Return(nil)

	pdfData := []byte("pdf-data")
	req := &participant.UploadFileRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UploadedBy:    uploaderID,
		FileName:      "doc.pdf",
		ContentType:   "application/pdf",
		Reader:        bytes.NewReader(pdfData),
		Size:          int64(len(pdfData)),
		FieldName:     "supporting_doc",
	}

	_, err := uc.UploadFile(context.Background(), req)
	require.NoError(t, err)

	require.NotNil(t, createdFile)
	assert.Equal(t, uploaderID, createdFile.UploadedBy)
	assert.Equal(t, tenantID, createdFile.TenantID)
	assert.Equal(t, productID, createdFile.ProductID)
	assert.Equal(t, "doc.pdf", createdFile.OriginalName)
	assert.Equal(t, "application/pdf", createdFile.ContentType)
}

func TestUploadFile_StorageError_NoDatabaseInsert(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()
	uploaderID := uuid.New()

	uc, participantRepo, fileRepo, fileStorage := makeUploadUsecase()

	p := &entity.Participant{
		ID:        participantID,
		TenantID:  tenantID,
		ProductID: productID,
		Status:    entity.ParticipantStatusDraft,
	}

	participantRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)
	fileStorage.On("UploadFile", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).
		Return("", errors.ErrInternal("storage unavailable"))

	req := &participant.UploadFileRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UploadedBy:    uploaderID,
		FileName:      "photo.jpg",
		ContentType:   "image/jpeg",
		Reader:        bytes.NewReader([]byte("data")),
		Size:          4,
	}

	result, err := uc.UploadFile(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)

	fileRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestUploadFile_DBInsertFails_CleansUpStorage(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()
	uploaderID := uuid.New()

	uc, participantRepo, fileRepo, fileStorage := makeUploadUsecase()

	p := &entity.Participant{
		ID:        participantID,
		TenantID:  tenantID,
		ProductID: productID,
		Status:    entity.ParticipantStatusDraft,
	}

	storageKey := "participants/path/photo.jpg"

	participantRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)
	fileStorage.On("UploadFile", mock.Anything, "participants", mock.AnythingOfType("string"),
		mock.Anything, mock.Anything, "image/jpeg").
		Return(storageKey, nil)
	fileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.File")).
		Return(errors.ErrInternal("db insert failed"))
	fileStorage.On("DeleteFile", mock.Anything, "participants", storageKey).Return(nil)

	req := &participant.UploadFileRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UploadedBy:    uploaderID,
		FileName:      "photo.jpg",
		ContentType:   "image/jpeg",
		Reader:        bytes.NewReader([]byte("data")),
		Size:          4,
		FieldName:     "ktp_photo",
	}

	result, err := uc.UploadFile(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)

	fileStorage.AssertCalled(t, "DeleteFile", mock.Anything, "participants", storageKey)
}
