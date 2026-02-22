package masterdata

import (
	"context"
)

func (uc *usecase) ValidateItemCodes(ctx context.Context, req *ValidateCodesRequest) (*ValidateCodesResponse, error) {
	results := make([]ValidationResult, len(req.Validations))
	allValid := true

	for i, v := range req.Validations {
		valid, err := uc.itemRepo.ValidateCode(ctx, v.CategoryCode, v.ItemCode, v.TenantID)
		if err != nil {
			return nil, err
		}

		result := ValidationResult{
			CategoryCode: v.CategoryCode,
			ItemCode:     v.ItemCode,
			Valid:        valid,
		}

		if !valid {
			allValid = false
			result.Message = "Item code not found"
		}

		results[i] = result
	}

	return &ValidateCodesResponse{
		AllValid: allValid,
		Results:  results,
	}, nil
}
