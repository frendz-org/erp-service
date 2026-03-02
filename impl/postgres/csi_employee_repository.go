package postgres

import (
	"context"

	"erp-service/entity"

	"gorm.io/gorm"
)

type csiEmployeeRepository struct {
	baseRepository
}

func NewCsiEmployeeRepository(db *gorm.DB) *csiEmployeeRepository {
	return &csiEmployeeRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *csiEmployeeRepository) GetByEmployeeNo(ctx context.Context, employeeNo string) (*entity.CsiEmployee, error) {
	var emp entity.CsiEmployee
	err := r.getDB(ctx).Where("employee_no = ?", employeeNo).First(&emp).Error
	if err != nil {
		return nil, translateError(err, "csi employee")
	}
	return &emp, nil
}
