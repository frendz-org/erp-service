package participant

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"

	"go.uber.org/zap"
)

const defaultBucket = "participants"

func (uc *usecase) UploadFile(ctx context.Context, req *UploadFileRequest) (*FileUploadResponse, error) {
	participant, err := uc.participantRepo.GetByID(ctx, req.ParticipantID)
	if err != nil {
		return nil, fmt.Errorf("get participant: %w", err)
	}

	if err := ValidateParticipantOwnership(participant, req.TenantID, req.ProductID); err != nil {
		return nil, err
	}

	if err := ValidateEditableState(participant); err != nil {
		return nil, err
	}

	bucket := defaultBucket
	objectKey := GenerateObjectKey(req.TenantID, req.ProductID, req.ParticipantID, req.FieldName, req.FileName)

	storageKey, err := uc.fileStorage.UploadFile(ctx, bucket, objectKey, req.Reader, req.Size, req.ContentType)
	if err != nil {
		return nil, fmt.Errorf("upload to storage: %w", err)
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	file := &entity.File{
		TenantID:     req.TenantID,
		ProductID:    req.ProductID,
		UploadedBy:   req.UploadedBy,
		Bucket:       bucket,
		StorageKey:   storageKey,
		OriginalName: req.FileName,
		ContentType:  req.ContentType,
		SizeBytes:    req.Size,
		ExpiresAt:    &expiresAt,
	}

	if err := uc.fileRepo.Create(ctx, file); err != nil {

		if delErr := uc.fileStorage.DeleteFile(ctx, bucket, storageKey); delErr != nil {
			uc.logger.Warn("failed to clean up orphaned storage object after DB insert failure",
				zap.String("bucket", bucket),
				zap.String("storage_key", storageKey),
				zap.Error(delErr),
			)
		}
		return nil, fmt.Errorf("persist file metadata: %w", err)
	}

	return &FileUploadResponse{
		FileID: file.ID,
	}, nil
}
