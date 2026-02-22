package infrastructure

import (
	"fmt"

	"erp-service/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOClient(cfg *config.Config) (*minio.Client, error) {
	client, err := minio.New(cfg.Infra.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Infra.Minio.AccessKey, cfg.Infra.Minio.SecretKey, ""),
		Secure: cfg.Infra.Minio.UseSSL,
		Region: cfg.Infra.Minio.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return client, nil
}
