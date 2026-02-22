package participant

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
)

func (uc *usecase) SaveAddresses(ctx context.Context, req *SaveAddressesRequest) ([]AddressResponse, error) {
	if len(req.Addresses) == 0 {
		return nil, errors.ErrBadRequest("addresses must contain at least one item")
	}

	var result []AddressResponse

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

		if err := uc.addressRepo.SoftDeleteAllByParticipantID(txCtx, req.ParticipantID); err != nil {
			return fmt.Errorf("soft delete addresses: %w", err)
		}

		now := time.Now()
		insertedAddresses := make([]*entity.ParticipantAddress, 0, len(req.Addresses))
		for _, item := range req.Addresses {
			address := &entity.ParticipantAddress{
				ParticipantID:   req.ParticipantID,
				AddressType:     item.AddressType,
				CountryCode:     item.CountryCode,
				ProvinceCode:    item.ProvinceCode,
				CityCode:        item.CityCode,
				DistrictCode:    item.DistrictCode,
				SubdistrictCode: item.SubdistrictCode,
				PostalCode:      item.PostalCode,
				RT:              item.RT,
				RW:              item.RW,
				AddressLine:     item.AddressLine,
				IsPrimary:       item.IsPrimary,
				Version:         1,
				CreatedAt:       now,
				UpdatedAt:       now,
			}
			if err := uc.addressRepo.Create(txCtx, address); err != nil {
				return fmt.Errorf("create address: %w", err)
			}
			insertedAddresses = append(insertedAddresses, address)
		}

		if participant.StepsCompleted == nil {
			participant.StepsCompleted = make(map[string]bool)
		}
		participant.StepsCompleted["address"] = true
		participant.UpdatedAt = now
		if err := uc.participantRepo.Update(txCtx, participant); err != nil {
			return fmt.Errorf("update participant steps: %w", err)
		}

		result = make([]AddressResponse, 0, len(insertedAddresses))
		for _, a := range insertedAddresses {
			result = append(result, mapAddressToResponse(a))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
