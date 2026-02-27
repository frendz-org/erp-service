package participant

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildQuarterlyDataPoints_BasicGrouping(t *testing.T) {
	balances := []CsiLedgerMonthlyBalance{
		{YearPeriod: 2024, MonthPeriod: 1, Balance: 100.0},
		{YearPeriod: 2024, MonthPeriod: 2, Balance: 200.0},
		{YearPeriod: 2024, MonthPeriod: 3, Balance: 300.0},
		{YearPeriod: 2024, MonthPeriod: 4, Balance: 400.0},
		{YearPeriod: 2024, MonthPeriod: 5, Balance: 500.0},
		{YearPeriod: 2024, MonthPeriod: 6, Balance: 600.0},
	}

	points := buildQuarterlyDataPoints(balances)

	require.Len(t, points, 2)

	assert.Equal(t, "2024-Q1", points[0].PeriodLabel)
	assert.Equal(t, 2024, points[0].Year)
	require.NotNil(t, points[0].Quarter)
	assert.Equal(t, 1, *points[0].Quarter)
	assert.Nil(t, points[0].Month)
	assert.Equal(t, int64(300), points[0].Balance)

	assert.Equal(t, "2024-Q2", points[1].PeriodLabel)
	assert.Equal(t, 2024, points[1].Year)
	require.NotNil(t, points[1].Quarter)
	assert.Equal(t, 2, *points[1].Quarter)
	assert.Nil(t, points[1].Month)
	assert.Equal(t, int64(600), points[1].Balance)
}

func TestBuildQuarterlyDataPoints_FullYear(t *testing.T) {
	var balances []CsiLedgerMonthlyBalance
	for month := 1; month <= 12; month++ {
		balances = append(balances, CsiLedgerMonthlyBalance{
			YearPeriod:  2023,
			MonthPeriod: month,
			Balance:     float64(month * 1000),
		})
	}

	points := buildQuarterlyDataPoints(balances)

	require.Len(t, points, 4)

	assert.Equal(t, "2023-Q1", points[0].PeriodLabel)
	assert.Equal(t, int64(3000), points[0].Balance)
	assert.Equal(t, 1, *points[0].Quarter)

	assert.Equal(t, "2023-Q2", points[1].PeriodLabel)
	assert.Equal(t, int64(6000), points[1].Balance)
	assert.Equal(t, 2, *points[1].Quarter)

	assert.Equal(t, "2023-Q3", points[2].PeriodLabel)
	assert.Equal(t, int64(9000), points[2].Balance)
	assert.Equal(t, 3, *points[2].Quarter)

	assert.Equal(t, "2023-Q4", points[3].PeriodLabel)
	assert.Equal(t, int64(12000), points[3].Balance)
	assert.Equal(t, 4, *points[3].Quarter)
}

func TestBuildQuarterlyDataPoints_IncompleteQuarter(t *testing.T) {

	balances := []CsiLedgerMonthlyBalance{
		{YearPeriod: 2024, MonthPeriod: 7, Balance: 700.0},
		{YearPeriod: 2024, MonthPeriod: 8, Balance: 800.0},
	}

	points := buildQuarterlyDataPoints(balances)

	require.Len(t, points, 1)
	assert.Equal(t, "2024-Q3", points[0].PeriodLabel)
	assert.Equal(t, 3, *points[0].Quarter)

	assert.Equal(t, int64(800), points[0].Balance)
}

func TestBuildQuarterlyDataPoints_MultiYear(t *testing.T) {
	balances := []CsiLedgerMonthlyBalance{
		{YearPeriod: 2022, MonthPeriod: 12, Balance: 1000.0},
		{YearPeriod: 2023, MonthPeriod: 3, Balance: 2000.0},
		{YearPeriod: 2023, MonthPeriod: 6, Balance: 3000.0},
		{YearPeriod: 2024, MonthPeriod: 9, Balance: 4000.0},
	}

	points := buildQuarterlyDataPoints(balances)

	require.Len(t, points, 4)

	assert.Equal(t, "2022-Q4", points[0].PeriodLabel)
	assert.Equal(t, 2022, points[0].Year)
	assert.Equal(t, 4, *points[0].Quarter)
	assert.Equal(t, int64(1000), points[0].Balance)

	assert.Equal(t, "2023-Q1", points[1].PeriodLabel)
	assert.Equal(t, 2023, points[1].Year)
	assert.Equal(t, 1, *points[1].Quarter)
	assert.Equal(t, int64(2000), points[1].Balance)

	assert.Equal(t, "2023-Q2", points[2].PeriodLabel)
	assert.Equal(t, 2023, points[2].Year)
	assert.Equal(t, 2, *points[2].Quarter)
	assert.Equal(t, int64(3000), points[2].Balance)

	assert.Equal(t, "2024-Q3", points[3].PeriodLabel)
	assert.Equal(t, 2024, points[3].Year)
	assert.Equal(t, 3, *points[3].Quarter)
	assert.Equal(t, int64(4000), points[3].Balance)
}

