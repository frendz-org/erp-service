package participant

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
)

func (uc *usecase) SaveBeneficiaries(ctx context.Context, req *SaveBeneficiariesRequest) ([]BeneficiaryResponse, error) {
	if len(req.Beneficiaries) == 0 {
		return nil, errors.ErrBadRequest("beneficiaries must contain at least one item")
	}

	var result []BeneficiaryResponse

	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
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

		if err := uc.beneficiaryRepo.SoftDeleteAllByParticipantID(txCtx, req.ParticipantID); err != nil {
			return fmt.Errorf("soft delete beneficiaries: %w", err)
		}

		now := time.Now()
		insertedBeneficiaries := make([]*entity.ParticipantBeneficiary, 0, len(req.Beneficiaries))
		for _, item := range req.Beneficiaries {

			fm, err := uc.familyMemberRepo.GetByID(txCtx, item.FamilyMemberID)
			if err != nil {
				return fmt.Errorf("validate family_member_id: %w", err)
			}
			if fm.ParticipantID != req.ParticipantID {
				return errors.ErrForbidden("family_member_id does not belong to this participant")
			}

			if item.IdentityPhotoFileID != nil {
				if err := uc.validateFileOwnership(txCtx, *item.IdentityPhotoFileID, req.TenantID, req.ProductID); err != nil {
					return fmt.Errorf("validate identity_photo_file_id: %w", err)
				}
			}
			if item.FamilyCardPhotoFileID != nil {
				if err := uc.validateFileOwnership(txCtx, *item.FamilyCardPhotoFileID, req.TenantID, req.ProductID); err != nil {
					return fmt.Errorf("validate family_card_photo_file_id: %w", err)
				}
			}
			if item.BankBookPhotoFileID != nil {
				if err := uc.validateFileOwnership(txCtx, *item.BankBookPhotoFileID, req.TenantID, req.ProductID); err != nil {
					return fmt.Errorf("validate bank_book_photo_file_id: %w", err)
				}
			}

			beneficiary := &entity.ParticipantBeneficiary{
				ParticipantID:         req.ParticipantID,
				FamilyMemberID:        item.FamilyMemberID,
				IdentityPhotoFileID:   item.IdentityPhotoFileID,
				FamilyCardPhotoFileID: item.FamilyCardPhotoFileID,
				BankBookPhotoFileID:   item.BankBookPhotoFileID,
				AccountNumber:         item.AccountNumber,
				Version:               1,
				CreatedAt:             now,
				UpdatedAt:             now,
			}
			if err := uc.beneficiaryRepo.Create(txCtx, beneficiary); err != nil {
				return fmt.Errorf("create beneficiary: %w", err)
			}
			insertedBeneficiaries = append(insertedBeneficiaries, beneficiary)

			if item.IdentityPhotoFileID != nil {
				if err := uc.fileRepo.SetPermanent(txCtx, *item.IdentityPhotoFileID); err != nil {
					return fmt.Errorf("set identity_photo_file permanent: %w", err)
				}
			}
			if item.FamilyCardPhotoFileID != nil {
				if err := uc.fileRepo.SetPermanent(txCtx, *item.FamilyCardPhotoFileID); err != nil {
					return fmt.Errorf("set family_card_photo_file permanent: %w", err)
				}
			}
			if item.BankBookPhotoFileID != nil {
				if err := uc.fileRepo.SetPermanent(txCtx, *item.BankBookPhotoFileID); err != nil {
					return fmt.Errorf("set bank_book_photo_file permanent: %w", err)
				}
			}
		}

		if participant.StepsCompleted == nil {
			participant.StepsCompleted = make(map[string]bool)
		}
		participant.StepsCompleted["beneficiaries"] = true
		participant.UpdatedAt = now
		if err := uc.participantRepo.Update(txCtx, participant); err != nil {
			return fmt.Errorf("update participant steps: %w", err)
		}

		result = make([]BeneficiaryResponse, 0, len(insertedBeneficiaries))
		for _, b := range insertedBeneficiaries {
			result = append(result, mapBeneficiaryToResponse(b))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
