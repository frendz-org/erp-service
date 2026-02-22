package participant

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
)

func (uc *usecase) SavePension(ctx context.Context, req *SavePensionRequest) (*PensionResponse, error) {
	var result *PensionResponse

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

		var pension *entity.ParticipantPension

		if req.ID != nil {
			pension, err = uc.pensionRepo.GetByID(txCtx, *req.ID)
			if err != nil {
				return fmt.Errorf("get pension: %w", err)
			}

			if pension.ParticipantID != req.ParticipantID {
				return errors.ErrForbidden("pension does not belong to this participant")
			}
		} else {
			pension, err = uc.pensionRepo.GetByParticipantID(txCtx, req.ParticipantID)
			if err != nil && !errors.IsNotFound(err) {
				return fmt.Errorf("get existing pension: %w", err)
			}
		}

		if pension != nil {
			pension.ParticipantNumber = req.ParticipantNumber
			pension.PensionCategory = req.PensionCategory
			pension.PensionStatus = req.PensionStatus
			pension.EffectiveDate = req.EffectiveDate
			pension.EndDate = req.EndDate
			pension.ProjectedRetirementDate = req.ProjectedRetirementDate
			pension.UpdatedAt = time.Now()

			if err := uc.pensionRepo.Update(txCtx, pension); err != nil {
				return fmt.Errorf("update pension: %w", err)
			}
		} else {
			now := time.Now()
			pension = &entity.ParticipantPension{
				ParticipantID:           req.ParticipantID,
				ParticipantNumber:       req.ParticipantNumber,
				PensionCategory:         req.PensionCategory,
				PensionStatus:           req.PensionStatus,
				EffectiveDate:           req.EffectiveDate,
				EndDate:                 req.EndDate,
				ProjectedRetirementDate: req.ProjectedRetirementDate,
				Version:                 1,
				CreatedAt:               now,
				UpdatedAt:               now,
			}

			if err := uc.pensionRepo.Create(txCtx, pension); err != nil {
				return fmt.Errorf("create pension: %w", err)
			}
		}

		resp := mapPensionToResponse(pension)
		result = &resp
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
