package entity

type EmployeeData struct {
	ID      int    `gorm:"column:id;primaryKey"`
	EmpNo   string `gorm:"column:emp_no"`
	EmpName string `gorm:"column:emp_name"`
}

func (EmployeeData) TableName() string {
	return "employee_data"
}
