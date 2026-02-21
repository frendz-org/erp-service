package internal

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant/participantdto"
)

func (uc *usecase) SaveEmployment(ctx context.Context, req *participantdto.SaveEmploymentRequest) (*participantdto.EmploymentResponse, error) {
	var result *participantdto.EmploymentResponse

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

		var employment *entity.ParticipantEmployment

		if req.ID != nil {
			employment, err = uc.employmentRepo.GetByID(txCtx, *req.ID)
			if err != nil {
				return fmt.Errorf("get employment: %w", err)
			}

			if employment.ParticipantID != req.ParticipantID {
				return errors.ErrForbidden("employment does not belong to this participant")
			}
		} else {
			employment, err = uc.employmentRepo.GetByParticipantID(txCtx, req.ParticipantID)
			if err != nil && !errors.IsNotFound(err) {
				return fmt.Errorf("get existing employment: %w", err)
			}
		}

		if employment != nil {
			employment.PersonnelNumber = req.PersonnelNumber
			employment.DateOfHire = req.DateOfHire
			employment.CorporateGroupName = req.CorporateGroupName
			employment.LegalEntityCode = req.LegalEntityCode
			employment.LegalEntityName = req.LegalEntityName
			employment.BusinessUnitCode = req.BusinessUnitCode
			employment.BusinessUnitName = req.BusinessUnitName
			employment.TenantName = req.TenantName
			employment.EmploymentStatus = req.EmploymentStatus
			employment.PositionName = req.PositionName
			employment.JobLevel = req.JobLevel
			employment.LocationCode = req.LocationCode
			employment.LocationName = req.LocationName
			employment.SubLocationName = req.SubLocationName
			employment.RetirementDate = req.RetirementDate
			employment.RetirementTypeCode = req.RetirementTypeCode

			if err := uc.employmentRepo.Update(txCtx, employment); err != nil {
				return fmt.Errorf("update employment: %w", err)
			}
		} else {
			now := time.Now()
			employment = &entity.ParticipantEmployment{
				ParticipantID:      req.ParticipantID,
				PersonnelNumber:    req.PersonnelNumber,
				DateOfHire:         req.DateOfHire,
				CorporateGroupName: req.CorporateGroupName,
				LegalEntityCode:    req.LegalEntityCode,
				LegalEntityName:    req.LegalEntityName,
				BusinessUnitCode:   req.BusinessUnitCode,
				BusinessUnitName:   req.BusinessUnitName,
				TenantName:         req.TenantName,
				EmploymentStatus:   req.EmploymentStatus,
				PositionName:       req.PositionName,
				JobLevel:           req.JobLevel,
				LocationCode:       req.LocationCode,
				LocationName:       req.LocationName,
				SubLocationName:    req.SubLocationName,
				RetirementDate:     req.RetirementDate,
				RetirementTypeCode: req.RetirementTypeCode,
				Version:            1,
				CreatedAt:          now,
				UpdatedAt:          now,
			}

			if err := uc.employmentRepo.Create(txCtx, employment); err != nil {
				return fmt.Errorf("create employment: %w", err)
			}
		}

		resp := mapEmploymentToResponse(employment)
		result = &resp
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
