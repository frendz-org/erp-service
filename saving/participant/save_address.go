package participant

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
)

func (uc *usecase) SaveAddress(ctx context.Context, req *SaveAddressRequest) (*AddressResponse, error) {
	var result *AddressResponse

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

		var address *entity.ParticipantAddress

		if req.ID != nil {
			address, err = uc.addressRepo.GetByID(txCtx, *req.ID)
			if err != nil {
				return fmt.Errorf("get address: %w", err)
			}

			if address.ParticipantID != req.ParticipantID {
				return errors.ErrForbidden("address does not belong to this participant")
			}

			address.AddressType = req.AddressType
			address.CountryCode = req.CountryCode
			address.ProvinceCode = req.ProvinceCode
			address.CityCode = req.CityCode
			address.DistrictCode = req.DistrictCode
			address.SubdistrictCode = req.SubdistrictCode
			address.PostalCode = req.PostalCode
			address.RT = req.RT
			address.RW = req.RW
			address.AddressLine = req.AddressLine
			address.IsPrimary = req.IsPrimary

			if err := uc.addressRepo.Update(txCtx, address); err != nil {
				return fmt.Errorf("update address: %w", err)
			}
		} else {
			now := time.Now()
			address = &entity.ParticipantAddress{
				ParticipantID:   req.ParticipantID,
				AddressType:     req.AddressType,
				CountryCode:     req.CountryCode,
				ProvinceCode:    req.ProvinceCode,
				CityCode:        req.CityCode,
				DistrictCode:    req.DistrictCode,
				SubdistrictCode: req.SubdistrictCode,
				PostalCode:      req.PostalCode,
				RT:              req.RT,
				RW:              req.RW,
				AddressLine:     req.AddressLine,
				IsPrimary:       req.IsPrimary,
				Version:         1,
				CreatedAt:       now,
				UpdatedAt:       now,
			}

			if err := uc.addressRepo.Create(txCtx, address); err != nil {
				return fmt.Errorf("create address: %w", err)
			}
		}

		resp := mapAddressToResponse(address)
		result = &resp
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
