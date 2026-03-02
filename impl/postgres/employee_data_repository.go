package postgres

import (
	"context"

	"erp-service/entity"
	"erp-service/saving/participant"

	"gorm.io/gorm"
)

type employeeDataRepository struct {
	baseRepository
}

func NewEmployeeDataRepository(db *gorm.DB) participant.EmployeeDataRepository {
	return &employeeDataRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *employeeDataRepository) GetByEmpNo(ctx context.Context, empNo string) (*entity.EmployeeData, error) {
	var emp entity.EmployeeData
	err := r.getDB(ctx).Where("emp_no = ?", empNo).First(&emp).Error
	if err != nil {
		return nil, translateError(err, "employee data")
	}
	return &emp, nil
}
