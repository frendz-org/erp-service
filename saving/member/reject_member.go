package member

import (
	"context"
	"encoding/json"

	"erp-service/entity"
	"erp-service/pkg/errors"
)

func (uc *usecase) RejectMember(ctx context.Context, req *RejectRequest) (*MemberDetailResponse, error) {
	reg, err := uc.utrRepo.GetByID(ctx, req.MemberID)
	if err != nil {
		return nil, err
	}

	if !uc.validateTenantBoundary(reg, req.TenantID, req.ProductID) {
		return nil, errors.ErrNotFound("member not found")
	}

	if reg.Status != entity.UTRStatusPendingApproval {
		return nil, errors.ErrBadRequest("only pending members can be rejected")
	}

	if reg.UserID == req.ApproverID {
		return nil, errors.ErrForbidden("you cannot reject your own registration")
	}

	metadata := map[string]interface{}{
		"rejection_reason": req.Reason,
		"rejected_by":      req.ApproverID.String(),
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, errors.ErrInternal("failed to marshal metadata")
	}

	reg.Status = entity.UTRStatusRejected
	reg.Metadata = metadataJSON

	if err := uc.utrRepo.UpdateStatus(ctx, reg); err != nil {
		return nil, err
	}

	profile, profileErr := uc.profileRepo.GetByUserID(ctx, reg.UserID)
	if profileErr != nil && !errors.IsNotFound(profileErr) {
		return nil, profileErr
	}
	email, emailErr := uc.getUserEmail(ctx, reg.UserID)
	if emailErr != nil {
		return nil, emailErr
	}

	return mapToDetailResponse(reg, profile, email, nil, nil), nil
}
