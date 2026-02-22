package participant

import (
	"context"
	"fmt"

)

func (uc *usecase) UpdatePersonalData(ctx context.Context, req *UpdatePersonalDataRequest) (*ParticipantResponse, error) {
	var result *ParticipantResponse

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

		participant.FullName = req.FullName
		participant.Gender = req.Gender
		participant.PlaceOfBirth = req.PlaceOfBirth
		participant.DateOfBirth = req.DateOfBirth
		participant.MaritalStatus = req.MaritalStatus
		participant.Citizenship = req.Citizenship
		participant.Religion = req.Religion
		participant.KTPNumber = req.KTPNumber
		participant.EmployeeNumber = req.EmployeeNumber
		participant.PhoneNumber = req.PhoneNumber

		if err := uc.participantRepo.Update(txCtx, participant); err != nil {
			return fmt.Errorf("update participant: %w", err)
		}

		resp, err := uc.buildFullParticipantResponse(txCtx, participant, false)
		if err != nil {
			return fmt.Errorf("build response: %w", err)
		}
		result = resp
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
