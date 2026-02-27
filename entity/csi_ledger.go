package entity

import "time"

type CsiLedger struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CsiEmployeeID int        `gorm:"column:csi_employee_id;not null" json:"csi_employee_id"`
	YearPeriod    int        `gorm:"column:year_period;not null" json:"year_period"`
	MonthPeriod   int        `gorm:"column:month_period;not null" json:"month_period"`
	CsiItemID     int        `gorm:"column:csi_item_id;not null" json:"csi_item_id"`
	AmountTrans   float64    `gorm:"column:amount_trans;type:numeric(18,4);not null" json:"amount_trans"`
	LastUpdated   *time.Time `gorm:"column:last_updated" json:"last_updated"`
	LastUpdater   *int       `gorm:"column:last_updater" json:"last_updater"`
}

func (CsiLedger) TableName() string {
	return "csi_ledgers"
}
