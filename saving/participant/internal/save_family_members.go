package internal

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant/participantdto"

	"github.com/google/uuid"
)

func (uc *usecase) validateFileOwnership(ctx context.Context, fileID uuid.UUID, tenantID, productID uuid.UUID) error {
	file, err := uc.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return fmt.Errorf("get file: %w", err)
	}
	if file.TenantID != tenantID || file.ProductID != productID {
		return errors.ErrForbidden("file does not belong to this tenant/product")
	}
	return nil
}

func (uc *usecase) SaveFamilyMembers(ctx context.Context, req *participantdto.SaveFamilyMembersRequest) ([]participantdto.FamilyMemberResponse, error) {
	if len(req.FamilyMembers) == 0 {
		return nil, errors.ErrBadRequest("family_members must contain at least one item")
	}

	var result []participantdto.FamilyMemberResponse

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

		if err := uc.familyMemberRepo.SoftDeleteAllByParticipantID(txCtx, req.ParticipantID); err != nil {
			return fmt.Errorf("soft delete family members: %w", err)
		}

		now := time.Now()
		insertedMembers := make([]*entity.ParticipantFamilyMember, 0, len(req.FamilyMembers))
		for _, item := range req.FamilyMembers {

			if item.SupportingDocFileID != nil {
				if err := uc.validateFileOwnership(txCtx, *item.SupportingDocFileID, req.TenantID, req.ProductID); err != nil {
					return fmt.Errorf("validate supporting_doc_file_id: %w", err)
				}
			}

			member := &entity.ParticipantFamilyMember{
				ParticipantID:       req.ParticipantID,
				FullName:            item.FullName,
				RelationshipType:    item.RelationshipType,
				IsDependent:         item.IsDependent,
				SupportingDocFileID: item.SupportingDocFileID,
				Version:             1,
				CreatedAt:           now,
				UpdatedAt:           now,
			}
			if err := uc.familyMemberRepo.Create(txCtx, member); err != nil {
				return fmt.Errorf("create family member: %w", err)
			}
			insertedMembers = append(insertedMembers, member)

			if item.SupportingDocFileID != nil {
				if err := uc.fileRepo.SetPermanent(txCtx, *item.SupportingDocFileID); err != nil {
					return fmt.Errorf("set supporting_doc_file permanent: %w", err)
				}
			}
		}

		if participant.StepsCompleted == nil {
			participant.StepsCompleted = make(map[string]bool)
		}
		participant.StepsCompleted["family_members"] = true
		participant.UpdatedAt = now
		if err := uc.participantRepo.Update(txCtx, participant); err != nil {
			return fmt.Errorf("update participant steps: %w", err)
		}

		result = make([]participantdto.FamilyMemberResponse, 0, len(insertedMembers))
		for _, m := range insertedMembers {
			result = append(result, mapFamilyMemberToResponse(m))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
