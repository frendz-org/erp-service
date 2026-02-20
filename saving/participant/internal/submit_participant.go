package internal

import (
	"context"
	"fmt"
	"time"

	"iam-service/entity"
	"iam-service/saving/participant/participantdto"
	"iam-service/pkg/errors"
)

func (uc *usecase) SubmitParticipant(ctx context.Context, req *participantdto.SubmitParticipantRequest) (*participantdto.ParticipantResponse, error) {
	var result *participantdto.ParticipantResponse

	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		participant, err := uc.participantRepo.GetByID(txCtx, req.ParticipantID)
		if err != nil {
			return fmt.Errorf("get participant: %w", err)
		}

		if err := validateParticipantOwnership(participant, req.TenantID, req.ProductID); err != nil {
			return err
		}

		if !participant.CanBeSubmitted() {
			return errors.ErrBadRequest(fmt.Sprintf("participant in %s status cannot be submitted", participant.Status))
		}

		now := time.Now()
		fromStatus := string(participant.Status)

		participant.Status = entity.ParticipantStatusPendingApproval
		participant.SubmittedBy = &req.UserID
		participant.SubmittedAt = &now

		if err := uc.participantRepo.Update(txCtx, participant); err != nil {
			return fmt.Errorf("update participant: %w", err)
		}

		history := &entity.ParticipantStatusHistory{
			ParticipantID: participant.ID,
			FromStatus:    &fromStatus,
			ToStatus:      string(entity.ParticipantStatusPendingApproval),
			ChangedBy:     req.UserID,
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
