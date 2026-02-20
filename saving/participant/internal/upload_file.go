package internal

import (
	"context"
	"fmt"
	"io"
	"time"

	"iam-service/saving/participant/participantdto"
)

func (uc *usecase) UploadFile(ctx context.Context, req *participantdto.UploadFileRequest, file io.Reader, fileSize int64, contentType, filename string) (*participantdto.FileUploadResponse, error) {
	participant, err := uc.participantRepo.GetByID(ctx, req.ParticipantID)
	if err != nil {
		return nil, fmt.Errorf("get participant: %w", err)
	}

	if err := validateParticipantOwnership(participant, req.TenantID, req.ProductID); err != nil {
		return nil, err
	}

	if err := validateEditableState(participant); err != nil {
		return nil, err
	}

	objectKey := generateObjectKey(req.TenantID, req.ProductID, req.ParticipantID, req.FieldName, filename)

	bucket := uc.cfg.Infra.Minio.Bucket
	filePath, err := uc.fileStorage.UploadFile(ctx, bucket, objectKey, file, fileSize, contentType)
	if err != nil {
		return nil, fmt.Errorf("upload file to storage: %w", err)
	}

	presignedURL, err := uc.fileStorage.GetPresignedURL(ctx, bucket, filePath, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("generate presigned URL: %w", err)
	}

	return &participantdto.FileUploadResponse{
		FilePath: filePath,
		URL:      presignedURL,
	}, nil
}
