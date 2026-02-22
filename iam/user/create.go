package user

import (
	"context"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	tenant, err := uc.TenantRepo.GetByID(ctx, req.TenantID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrTenantNotFound()
		}
		return nil, errors.ErrInternal("failed to verify tenant").WithError(err)
	}
	if !tenant.IsActive() {
		return nil, errors.ErrTenantInactive()
	}

	role, err := uc.RoleRepo.GetByCode(ctx, req.TenantID, req.RoleCode)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrRoleNotFound()
		}
		return nil, err
	}

	if !role.IsSystem {
		return nil, errors.ErrBadRequest("Only system roles can be assigned through this endpoint")
	}

	emailExists, err := uc.UserRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to check email").WithError(err)
	}
	if emailExists {
		return nil, errors.ErrUserAlreadyExists()
	}

	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal("failed to hash password").WithError(err)
	}

	var response *CreateResponse
	now := time.Now()
	passwordHashStr := string(passwordHash)

	err = uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		user := &entity.User{
			Email:              req.Email,
			Status:             entity.UserStatusActive,
			StatusChangedAt:    &now,
			RegistrationSource: "ADMIN",
		}
		if err := uc.UserRepo.Create(txCtx, user); err != nil {
			return err
		}

		authMethod := entity.NewPasswordAuthMethod(user.ID, passwordHashStr)
		if err := uc.UserAuthMethodRepo.Create(txCtx, authMethod); err != nil {
			return err
		}

		profile := &entity.UserProfile{
			UserID:    user.ID,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			UpdatedAt: now,
		}
		if err := uc.UserProfileRepo.Create(txCtx, profile); err != nil {
			return err
		}

		userRole := &entity.UserRole{
			UserID:     user.ID,
			RoleID:     role.ID,
			AssignedAt: now,
			CreatedAt:  now,
		}
		if err := uc.UserRoleRepo.Create(txCtx, userRole); err != nil {
			return err
		}

		securityState := &entity.UserSecurityState{
			UserID:          user.ID,
			EmailVerified:   true,
			EmailVerifiedAt: &now,
			UpdatedAt:       now,
		}
		if err := uc.UserSecurityStateRepo.Create(txCtx, securityState); err != nil {
			return err
		}

		response = &CreateResponse{
			UserID:   user.ID,
			Email:    req.Email,
			FullName: req.FirstName + " " + req.LastName,
			RoleCode: req.RoleCode,
			TenantID: req.TenantID,
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrInternal("failed to create user").WithError(err)
	}

	return response, nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.ErrValidation("password must be at least 8 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*':
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.ErrValidation("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.ErrValidation("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.ErrValidation("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.ErrValidation("password must contain at least one special character")
	}

	return nil
}
