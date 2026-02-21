package internal

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant/participantdto"
)

func (uc *usecase) SaveIdentity(ctx context.Context, req *participantdto.SaveIdentityRequest) (*participantdto.IdentityResponse, error) {
	var result *participantdto.IdentityResponse

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

		var identity *entity.ParticipantIdentity

		if req.ID != nil {
			identity, err = uc.identityRepo.GetByID(txCtx, *req.ID)
			if err != nil {
				return fmt.Errorf("get identity: %w", err)
			}

			if identity.ParticipantID != req.ParticipantID {
				return errors.ErrForbidden("identity does not belong to this participant")
			}

			identity.IdentityType = req.IdentityType
			identity.IdentityNumber = req.IdentityNumber
			identity.IdentityAuthority = req.IdentityAuthority
			identity.IssueDate = req.IssueDate
			identity.ExpiryDate = req.ExpiryDate
			identity.PhotoFilePath = req.PhotoFilePath

			if err := uc.identityRepo.Update(txCtx, identity); err != nil {
				return fmt.Errorf("update identity: %w", err)
			}
		} else {
			now := time.Now()
			identity = &entity.ParticipantIdentity{
				ParticipantID:     req.ParticipantID,
				IdentityType:      req.IdentityType,
				IdentityNumber:    req.IdentityNumber,
				IdentityAuthority: req.IdentityAuthority,
				IssueDate:         req.IssueDate,
				ExpiryDate:        req.ExpiryDate,
				PhotoFilePath:     req.PhotoFilePath,
				Version:           1,
				CreatedAt:         now,
				UpdatedAt:         now,
			}

			if err := uc.identityRepo.Create(txCtx, identity); err != nil {
				return fmt.Errorf("create identity: %w", err)
			}
		}

		resp := mapIdentityToResponse(identity)
		result = &resp
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
