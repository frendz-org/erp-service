package participant

import (
	"context"
	"fmt"

	"erp-service/pkg/errors"
)

func (uc *usecase) DeleteBankAccount(ctx context.Context, req *DeleteChildEntityRequest) error {
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

		account, err := uc.bankAccountRepo.GetByID(txCtx, req.ChildID)
		if err != nil {
			return fmt.Errorf("get bank account: %w", err)
		}
		if account.ParticipantID != req.ParticipantID {
			return errors.ErrForbidden("bank account does not belong to this participant")
		}

		if err := uc.bankAccountRepo.SoftDelete(txCtx, req.ChildID); err != nil {
			return fmt.Errorf("soft delete bank account: %w", err)
		}

		return nil
	})
}
