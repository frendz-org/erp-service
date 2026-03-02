package participant

import (
	"context"
)

func (uc *usecase) GetMyStatusHistory(ctx context.Context, req *GetMyParticipantRequest) ([]StatusHistoryResponse, error) {
	p, err := uc.participantRepo.GetByUserAndTenantProduct(ctx, req.UserID, req.TenantID, req.ProductID)
	if err != nil {
		return nil, err
	}

	histories, err := uc.statusHistoryRepo.ListByParticipantID(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	results := make([]StatusHistoryResponse, 0, len(histories))
	for _, h := range histories {
		results = append(results, mapStatusHistoryToResponse(h))
	}

	return results, nil
}
