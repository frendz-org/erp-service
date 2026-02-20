package internal

import (
	"context"
	"fmt"

	"iam-service/pkg/errors"
	"iam-service/saving/participant/participantdto"
)

func (uc *usecase) DeleteFamilyMember(ctx context.Context, req *participantdto.DeleteChildEntityRequest) error {
	return uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		participant, err := uc.participantRepo.GetByID(txCtx, req.ParticipantID)
		if err != nil {
			return fmt.Errorf("get participant: %w", err)
		}

		if err := validateParticipantOwnership(participant, req.TenantID, req.ProductID); err != nil {
			return err
		}

		if err := validateEditableState(participant); err != nil {
			return err
		}

		member, err := uc.familyMemberRepo.GetByID(txCtx, req.ChildID)
		if err != nil {
			return fmt.Errorf("get family member: %w", err)
		}
		if member.ParticipantID != req.ParticipantID {
			return errors.ErrForbidden("family member does not belong to this participant")
		}

		beneficiaries, err := uc.beneficiaryRepo.ListByParticipantID(txCtx, req.ParticipantID)
		if err != nil {
			return fmt.Errorf("list beneficiaries: %w", err)
		}

		for _, b := range beneficiaries {
			if b.FamilyMemberID == req.ChildID {
				return errors.ErrBadRequest("cannot delete family member that is referenced as a beneficiary")
			}
		}

		if err := uc.familyMemberRepo.SoftDelete(txCtx, req.ChildID); err != nil {
			return fmt.Errorf("soft delete family member: %w", err)
		}

		return nil
	})
}
