package internal

import (
	"context"

	"erp-service/masterdata/masterdatadto"
)

func (uc *usecase) ValidateItemCodes(ctx context.Context, req *masterdatadto.ValidateCodesRequest) (*masterdatadto.ValidateCodesResponse, error) {
	results := make([]masterdatadto.ValidationResult, len(req.Validations))
	allValid := true

	for i, v := range req.Validations {
		valid, err := uc.itemRepo.ValidateCode(ctx, v.CategoryCode, v.ItemCode, v.TenantID)
		if err != nil {
			return nil, err
		}

		result := masterdatadto.ValidationResult{
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

	return &masterdatadto.ValidateCodesResponse{
		AllValid: allValid,
		Results:  results,
	}, nil
}
