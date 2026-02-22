package participant

import (
	"context"
	"fmt"

	"erp-service/pkg/errors"
)

func (uc *usecase) DeleteAddress(ctx context.Context, req *DeleteChildEntityRequest) error {
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

		address, err := uc.addressRepo.GetByID(txCtx, req.ChildID)
		if err != nil {
			return fmt.Errorf("get address: %w", err)
		}
		if address.ParticipantID != req.ParticipantID {
			return errors.ErrForbidden("address does not belong to this participant")
		}

		if err := uc.addressRepo.SoftDelete(txCtx, req.ChildID); err != nil {
			return fmt.Errorf("soft delete address: %w", err)
		}

		return nil
	})
}
