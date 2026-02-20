package internal

import (
	"context"
	"fmt"

	"iam-service/saving/participant/participantdto"
)

func (uc *usecase) UpdatePersonalData(ctx context.Context, req *participantdto.UpdatePersonalDataRequest) (*participantdto.ParticipantResponse, error) {
	var result *participantdto.ParticipantResponse

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
