package internal

import (
	"context"

	"iam-service/saving/participant/participantdto"
)

func (uc *usecase) GetParticipant(ctx context.Context, req *participantdto.GetParticipantRequest) (*participantdto.ParticipantResponse, error) {
	participant, err := uc.participantRepo.GetByID(ctx, req.ParticipantID)
	if err != nil {
		return nil, err
	}

	if err := validateParticipantOwnership(participant, req.TenantID, req.ProductID); err != nil {
		return nil, err
	}

	return uc.buildFullParticipantResponse(ctx, participant, true)
}
