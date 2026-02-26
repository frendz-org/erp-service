package entity

import "time"

type CsiEmployee struct {
	CsiEmployeeID     int        `gorm:"column:csi_employee_id;primaryKey" json:"csi_employee_id"`
	EmployeeNo        string     `gorm:"column:employee_no" json:"employee_no"`
	EmployeeName      string     `gorm:"column:employee_name" json:"employee_name"`
	Gender            string     `gorm:"column:gender" json:"gender"`
	BirthPlace        string     `gorm:"column:birth_place" json:"birth_place"`
	BirthDate         *time.Time `gorm:"column:birth_date" json:"birth_date"`
	RetirementDate    *time.Time `gorm:"column:retirement_date" json:"retirement_date"`
	MaritalStatus     string     `gorm:"column:marital_status" json:"marital_status"`
	MobilePhoneNo     string     `gorm:"column:mobile_phone_no" json:"mobile_phone_no"`
	OpuNo             string     `gorm:"column:opu_no" json:"opu_no"`
	GroupName         string     `gorm:"column:group_name" json:"group_name"`
	StatusName        string     `gorm:"column:status_name" json:"status_name"`
	CostCenterNo      string     `gorm:"column:cost_center_no" json:"cost_center_no"`
	JoinDate          *time.Time `gorm:"column:join_date" json:"join_date"`
	StartDate         *time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate           *time.Time `gorm:"column:end_date" json:"end_date"`
	PensionCategoryNo string     `gorm:"column:pension_category_no" json:"pension_category_no"`
	AmountBalance     *float64   `gorm:"column:amount_balance;type:numeric(18,4)" json:"amount_balance"`
	AccountNo         string     `gorm:"column:account_no" json:"account_no"`
	Photo             string     `gorm:"column:photo" json:"photo"`
	IsActive          bool       `gorm:"column:is_active" json:"is_active"`
	LastUpdated       *time.Time `gorm:"column:last_updated" json:"last_updated"`
	LastUpdater       *int       `gorm:"column:last_updater" json:"last_updater"`
}

func (CsiEmployee) TableName() string {
	return "csi_employees"
}
