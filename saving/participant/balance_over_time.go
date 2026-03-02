package participant

import (
	"context"
	"fmt"
	"math"
	"strings"

	apperrors "erp-service/pkg/errors"
)

func (uc *usecase) GetBalanceOverTime(ctx context.Context, req *BalanceOverTimeRequest) ([]BalanceOverTimeResponse, error) {
	profile, err := uc.userProfileRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve user profile: %w", err)
	}
	fullName := strings.TrimSpace(profile.FirstName + " " + profile.LastName)

	participants, err := uc.participantRepo.ListByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to list participants: %w", err)
	}

	if len(participants) == 0 {
		return []BalanceOverTimeResponse{buildEmptyBalanceOverTime(fullName, "", "", req.Granularity)}, nil
	}

	responses := make([]BalanceOverTimeResponse, 0, len(participants))
	for _, p := range participants {
		tenantName := ""
		tenant, tErr := uc.tenantRepo.GetByID(ctx, p.TenantID)
		if tErr != nil {
			if !apperrors.IsNotFound(tErr) {
				return nil, fmt.Errorf("failed to resolve tenant: %w", tErr)
			}

		} else {
			tenantName = tenant.Name
		}

		if p.EmployeeNumber == nil || *p.EmployeeNumber == "" {
			responses = append(responses, buildEmptyBalanceOverTime(fullName, p.TenantID.String(), tenantName, req.Granularity))
			continue
		}

		bot, bErr := uc.buildBalanceOverTime(ctx, fullName, p.TenantID.String(), tenantName, *p.EmployeeNumber, req)
		if bErr != nil {
			if apperrors.IsNotFound(bErr) {
				responses = append(responses, buildEmptyBalanceOverTime(fullName, p.TenantID.String(), tenantName, req.Granularity))
				continue
			}
			return nil, bErr
		}
		responses = append(responses, *bot)
	}

	return responses, nil
}

func (uc *usecase) buildBalanceOverTime(ctx context.Context, fullName, tenantID, tenantName, empNo string, req *BalanceOverTimeRequest) (*BalanceOverTimeResponse, error) {
	csiEmp, err := uc.csiEmployeeRepo.GetByEmployeeNo(ctx, empNo)
	if err != nil {
		return nil, err
	}

	monthlyBalances, err := uc.csiLedgerRepo.GetMonthlyBalances(ctx, csiEmp.CsiEmployeeID, req.YearFrom, req.YearTo)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly balances: %w", err)
	}

	if len(monthlyBalances) == 0 {
		empty := buildEmptyBalanceOverTime(fullName, tenantID, tenantName, req.Granularity)
		return &empty, nil
	}

	var dataPoints []BalanceOverTimeDataPoint

	switch req.Granularity {
	case "monthly":
		dataPoints = buildMonthlyDataPoints(monthlyBalances)
	case "quarterly":
		dataPoints = buildQuarterlyDataPoints(monthlyBalances)
	case "yearly":
		dataPoints = buildYearlyDataPoints(monthlyBalances)
	default:
		return nil, apperrors.ErrBadRequest(fmt.Sprintf("unsupported granularity: %q", req.Granularity))
	}

	return &BalanceOverTimeResponse{
		TenantID:    parseUUID(tenantID),
		TenantName:  tenantName,
		User:        fullName,
		Granularity: req.Granularity,
		DataPoints:  dataPoints,
	}, nil
}

func buildMonthlyDataPoints(balances []CsiLedgerMonthlyBalance) []BalanceOverTimeDataPoint {
	points := make([]BalanceOverTimeDataPoint, 0, len(balances))
	for _, b := range balances {
		month := b.MonthPeriod
		points = append(points, BalanceOverTimeDataPoint{
			PeriodLabel: fmt.Sprintf("%d-%02d", b.YearPeriod, b.MonthPeriod),
			Year:        b.YearPeriod,
			Month:       &month,
			Balance:     int64(math.Round(b.Balance)),
		})
	}
	return points
}

func buildQuarterlyDataPoints(balances []CsiLedgerMonthlyBalance) []BalanceOverTimeDataPoint {
	type quarterKey struct {
		year    int
		quarter int
	}

	initialCap := len(balances)/3 + 1
	quarterBalance := make(map[quarterKey]int64, initialCap)
	keys := make([]quarterKey, 0, initialCap)

	for _, b := range balances {
		q := (b.MonthPeriod-1)/3 + 1
		k := quarterKey{year: b.YearPeriod, quarter: q}
		if _, exists := quarterBalance[k]; !exists {
			keys = append(keys, k)
		}

		quarterBalance[k] = int64(math.Round(b.Balance))
	}

	points := make([]BalanceOverTimeDataPoint, 0, len(keys))
	for _, k := range keys {
		q := k.quarter
		points = append(points, BalanceOverTimeDataPoint{
			PeriodLabel: fmt.Sprintf("%d-Q%d", k.year, k.quarter),
			Year:        k.year,
			Quarter:     &q,
			Balance:     quarterBalance[k],
		})
	}
	return points
}

func buildYearlyDataPoints(balances []CsiLedgerMonthlyBalance) []BalanceOverTimeDataPoint {
	initialCap := len(balances)/12 + 1
	yearBalance := make(map[int]int64, initialCap)
	years := make([]int, 0, initialCap)

	for _, b := range balances {
		if _, exists := yearBalance[b.YearPeriod]; !exists {
			years = append(years, b.YearPeriod)
		}

		yearBalance[b.YearPeriod] = int64(math.Round(b.Balance))
	}

	points := make([]BalanceOverTimeDataPoint, 0, len(years))
	for _, year := range years {
		points = append(points, BalanceOverTimeDataPoint{
			PeriodLabel: fmt.Sprintf("%d", year),
			Year:        year,
			Balance:     yearBalance[year],
		})
	}
	return points
}

func buildEmptyBalanceOverTime(fullName, tenantID, tenantName, granularity string) BalanceOverTimeResponse {
	return BalanceOverTimeResponse{
		TenantID:    parseUUID(tenantID),
		TenantName:  tenantName,
		User:        fullName,
		Granularity: granularity,
		DataPoints:  []BalanceOverTimeDataPoint{},
	}
}
