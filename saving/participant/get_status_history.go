package participant

import (
	"context"

)

func (uc *usecase) GetStatusHistory(ctx context.Context, req *GetParticipantRequest) ([]StatusHistoryResponse, error) {
	participant, err := uc.participantRepo.GetByID(ctx, req.ParticipantID)
	if err != nil {
		return nil, err
	}

	if err := ValidateParticipantOwnership(participant, req.TenantID, req.ProductID); err != nil {
		return nil, err
	}

	histories, err := uc.statusHistoryRepo.ListByParticipantID(ctx, req.ParticipantID)
	if err != nil {
		return nil, err
	}

	results := make([]StatusHistoryResponse, 0, len(histories))
	for _, h := range histories {
		results = append(results, mapStatusHistoryToResponse(h))
	}

	return results, nil
}
