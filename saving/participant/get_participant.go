package participant

import (
	"context"

)

func (uc *usecase) GetParticipant(ctx context.Context, req *GetParticipantRequest) (*ParticipantResponse, error) {
	participant, err := uc.participantRepo.GetByID(ctx, req.ParticipantID)
	if err != nil {
		return nil, err
	}

	if err := ValidateParticipantOwnership(participant, req.TenantID, req.ProductID); err != nil {
		return nil, err
	}

	return uc.buildFullParticipantResponse(ctx, participant, true)
}
