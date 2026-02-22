package participant

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
)

func (uc *usecase) SaveBeneficiary(ctx context.Context, req *SaveBeneficiaryRequest) (*BeneficiaryResponse, error) {
	var result *BeneficiaryResponse

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

		familyMember, err := uc.familyMemberRepo.GetByID(txCtx, req.FamilyMemberID)
		if err != nil {
			return fmt.Errorf("validate family member exists: %w", err)
		}

		if familyMember.ParticipantID != req.ParticipantID {
			return errors.ErrForbidden("family member does not belong to this participant")
		}

		if req.IdentityPhotoFileID != nil {
			file, err := uc.fileRepo.GetByID(txCtx, *req.IdentityPhotoFileID)
			if err != nil {
				return fmt.Errorf("get identity photo file: %w", err)
			}
			if file.TenantID != req.TenantID || file.ProductID != req.ProductID {
				return errors.ErrForbidden("identity_photo_file_id does not belong to this tenant/product")
			}
			if err := uc.fileRepo.SetPermanent(txCtx, *req.IdentityPhotoFileID); err != nil {
				return fmt.Errorf("set identity photo file permanent: %w", err)
			}
		}

		if req.FamilyCardPhotoFileID != nil {
			file, err := uc.fileRepo.GetByID(txCtx, *req.FamilyCardPhotoFileID)
			if err != nil {
				return fmt.Errorf("get family card photo file: %w", err)
			}
			if file.TenantID != req.TenantID || file.ProductID != req.ProductID {
				return errors.ErrForbidden("family_card_photo_file_id does not belong to this tenant/product")
			}
			if err := uc.fileRepo.SetPermanent(txCtx, *req.FamilyCardPhotoFileID); err != nil {
				return fmt.Errorf("set family card photo file permanent: %w", err)
			}
		}

		if req.BankBookPhotoFileID != nil {
			file, err := uc.fileRepo.GetByID(txCtx, *req.BankBookPhotoFileID)
			if err != nil {
				return fmt.Errorf("get bank book photo file: %w", err)
			}
			if file.TenantID != req.TenantID || file.ProductID != req.ProductID {
				return errors.ErrForbidden("bank_book_photo_file_id does not belong to this tenant/product")
			}
			if err := uc.fileRepo.SetPermanent(txCtx, *req.BankBookPhotoFileID); err != nil {
				return fmt.Errorf("set bank book photo file permanent: %w", err)
			}
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
			beneficiary.IdentityPhotoFileID = req.IdentityPhotoFileID
			beneficiary.FamilyCardPhotoFileID = req.FamilyCardPhotoFileID
			beneficiary.BankBookPhotoFileID = req.BankBookPhotoFileID
			beneficiary.AccountNumber = req.AccountNumber

			if err := uc.beneficiaryRepo.Update(txCtx, beneficiary); err != nil {
				return fmt.Errorf("update beneficiary: %w", err)
			}
		} else {
			now := time.Now()
			beneficiary = &entity.ParticipantBeneficiary{
				ParticipantID:        req.ParticipantID,
				FamilyMemberID:       req.FamilyMemberID,
				IdentityPhotoFileID:  req.IdentityPhotoFileID,
				FamilyCardPhotoFileID: req.FamilyCardPhotoFileID,
				BankBookPhotoFileID:  req.BankBookPhotoFileID,
				AccountNumber:        req.AccountNumber,
				Version:              1,
				CreatedAt:            now,
				UpdatedAt:            now,
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
