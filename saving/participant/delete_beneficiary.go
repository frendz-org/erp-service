package participant

import (
	"context"
	"fmt"

	"erp-service/pkg/errors"
)

func (uc *usecase) DeleteBeneficiary(ctx context.Context, req *DeleteChildEntityRequest) error {
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

		beneficiary, err := uc.beneficiaryRepo.GetByID(txCtx, req.ChildID)
		if err != nil {
			return fmt.Errorf("get beneficiary: %w", err)
		}
		if beneficiary.ParticipantID != req.ParticipantID {
			return errors.ErrForbidden("beneficiary does not belong to this participant")
		}

		if err := uc.beneficiaryRepo.SoftDelete(txCtx, req.ChildID); err != nil {
			return fmt.Errorf("soft delete beneficiary: %w", err)
		}

		return nil
	})
}
