package presenter

import (
	"erp-service/delivery/http/dto/response"
	"erp-service/masterdata/masterdatadto"
)

func ToItemResponse(item *masterdatadto.ItemResponse) *response.MasterdataItemResponse {
	if item == nil {
		return nil
	}
	return &response.MasterdataItemResponse{
		ID:           item.ID,
		Code:         item.Code,
		Name:         item.Name,
		AltName:      item.AltName,
		Description:  item.Description,
		ParentItemID: item.ParentItemID,
		Status:       item.Status,
		IsDefault:    item.IsDefault,
		Metadata:     item.Metadata,
	}
}

func ToItemListResponse(items []*masterdatadto.ItemResponse) []*response.MasterdataItemResponse {
	if items == nil {
		return nil
	}
	result := make([]*response.MasterdataItemResponse, len(items))
	for i, item := range items {
		result[i] = ToItemResponse(item)
	}
	return result
}

func ToItemTreeResponse(items []*masterdatadto.ItemResponse) []*response.MasterdataItemTreeResponse {
	if items == nil {
		return nil
	}

	itemMap := make(map[string]*response.MasterdataItemTreeResponse)
	var roots []*response.MasterdataItemTreeResponse

	for _, item := range items {
		node := &response.MasterdataItemTreeResponse{
			ID:           item.ID,
			Code:         item.Code,
			Name:         item.Name,
			AltName:      item.AltName,
			ParentItemID: item.ParentItemID,
			Status:       item.Status,
			Metadata:     item.Metadata,
			Children:     []*response.MasterdataItemTreeResponse{},
		}
		itemMap[item.ID.String()] = node
	}

	for _, item := range items {
		node := itemMap[item.ID.String()]
		if item.ParentItemID == nil {
			roots = append(roots, node)
		} else {
			if parent, ok := itemMap[item.ParentItemID.String()]; ok {
				parent.Children = append(parent.Children, node)
			} else {

				roots = append(roots, node)
			}
		}
	}

	return roots
}

func ToValidateCodeResponse(resp *masterdatadto.ValidateCodeResponse) *response.ValidateCodeResponse {
	if resp == nil {
		return nil
	}
	return &response.ValidateCodeResponse{
		Valid:   resp.Valid,
		Message: resp.Message,
	}
}

func ToValidateCodesResponse(resp *masterdatadto.ValidateCodesResponse) *response.ValidateCodesResponse {
	if resp == nil {
		return nil
	}

	results := make([]response.ValidationResult, len(resp.Results))
	for i, r := range resp.Results {
		results[i] = response.ValidationResult{
			CategoryCode: r.CategoryCode,
			ItemCode:     r.ItemCode,
			Valid:        r.Valid,
			Message:      r.Message,
		}
	}

	return &response.ValidateCodesResponse{
		AllValid: resp.AllValid,
		Results:  results,
	}
}
