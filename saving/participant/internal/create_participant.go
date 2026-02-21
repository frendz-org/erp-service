package internal

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant/participantdto"
)

func (uc *usecase) CreateParticipant(ctx context.Context, req *participantdto.CreateParticipantRequest) (*participantdto.ParticipantResponse, error) {
	var result *participantdto.ParticipantResponse

	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {

		existing, err := uc.participantRepo.GetByKTPNumber(txCtx, req.TenantID, req.ProductID, req.KTPNumber)
		if err != nil && !errors.IsNotFound(err) {
			return fmt.Errorf("check ktp number: %w", err)
		}
		if existing != nil {
			if existing.Status == entity.ParticipantStatusDraft {
				return errors.ErrParticipantDraftExists(existing.ID)
			}
			return errors.ErrParticipantAlreadyRegistered()
		}

		existingByEmp, err := uc.participantRepo.GetByEmployeeNumber(txCtx, req.TenantID, req.ProductID, req.EmployeeNumber)
		if err != nil && !errors.IsNotFound(err) {
			return fmt.Errorf("check employee number: %w", err)
		}
		if existingByEmp != nil {
			if existingByEmp.Status == entity.ParticipantStatusDraft {
				return errors.ErrParticipantDraftExists(existingByEmp.ID)
			}
			return errors.ErrParticipantAlreadyRegistered()
		}

		now := time.Now()
		ktpNumber := req.KTPNumber
		employeeNumber := req.EmployeeNumber

		participant := &entity.Participant{
			TenantID:       req.TenantID,
			ProductID:      req.ProductID,
			UserID:         &req.UserID,
			FullName:       req.FullName,
			KTPNumber:      &ktpNumber,
			EmployeeNumber: &employeeNumber,
			Status:         entity.ParticipantStatusDraft,
			CreatedBy:      req.UserID,
			Version:        1,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if err := uc.participantRepo.Create(txCtx, participant); err != nil {
			return fmt.Errorf("create participant: %w", err)
		}

		history := &entity.ParticipantStatusHistory{
			ParticipantID: participant.ID,
			FromStatus:    nil,
			ToStatus:      string(entity.ParticipantStatusDraft),
			ChangedBy:     req.UserID,
			ChangedAt:     now,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := uc.statusHistoryRepo.Create(txCtx, history); err != nil {
			return fmt.Errorf("create status history: %w", err)
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
