package internal

import (
	"context"
	"fmt"
	"time"

	"iam-service/entity"
	"iam-service/pkg/errors"
	"iam-service/saving/participant/participantdto"
)

func (uc *usecase) ApproveParticipant(ctx context.Context, req *participantdto.ApproveParticipantRequest) (*participantdto.ParticipantResponse, error) {
	var result *participantdto.ParticipantResponse

	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		participant, err := uc.participantRepo.GetByID(txCtx, req.ParticipantID)
		if err != nil {
			return fmt.Errorf("get participant: %w", err)
		}

		if err := validateParticipantOwnership(participant, req.TenantID, req.ProductID); err != nil {
			return err
		}

		if !participant.CanBeApproved() {
			return errors.ErrBadRequest(fmt.Sprintf("participant in %s status cannot be approved", participant.Status))
		}

		now := time.Now()
		fromStatus := string(participant.Status)

		participant.Status = entity.ParticipantStatusApproved
		participant.ApprovedBy = &req.UserID
		participant.ApprovedAt = &now

		if err := uc.participantRepo.Update(txCtx, participant); err != nil {
			return fmt.Errorf("update participant: %w", err)
		}

		history := &entity.ParticipantStatusHistory{
			ParticipantID: participant.ID,
			FromStatus:    &fromStatus,
			ToStatus:      string(entity.ParticipantStatusApproved),
			ChangedBy:     req.UserID,
			ChangedAt:     now,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := uc.statusHistoryRepo.Create(txCtx, history); err != nil {
			return fmt.Errorf("create status history: %w", err)
		}

		resp, err := uc.buildFullParticipantResponse(txCtx, participant, true)
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
