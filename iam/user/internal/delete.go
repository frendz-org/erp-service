package internal

import (
	"context"

	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) Delete(ctx context.Context, callerTenantID *uuid.UUID, id uuid.UUID) error {
	_, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return errors.ErrUserNotFound()
		}
		return err
	}

	if err := uc.UserRepo.Delete(ctx, id); err != nil {
		if errors.IsNotFound(err) {
			return errors.ErrUserNotFound()
		}
		return err
	}

	return nil
}
