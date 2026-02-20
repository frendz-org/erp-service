package contract

import (
	"context"
	"io"
	"time"
)

type FileStorageAdapter interface {
	UploadFile(ctx context.Context, bucket, objectKey string, data io.Reader, size int64, contentType string) (string, error)
	DeleteFile(ctx context.Context, bucket, objectKey string) error
	GetPresignedURL(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error)
}
