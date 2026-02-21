package postgres

import (
	"context"

	"erp-service/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	baseRepository
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{
		baseRepository: baseRepository{db: db},
	}
}

func NewProductsByTenantRepository(db *gorm.DB) *productRepository {
	return &productRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *productRepository) GetByCodeAndTenant(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Product, error) {
	var product entity.Product
	err := r.getDB(ctx).Where("tenant_id = ? AND code = ? AND status = 'ACTIVE' AND deleted_at IS NULL",
		tenantID, code).First(&product).Error
	if err != nil {
		return nil, translateError(err, "product")
	}
	return &product, nil
}

func (r *productRepository) GetByIDAndTenant(ctx context.Context, productID, tenantID uuid.UUID) (*entity.Product, error) {
	var product entity.Product
	err := r.getDB(ctx).Where("id = ? AND tenant_id = ? AND status = 'ACTIVE' AND deleted_at IS NULL",
		productID, tenantID).First(&product).Error
	if err != nil {
		return nil, translateError(err, "product")
	}
	return &product, nil
}

func (r *productRepository) ListActiveByTenantID(ctx context.Context, tenantID uuid.UUID) ([]entity.Product, error) {
	var products []entity.Product
	err := r.getDB(ctx).
		Where("tenant_id = ? AND status = 'ACTIVE' AND deleted_at IS NULL", tenantID).
		Find(&products).Error
	if err != nil {
		return nil, translateError(err, "product")
	}
	return products, nil
}
