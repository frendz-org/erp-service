package participant

import (
	"context"
	"fmt"

	"erp-service/pkg/errors"
)

func (uc *usecase) DeleteIdentity(ctx context.Context, req *DeleteChildEntityRequest) error {
	return uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		participant, err := uc.participantRepo.GetByID(txCtx, req.ParticipantID)
		if err != nil {
			return fmt.Errorf("get participant: %w", err)
		}

		if err := ValidateParticipantOwnership(participant, req.TenantID, req.ProductID); err != nil {
			return err
		}

		if err := ValidateEditableState(participant); err != nil {
			return err
		}

		identity, err := uc.identityRepo.GetByID(txCtx, req.ChildID)
		if err != nil {
			return fmt.Errorf("get identity: %w", err)
		}
		if identity.ParticipantID != req.ParticipantID {
			return errors.ErrForbidden("identity does not belong to this participant")
		}

		if err := uc.identityRepo.SoftDelete(txCtx, req.ChildID); err != nil {
			return fmt.Errorf("soft delete identity: %w", err)
		}

		return nil
	})
}