func TestBuildQuarterlyDataPoints_EmptyInput(t *testing.T) {
	points := buildQuarterlyDataPoints(nil)
	assert.Empty(t, points)

	points = buildQuarterlyDataPoints([]CsiLedgerMonthlyBalance{})
	assert.Empty(t, points)
}

func TestBuildQuarterlyDataPoints_QuarterBoundaryAllMonths(t *testing.T) {

	expectedQuarters := map[int]int{
		1: 1, 2: 1, 3: 1,
		4: 2, 5: 2, 6: 2,
		7: 3, 8: 3, 9: 3,
		10: 4, 11: 4, 12: 4,
	}

	for month, expectedQ := range expectedQuarters {
		balances := []CsiLedgerMonthlyBalance{
			{YearPeriod: 2024, MonthPeriod: month, Balance: float64(month * 100)},
		}
		points := buildQuarterlyDataPoints(balances)
		require.Len(t, points, 1, "month %d should produce exactly 1 data point", month)
		assert.Equal(t, expectedQ, *points[0].Quarter, "month %d should map to Q%d", month, expectedQ)
	}
}

func TestBuildQuarterlyDataPoints_BalanceRounding(t *testing.T) {
	balances := []CsiLedgerMonthlyBalance{
		{YearPeriod: 2024, MonthPeriod: 3, Balance: 1234567.678},
	}

	points := buildQuarterlyDataPoints(balances)

	require.Len(t, points, 1)

	assert.Equal(t, int64(1234568), points[0].Balance)
}

func TestBuildQuarterlyDataPoints_SingleMonthPerQuarter(t *testing.T) {

	balances := []CsiLedgerMonthlyBalance{
		{YearPeriod: 2024, MonthPeriod: 2, Balance: 999.5},
	}

	points := buildQuarterlyDataPoints(balances)

	require.Len(t, points, 1)
	assert.Equal(t, "2024-Q1", points[0].PeriodLabel)
	assert.Equal(t, int64(1000), points[0].Balance)
}

func TestBuildMonthlyDataPoints(t *testing.T) {

	balances := []CsiLedgerMonthlyBalance{
		{YearPeriod: 2024, MonthPeriod: 1, Balance: 100.7},
		{YearPeriod: 2024, MonthPeriod: 2, Balance: 200.3},
	}

	points := buildMonthlyDataPoints(balances)

	require.Len(t, points, 2)
	assert.Equal(t, "2024-01", points[0].PeriodLabel)
	assert.Equal(t, 2024, points[0].Year)
	require.NotNil(t, points[0].Month)
	assert.Equal(t, 1, *points[0].Month)
	assert.Nil(t, points[0].Quarter)
	assert.Equal(t, int64(101), points[0].Balance)

	assert.Equal(t, "2024-02", points[1].PeriodLabel)
	assert.Equal(t, int64(200), points[1].Balance)
}

func TestBuildYearlyDataPoints(t *testing.T) {

	balances := []CsiLedgerMonthlyBalance{
		{YearPeriod: 2023, MonthPeriod: 6, Balance: 500.0},
		{YearPeriod: 2023, MonthPeriod: 12, Balance: 1000.0},
		{YearPeriod: 2024, MonthPeriod: 3, Balance: 1500.0},
	}

	points := buildYearlyDataPoints(balances)

	require.Len(t, points, 2)
	assert.Equal(t, "2023", points[0].PeriodLabel)
	assert.Equal(t, 2023, points[0].Year)
	assert.Nil(t, points[0].Month)
	assert.Nil(t, points[0].Quarter)
	assert.Equal(t, int64(1000), points[0].Balance)

	assert.Equal(t, "2024", points[1].PeriodLabel)
	assert.Equal(t, int64(1500), points[1].Balance)
}
