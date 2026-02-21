package contract

import (
	"context"
	"time"
)

type FileCleanupStorage interface {
	DeleteObject(ctx context.Context, bucket, key string) error
	PresignGetURL(ctx context.Context, bucket, key string, ttl time.Duration) (string, error)
}
