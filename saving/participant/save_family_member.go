package participant

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
)

func (uc *usecase) SaveFamilyMember(ctx context.Context, req *SaveFamilyMemberRequest) (*FamilyMemberResponse, error) {
	var result *FamilyMemberResponse

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

		if req.SupportingDocFileID != nil {
			file, err := uc.fileRepo.GetByID(txCtx, *req.SupportingDocFileID)
			if err != nil {
				return fmt.Errorf("get supporting doc file: %w", err)
			}
			if file.TenantID != req.TenantID || file.ProductID != req.ProductID {
				return errors.ErrForbidden("supporting_doc_file_id does not belong to this tenant/product")
			}
			if err := uc.fileRepo.SetPermanent(txCtx, *req.SupportingDocFileID); err != nil {
				return fmt.Errorf("set supporting doc file permanent: %w", err)
			}
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
			member.SupportingDocFileID = req.SupportingDocFileID

			if err := uc.familyMemberRepo.Update(txCtx, member); err != nil {
				return fmt.Errorf("update family member: %w", err)
			}
		} else {
			now := time.Now()
			member = &entity.ParticipantFamilyMember{
				ParticipantID:       req.ParticipantID,
				FullName:            req.FullName,
				RelationshipType:    req.RelationshipType,
				IsDependent:         req.IsDependent,
				SupportingDocFileID: req.SupportingDocFileID,
				Version:             1,
				CreatedAt:           now,
				UpdatedAt:           now,
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
