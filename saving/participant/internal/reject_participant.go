package internal

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant/participantdto"
)

func (uc *usecase) RejectParticipant(ctx context.Context, req *participantdto.RejectParticipantRequest) (*participantdto.ParticipantResponse, error) {
	var result *participantdto.ParticipantResponse

	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		participant, err := uc.participantRepo.GetByID(txCtx, req.ParticipantID)
		if err != nil {
			return fmt.Errorf("get participant: %w", err)
		}

		if err := validateParticipantOwnership(participant, req.TenantID, req.ProductID); err != nil {
			return err
		}

		if !participant.CanBeRejected() {
			return errors.ErrBadRequest(fmt.Sprintf("participant in %s status cannot be rejected", participant.Status))
		}

		now := time.Now()
		fromStatus := string(participant.Status)

		participant.Status = entity.ParticipantStatusRejected
		participant.RejectedBy = &req.UserID
		participant.RejectedAt = &now
		participant.RejectionReason = &req.Reason

		if err := uc.participantRepo.Update(txCtx, participant); err != nil {
			return fmt.Errorf("update participant: %w", err)
		}

		history := &entity.ParticipantStatusHistory{
			ParticipantID: participant.ID,
			FromStatus:    &fromStatus,
			ToStatus:      string(entity.ParticipantStatusRejected),
			ChangedBy:     req.UserID,
			Reason:        &req.Reason,
			ChangedAt:     now,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := uc.statusHistoryRepo.Create(txCtx, history); err != nil {
			return fmt.Errorf("create status history: %w", err)
		}

		resp, err := uc.buildFullParticipantResponse(txCtx, participant, false)
		if err != nil {
			return fmt.Errorf("build response: %w", err)
		}
		result = resp
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
