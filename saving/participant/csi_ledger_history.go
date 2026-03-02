package participant

import (
	"context"
	"fmt"
	"math"
	"strings"
)

func (uc *usecase) GetCsiLedgerHistory(ctx context.Context, req *CsiLedgerHistoryRequest) ([]CsiLedgerHistoryResponse, error) {
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
		return []CsiLedgerHistoryResponse{buildEmptyHistory(fullName, "", "", req)}, nil
	}

	var responses []CsiLedgerHistoryResponse
	for _, p := range participants {
		tenantName := ""
		tenant, tErr := uc.tenantRepo.GetByID(ctx, p.TenantID)
		if tErr == nil {
			tenantName = tenant.Name
		}

		if p.EmployeeNumber == nil || *p.EmployeeNumber == "" {
			responses = append(responses, buildEmptyHistory(fullName, p.TenantID.String(), tenantName, req))
			continue
		}

		history, hErr := uc.buildLedgerHistory(ctx, fullName, p.TenantID.String(), tenantName, *p.EmployeeNumber, req)
		if hErr != nil {
			responses = append(responses, buildEmptyHistory(fullName, p.TenantID.String(), tenantName, req))
			continue
		}
		responses = append(responses, *history)
	}

	return responses, nil
}

func (uc *usecase) buildLedgerHistory(ctx context.Context, fullName, tenantID, tenantName, empNo string, req *CsiLedgerHistoryRequest) (*CsiLedgerHistoryResponse, error) {
	csiEmp, err := uc.csiEmployeeRepo.GetByEmployeeNo(ctx, empNo)
	if err != nil {
		return nil, err
	}

	filter := &CsiLedgerPeriodFilter{
		CsiEmployeeID: csiEmp.CsiEmployeeID,
		YearFrom:      req.YearFrom,
		YearTo:        req.YearTo,
	}

	totalRecords, err := uc.csiLedgerRepo.CountDistinctPeriods(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count periods: %w", err)
	}

	if totalRecords == 0 {
		return ptrHistory(buildEmptyHistory(fullName, tenantID, tenantName, req)), nil
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(req.PerPage)))
	offset := (req.Page - 1) * req.PerPage

	periods, err := uc.csiLedgerRepo.GetDistinctPeriods(ctx, filter, offset, req.PerPage)
	if err != nil {
		return nil, fmt.Errorf("failed to get periods: %w", err)
	}

	if len(periods) == 0 {
		return ptrHistory(buildEmptyHistory(fullName, tenantID, tenantName, req)), nil
	}

	rows, err := uc.csiLedgerRepo.GetLedgersByPeriods(ctx, csiEmp.CsiEmployeeID, periods)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger rows: %w", err)
	}

	periodResponses, err := uc.buildPeriodResponses(ctx, csiEmp.CsiEmployeeID, periods, rows)
	if err != nil {
		return nil, err
	}

	return &CsiLedgerHistoryResponse{
		TenantID:   parseUUID(tenantID),
		TenantName: tenantName,
		User:       fullName,
		Periods:    periodResponses,
		Pagination: CsiLedgerPaginationResponse{
			Page:         req.Page,
			PerPage:      req.PerPage,
			TotalRecords: totalRecords,
			TotalPages:   totalPages,
		},
	}, nil
}

func (uc *usecase) buildPeriodResponses(ctx context.Context, csiEmployeeID int, periods []CsiLedgerPeriod, rows []CsiLedgerRow) ([]CsiLedgerPeriodResponse, error) {
	rowsByPeriod := make(map[string][]CsiLedgerRow)
	for _, r := range rows {
		key := fmt.Sprintf("%d-%02d", r.YearPeriod, r.MonthPeriod)
		rowsByPeriod[key] = append(rowsByPeriod[key], r)
	}

	responses := make([]CsiLedgerPeriodResponse, 0, len(periods))
	for _, p := range periods {
		key := fmt.Sprintf("%d-%02d", p.YearPeriod, p.MonthPeriod)

		cumBalance, err := uc.csiLedgerRepo.GetCumulativeBalance(ctx, csiEmployeeID, p.YearPeriod, p.MonthPeriod)
		if err != nil {
			return nil, fmt.Errorf("failed to get cumulative balance for %s: %w", key, err)
		}

		periodRows := rowsByPeriod[key]
		items, periodTotal := buildPeriodItems(periodRows)

		responses = append(responses, CsiLedgerPeriodResponse{
			YearPeriod:        p.YearPeriod,
			MonthPeriod:       p.MonthPeriod,
			PeriodLabel:       key,
			Items:             items,
			PeriodTotal:       periodTotal,
			CumulativeBalance: int64(math.Round(cumBalance)),
		})
	}

	return responses, nil
}

func buildPeriodItems(rows []CsiLedgerRow) ([]CsiLedgerPeriodItem, int64) {
	amountByItem := make(map[int]float64)
	for _, r := range rows {
		amountByItem[r.CsiItemID] = r.AmountTrans
	}

	var totalRaw float64
	items := make([]CsiLedgerPeriodItem, 0, 3)
	for _, itemID := range []int{1, 2, 3} {
		amt := amountByItem[itemID]
		totalRaw += amt
		name := csiItemNames[itemID]
		items = append(items, CsiLedgerPeriodItem{
			CsiItemID: itemID,
			Name:      name,
			Amount:    int64(math.Round(amt)),
		})
	}

	return items, int64(math.Round(totalRaw))
}

func buildEmptyHistory(fullName, tenantID, tenantName string, req *CsiLedgerHistoryRequest) CsiLedgerHistoryResponse {
	return CsiLedgerHistoryResponse{
		TenantID:   parseUUID(tenantID),
		TenantName: tenantName,
		User:       fullName,
		Periods:    []CsiLedgerPeriodResponse{},
		Pagination: CsiLedgerPaginationResponse{
			Page:         req.Page,
			PerPage:      req.PerPage,
			TotalRecords: 0,
			TotalPages:   0,
		},
	}
}

func ptrHistory(h CsiLedgerHistoryResponse) *CsiLedgerHistoryResponse {
	return &h
}
