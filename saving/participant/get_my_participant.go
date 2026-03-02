package participant

import (
	"context"

	"erp-service/pkg/errors"

	"go.uber.org/zap"
)

func (uc *usecase) GetMyParticipant(ctx context.Context, req *GetMyParticipantRequest) (*MyParticipantResponse, error) {
	p, err := uc.participantRepo.GetByUserAndTenantProduct(ctx, req.UserID, req.TenantID, req.ProductID)
	if err != nil {
		return nil, err
	}

	participantResp, err := uc.buildFullParticipantResponse(ctx, p, true)
	if err != nil {
		return nil, err
	}

	registrationStatus := ""
	utr, err := uc.utrRepo.GetByUserAndProduct(ctx, req.UserID, req.TenantID, req.ProductID, "PARTICIPANT")
	if err != nil {
		if !errors.IsNotFound(err) {
			return nil, err
		}
		uc.logger.Warn("participant exists but UTR not found",
			zap.String("user_id", req.UserID.String()),
			zap.String("tenant_id", req.TenantID.String()),
		)
	} else {
		registrationStatus = string(utr.Status)
	}

	return &MyParticipantResponse{
		ParticipantResponse: *participantResp,
		RegistrationStatus:  registrationStatus,
	}, nil
}
