package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"erp-service/saving/participant/contract"

	"github.com/minio/minio-go/v7"
)

type fileStorage struct {
	client *minio.Client
}

func NewFileStorage(client *minio.Client) contract.FileStorageAdapter {
	return &fileStorage{
		client: client,
	}
}

func (fs *fileStorage) UploadFile(ctx context.Context, bucket, objectKey string, data io.Reader, size int64, contentType string) (string, error) {
	if err := fs.ensureBucketExists(ctx, bucket); err != nil {
		return "", fmt.Errorf("ensure bucket exists: %w", err)
	}

	_, err := fs.client.PutObject(ctx, bucket, objectKey, data, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("upload file to MinIO: %w", err)
	}

	return objectKey, nil
}

func (fs *fileStorage) DeleteFile(ctx context.Context, bucket, objectKey string) error {
	err := fs.client.RemoveObject(ctx, bucket, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("delete file from MinIO: %w", err)
	}
	return nil
}

func (fs *fileStorage) DeleteObject(ctx context.Context, bucket, key string) error {
	return fs.DeleteFile(ctx, bucket, key)
}

func (fs *fileStorage) GetPresignedURL(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error) {
	url, err := fs.client.PresignedGetObject(ctx, bucket, objectKey, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("generate presigned URL: %w", err)
	}
	return url.String(), nil
}

func (fs *fileStorage) PresignGetURL(ctx context.Context, bucket, key string, ttl time.Duration) (string, error) {
	return fs.GetPresignedURL(ctx, bucket, key, ttl)
}

func (fs *fileStorage) ensureBucketExists(ctx context.Context, bucket string) error {
	exists, err := fs.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("check bucket existence: %w", err)
	}

	if !exists {
		err = fs.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			exists2, err2 := fs.client.BucketExists(ctx, bucket)
			if err2 != nil || !exists2 {
				return fmt.Errorf("create bucket: %w", err)
			}
		}
	}

	return nil
}
