package internal

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant/participantdto"
)

func (uc *usecase) SaveBeneficiary(ctx context.Context, req *participantdto.SaveBeneficiaryRequest) (*participantdto.BeneficiaryResponse, error) {
	var result *participantdto.BeneficiaryResponse

	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
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

		familyMember, err := uc.familyMemberRepo.GetByID(txCtx, req.FamilyMemberID)
		if err != nil {
			return fmt.Errorf("validate family member exists: %w", err)
		}

		if familyMember.ParticipantID != req.ParticipantID {
			return errors.ErrForbidden("family member does not belong to this participant")
		}

		var beneficiary *entity.ParticipantBeneficiary

		if req.ID != nil {
			beneficiary, err = uc.beneficiaryRepo.GetByID(txCtx, *req.ID)
			if err != nil {
				return fmt.Errorf("get beneficiary: %w", err)
			}

			if beneficiary.ParticipantID != req.ParticipantID {
				return errors.ErrForbidden("beneficiary does not belong to this participant")
			}

			beneficiary.FamilyMemberID = req.FamilyMemberID
			beneficiary.IdentityPhotoFilePath = req.IdentityPhotoFilePath
			beneficiary.FamilyCardPhotoFilePath = req.FamilyCardPhotoFilePath
			beneficiary.BankBookPhotoFilePath = req.BankBookPhotoFilePath
			beneficiary.AccountNumber = req.AccountNumber

			if err := uc.beneficiaryRepo.Update(txCtx, beneficiary); err != nil {
				return fmt.Errorf("update beneficiary: %w", err)
			}
		} else {
			now := time.Now()
			beneficiary = &entity.ParticipantBeneficiary{
				ParticipantID:           req.ParticipantID,
				FamilyMemberID:          req.FamilyMemberID,
				IdentityPhotoFilePath:   req.IdentityPhotoFilePath,
				FamilyCardPhotoFilePath: req.FamilyCardPhotoFilePath,
				BankBookPhotoFilePath:   req.BankBookPhotoFilePath,
				AccountNumber:           req.AccountNumber,
				Version:                 1,
				CreatedAt:               now,
				UpdatedAt:               now,
			}

			if err := uc.beneficiaryRepo.Create(txCtx, beneficiary); err != nil {
				return fmt.Errorf("create beneficiary: %w", err)
			}
		}

		resp := mapBeneficiaryToResponse(beneficiary)
		result = &resp
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
