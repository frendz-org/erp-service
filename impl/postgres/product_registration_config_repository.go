package postgres

import (
	"context"

	"erp-service/entity"
	membercontract "erp-service/saving/member/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRegistrationConfigRepository struct {
	baseRepository
}

func NewProductRegistrationConfigRepository(db *gorm.DB) membercontract.ProductRegistrationConfigRepository {
	return &productRegistrationConfigRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *productRegistrationConfigRepository) GetByProductAndType(ctx context.Context, productID uuid.UUID, regType string) (*entity.ProductRegistrationConfig, error) {
	var config entity.ProductRegistrationConfig
	err := r.getDB(ctx).
		Where("product_id = ? AND registration_type = ?", productID, regType).
		First(&config).Error
	if err != nil {
		return nil, translateError(err, "product registration config")
	}
	return &config, nil
}
