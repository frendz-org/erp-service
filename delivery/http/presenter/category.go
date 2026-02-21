package presenter

import (
	"erp-service/delivery/http/dto/response"
	"erp-service/masterdata/masterdatadto"
)

func ToCategoryResponse(category *masterdatadto.CategoryResponse) *response.MasterdataCategoryResponse {
	if category == nil {
		return nil
	}
	return &response.MasterdataCategoryResponse{
		ID:          category.ID,
		Code:        category.Code,
		Name:        category.Name,
		Description: category.Description,
		Status:      category.Status,
		Metadata:    category.Metadata,
	}
}

func ToCategoryListResponse(categories []*masterdatadto.CategoryResponse) []*response.MasterdataCategoryResponse {
	if categories == nil {
		return nil
	}
	result := make([]*response.MasterdataCategoryResponse, len(categories))
	for i, category := range categories {
		result[i] = ToCategoryResponse(category)
	}
	return result
}
