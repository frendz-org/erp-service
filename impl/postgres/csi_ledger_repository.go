package postgres

import (
	"context"

	"erp-service/entity"
	"erp-service/saving/participant"

	"gorm.io/gorm"
)

type csiLedgerRepository struct {
	baseRepository
}

func NewCsiLedgerRepository(db *gorm.DB) *csiLedgerRepository {
	return &csiLedgerRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *csiLedgerRepository) GetLatestPeriod(ctx context.Context, csiEmployeeID int) (int, int, error) {
	var result struct {
		YearPeriod  int
		MonthPeriod int
	}

	err := r.getDB(ctx).
		Model(&entity.CsiLedger{}).
		Where("csi_employee_id = ?", csiEmployeeID).
		Select("year_period, month_period").
		Order("year_period DESC, month_period DESC").
		Limit(1).
		Scan(&result).Error
	if err != nil {
		return 0, 0, translateError(err, "csi ledger")
	}

	return result.YearPeriod, result.MonthPeriod, nil
}

func (r *csiLedgerRepository) GetAmountsByItemID(ctx context.Context, csiEmployeeID int, upToYear int, upToMonth int) ([]participant.CsiLedgerItemAmount, error) {
	var results []participant.CsiLedgerItemAmount

	err := r.getDB(ctx).
		Model(&entity.CsiLedger{}).
		Where("csi_employee_id = ?", csiEmployeeID).
		Where("year_period < ? OR (year_period = ? AND month_period <= ?)", upToYear, upToYear, upToMonth).
		Select("csi_item_id, SUM(amount_trans) as total_amount").
		Group("csi_item_id").
		Order("csi_item_id").
		Scan(&results).Error
	if err != nil {
		return nil, translateError(err, "csi ledger")
	}

	return results, nil
}

func (r *csiLedgerRepository) applyPeriodFilter(db *gorm.DB, filter *participant.CsiLedgerPeriodFilter) *gorm.DB {
	db = db.Where("csi_employee_id = ?", filter.CsiEmployeeID)
	if filter.YearFrom != nil {
		db = db.Where("year_period >= ?", *filter.YearFrom)
	}
	if filter.YearTo != nil {
		db = db.Where("year_period <= ?", *filter.YearTo)
	}
	return db
}

func (r *csiLedgerRepository) CountDistinctPeriods(ctx context.Context, filter *participant.CsiLedgerPeriodFilter) (int64, error) {
	var count int64
	q := r.getDB(ctx).Model(&entity.CsiLedger{})
	q = r.applyPeriodFilter(q, filter)
	err := q.Select("COUNT(DISTINCT (year_period, month_period))").Scan(&count).Error
	if err != nil {
		return 0, translateError(err, "csi ledger")
	}
	return count, nil
}

func (r *csiLedgerRepository) GetDistinctPeriods(ctx context.Context, filter *participant.CsiLedgerPeriodFilter, offset int, limit int) ([]participant.CsiLedgerPeriod, error) {
	var results []participant.CsiLedgerPeriod
	q := r.getDB(ctx).Model(&entity.CsiLedger{})
	q = r.applyPeriodFilter(q, filter)
	err := q.Select("DISTINCT year_period, month_period").
		Order("year_period DESC, month_period DESC").
		Offset(offset).
		Limit(limit).
		Scan(&results).Error
	if err != nil {
		return nil, translateError(err, "csi ledger")
	}
	return results, nil
}

func (r *csiLedgerRepository) GetLedgersByPeriods(ctx context.Context, csiEmployeeID int, periods []participant.CsiLedgerPeriod) ([]participant.CsiLedgerRow, error) {
	if len(periods) == 0 {
		return nil, nil
	}

	periodPairs := make([][]int, 0, len(periods))
	for _, p := range periods {
		periodPairs = append(periodPairs, []int{p.YearPeriod, p.MonthPeriod})
	}

	var results []participant.CsiLedgerRow
	err := r.getDB(ctx).Model(&entity.CsiLedger{}).
		Where("csi_employee_id = ?", csiEmployeeID).
		Where("(year_period, month_period) IN ?", periodPairs).
		Select("year_period, month_period, csi_item_id, amount_trans").
		Order("year_period DESC, month_period DESC, csi_item_id ASC").
		Scan(&results).Error
	if err != nil {
		return nil, translateError(err, "csi ledger")
	}
	return results, nil
}

func (r *csiLedgerRepository) GetMonthlyBalances(ctx context.Context, csiEmployeeID int, yearFrom *int, yearTo *int) ([]participant.CsiLedgerMonthlyBalance, error) {
	var results []participant.CsiLedgerMonthlyBalance

	// CTE approach: compute full cumulative from all-time data, then filter output by year range.
	// This ensures the cumulative balance is correct even when yearFrom filters the returned rows.
	query := `
		WITH monthly_totals AS (
			SELECT year_period, month_period,
				SUM(amount_trans) AS period_amount
			FROM csi_ledgers
			WHERE csi_employee_id = ?
			GROUP BY year_period, month_period
		),
		cumulative AS (
			SELECT year_period, month_period,
				SUM(period_amount) OVER (ORDER BY year_period, month_period) AS balance
			FROM monthly_totals
		)
		SELECT year_period, month_period, balance
		FROM cumulative
		WHERE 1=1`

	args := []interface{}{csiEmployeeID}

	if yearFrom != nil {
		query += " AND year_period >= ?"
		args = append(args, *yearFrom)
	}
	if yearTo != nil {
		query += " AND year_period <= ?"
		args = append(args, *yearTo)
	}

	query += " ORDER BY year_period, month_period"

	err := r.getDB(ctx).Raw(query, args...).Scan(&results).Error
	if err != nil {
		return nil, translateError(err, "csi ledger")
	}

	return results, nil
}

func (r *csiLedgerRepository) GetCumulativeBalance(ctx context.Context, csiEmployeeID int, upToYear int, upToMonth int) (float64, error) {
	var balance float64
	err := r.getDB(ctx).Model(&entity.CsiLedger{}).
		Where("csi_employee_id = ?", csiEmployeeID).
		Where("year_period < ? OR (year_period = ? AND month_period <= ?)", upToYear, upToYear, upToMonth).
		Select("COALESCE(SUM(amount_trans), 0)").
		Scan(&balance).Error
	if err != nil {
		return 0, translateError(err, "csi ledger")
	}
	return balance, nil
}
