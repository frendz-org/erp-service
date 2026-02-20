package internal

import (
	"context"
	"fmt"
	"time"

	"iam-service/entity"
	"iam-service/saving/participant/participantdto"
	"iam-service/pkg/errors"
)

func (uc *usecase) SaveFamilyMember(ctx context.Context, req *participantdto.SaveFamilyMemberRequest) (*participantdto.FamilyMemberResponse, error) {
	var result *participantdto.FamilyMemberResponse

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

		var member *entity.ParticipantFamilyMember

		if req.ID != nil {
			member, err = uc.familyMemberRepo.GetByID(txCtx, *req.ID)
			if err != nil {
				return fmt.Errorf("get family member: %w", err)
			}

			if member.ParticipantID != req.ParticipantID {
				return errors.ErrForbidden("family member does not belong to this participant")
			}

			member.FullName = req.FullName
			member.RelationshipType = req.RelationshipType
			member.IsDependent = req.IsDependent
			member.SupportingDocFilePath = req.SupportingDocFilePath

			if err := uc.familyMemberRepo.Update(txCtx, member); err != nil {
				return fmt.Errorf("update family member: %w", err)
			}
		} else {
			now := time.Now()
			member = &entity.ParticipantFamilyMember{
				ParticipantID:         req.ParticipantID,
				FullName:              req.FullName,
				RelationshipType:      req.RelationshipType,
				IsDependent:           req.IsDependent,
				SupportingDocFilePath: req.SupportingDocFilePath,
				Version:               1,
				CreatedAt:             now,
				UpdatedAt:             now,
			}

			if err := uc.familyMemberRepo.Create(txCtx, member); err != nil {
				return fmt.Errorf("create family member: %w", err)
			}
		}

		resp := mapFamilyMemberToResponse(member)
		result = &resp
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
