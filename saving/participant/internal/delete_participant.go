package internal

import (
	"context"
	"fmt"

	"iam-service/pkg/errors"
	"iam-service/saving/participant/participantdto"
)

func (uc *usecase) DeleteParticipant(ctx context.Context, req *participantdto.DeleteParticipantRequest) error {
	return uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		participant, err := uc.participantRepo.GetByID(txCtx, req.ParticipantID)
		if err != nil {
			return fmt.Errorf("get participant: %w", err)
		}

		if err := validateParticipantOwnership(participant, req.TenantID, req.ProductID); err != nil {
			return err
		}

		if !participant.IsDraft() {
			return errors.ErrBadRequest("only DRAFT participants can be deleted")
		}

		if err := uc.participantRepo.SoftDelete(txCtx, req.ParticipantID); err != nil {
			return fmt.Errorf("soft delete participant: %w", err)
		}

		return nil
	})
}
