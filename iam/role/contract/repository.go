package contract

import (
	"context"

	"iam-service/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	EmailExists(ctx context.Context, email string) (bool, error)
}
type UserProfileRepository interface {
	Create(ctx context.Context, profile *entity.UserProfile) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error)
	Update(ctx context.Context, profile *entity.UserProfile) error
}
type UserAuthMethodRepository interface {
	Create(ctx context.Context, authMethod *entity.UserAuthMethod) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserAuthMethod, error)
	Update(ctx context.Context, authMethod *entity.UserAuthMethod) error
}
type UserSecurityStateRepository interface {
	Create(ctx context.Context, securityState *entity.UserSecurityState) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserSecurityState, error)
	Update(ctx context.Context, securityState *entity.UserSecurityState) error
}
type EmailVerificationRepository interface {
	Create(ctx context.Context, verification *entity.EmailVerification) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.EmailVerification, error)
	GetLatestByEmail(ctx context.Context, email string, otpType entity.OTPType) (*entity.EmailVerification, error)
	GetLatestByUserID(ctx context.Context, userID uuid.UUID, otpType entity.OTPType) (*entity.EmailVerification, error)
	MarkAsVerified(ctx context.Context, id uuid.UUID) error
	CountActiveOTPsByEmail(ctx context.Context, email string, otpType entity.OTPType) (int, error)
	DeleteExpiredByEmail(ctx context.Context, email string) error
}
type TenantRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error)
	GetByCode(ctx context.Context, code string) (*entity.Tenant, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	GetByName(ctx context.Context, productID uuid.UUID, name string) (*entity.Role, error)
	GetByCode(ctx context.Context, productID uuid.UUID, code string) (*entity.Role, error)
}
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID, reason string) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID, reason string) error
	RevokeByFamily(ctx context.Context, tokenFamily uuid.UUID, reason string) error
}
type PINVerificationLogRepository interface {
	Create(ctx context.Context, log *entity.PINVerificationLog) error
	CountRecentFailures(ctx context.Context, userID uuid.UUID, since int) (int, error)
}

type RolePermissionRepository interface {
	Create(ctx context.Context, rolePermission *entity.RolePermission) error
}
