package internal

import (
	"context"

	"erp-service/saving/participant/participantdto"
)

func (uc *usecase) GetStatusHistory(ctx context.Context, req *participantdto.GetParticipantRequest) ([]participantdto.StatusHistoryResponse, error) {
	participant, err := uc.participantRepo.GetByID(ctx, req.ParticipantID)
	if err != nil {
		return nil, err
	}

	if err := validateParticipantOwnership(participant, req.TenantID, req.ProductID); err != nil {
		return nil, err
	}

	histories, err := uc.statusHistoryRepo.ListByParticipantID(ctx, req.ParticipantID)
	if err != nil {
		return nil, err
	}

	results := make([]participantdto.StatusHistoryResponse, 0, len(histories))
	for _, h := range histories {
		results = append(results, mapStatusHistoryToResponse(h))
	}

	return results, nil
}
