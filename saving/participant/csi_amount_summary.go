package participant

import (
	"context"
	"fmt"
	"math"
	"strings"

	apperrors "erp-service/pkg/errors"

	"github.com/google/uuid"
)

var csiItemNames = map[int]string{
	1: "Kontribusi Peserta",
	2: "Kontribusi Perusahaan",
	3: "Hasil Pengembangan",
}

func (uc *usecase) GetCsiAmountSummary(ctx context.Context, req *CsiAmountSummaryRequest) ([]CsiAmountSummaryResponse, error) {
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
		return []CsiAmountSummaryResponse{buildZeroSummary(fullName, "", "")}, nil
	}

	var summaries []CsiAmountSummaryResponse
	for _, p := range participants {
		tenantName := ""
		tenant, tErr := uc.tenantRepo.GetByID(ctx, p.TenantID)
		if tErr == nil {
			tenantName = tenant.Name
		}

		if p.EmployeeNumber == nil || *p.EmployeeNumber == "" {
			summaries = append(summaries, buildZeroSummary(fullName, p.TenantID.String(), tenantName))
			continue
		}

		summary, sErr := uc.buildCsiSummary(ctx, fullName, p.TenantID.String(), tenantName, *p.EmployeeNumber)
		if sErr != nil {
			if apperrors.IsNotFound(sErr) {
				summaries = append(summaries, buildZeroSummary(fullName, p.TenantID.String(), tenantName))
				continue
			}
			return nil, sErr
		}
		summaries = append(summaries, *summary)
	}

	return summaries, nil
}

func (uc *usecase) buildCsiSummary(ctx context.Context, fullName, tenantID, tenantName, empNo string) (*CsiAmountSummaryResponse, error) {
	csiEmp, err := uc.csiEmployeeRepo.GetByEmployeeNo(ctx, empNo)
	if err != nil {
		return nil, err
	}

	latestYear, latestMonth, err := uc.csiLedgerRepo.GetLatestPeriod(ctx, csiEmp.CsiEmployeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest period: %w", err)
	}
	if latestYear == 0 && latestMonth == 0 {
		return ptrSummary(buildZeroSummary(fullName, tenantID, tenantName)), nil
	}

	currentAmounts, err := uc.csiLedgerRepo.GetAmountsByItemID(ctx, csiEmp.CsiEmployeeID, latestYear, latestMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get current amounts: %w", err)
	}

	prevYear, prevMonth := previousPeriod(latestYear, latestMonth)
	prevAmounts, err := uc.csiLedgerRepo.GetAmountsByItemID(ctx, csiEmp.CsiEmployeeID, prevYear, prevMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous amounts: %w", err)
	}

	var totalSaldo float64
	for _, a := range currentAmounts {
		totalSaldo += a.TotalAmount
	}

	var prevTotal float64
	for _, a := range prevAmounts {
		prevTotal += a.TotalAmount
	}

	growthPercentage := calculateGrowth(totalSaldo, prevTotal)

	items := buildItems(currentAmounts, totalSaldo)
	latestPeriod := fmt.Sprintf("%d-%02d", latestYear, latestMonth)

	return &CsiAmountSummaryResponse{
		TenantID:         parseUUID(tenantID),
		TenantName:       tenantName,
		User:             fullName,
		TotalSaldo:       int64(math.Round(totalSaldo)),
		GrowthPercentage: growthPercentage,
		Items:            items,
		LatestPeriod:     latestPeriod,
	}, nil
}

func buildZeroSummary(fullName, tenantID, tenantName string) CsiAmountSummaryResponse {
	return CsiAmountSummaryResponse{
		TenantID:         parseUUID(tenantID),
		TenantName:       tenantName,
		User:             fullName,
		TotalSaldo:       0,
		GrowthPercentage: 0,
		Items:            []CsiAmountSummaryItem{},
		LatestPeriod:     "",
	}
}

func buildItems(amounts []CsiLedgerItemAmount, totalSaldo float64) []CsiAmountSummaryItem {
	items := make([]CsiAmountSummaryItem, 0, len(amounts))
	for _, a := range amounts {
		name, ok := csiItemNames[a.CsiItemID]
		if !ok {
			continue
		}

		var pct float64
		if totalSaldo > 0 {
			pct = math.Round((a.TotalAmount/totalSaldo)*1000) / 10
		}

		items = append(items, CsiAmountSummaryItem{
			CsiItemID:  a.CsiItemID,
			Name:       name,
			Amount:     int64(math.Round(a.TotalAmount)),
			Percentage: pct,
		})
	}
	return items
}

func calculateGrowth(current, previous float64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100.0
		}
		return 0
	}
	growth := ((current - previous) / previous) * 100
	return math.Round(growth*100) / 100
}

func previousPeriod(year, month int) (int, int) {
	if month == 1 {
		return year - 1, 12
	}
	return year, month - 1
}

func ptrSummary(s CsiAmountSummaryResponse) *CsiAmountSummaryResponse {
	return &s
}

func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
