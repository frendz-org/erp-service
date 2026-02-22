package files

import (
	"context"
	"time"

	"erp-service/entity"

	"github.com/google/uuid"
)

type FileRepository interface {
	ClaimExpired(ctx context.Context, limit int) ([]*entity.File, error)
	ReleaseStaleClaimsOlderThan(ctx context.Context, age time.Duration) error
	IncrementFailedAttempts(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type FileStorageAdapter interface {
	DeleteFile(ctx context.Context, bucket, objectKey string) error
}
