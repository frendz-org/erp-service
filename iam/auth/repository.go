package auth

import (
	"context"
	"time"

	"erp-service/entity"
	usercontract "erp-service/iam/user"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	EmailExists(ctx context.Context, email string) (bool, error)
	List(ctx context.Context, filter *usercontract.UserListFilter) ([]*entity.User, int64, error)
}
type UserProfileRepository interface {
	Create(ctx context.Context, profile *entity.UserProfile) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error)
	Update(ctx context.Context, profile *entity.UserProfile) error
}
type UserAuthMethodRepository interface {
	Create(ctx context.Context, authMethod *entity.UserAuthMethod) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserAuthMethod, error)
	GetByUserIDAndMethodType(ctx context.Context, userID uuid.UUID, methodType string) (*entity.UserAuthMethod, error)
	GetByCredentialField(ctx context.Context, methodType, jsonField, value string) (*entity.UserAuthMethod, error)
	Update(ctx context.Context, authMethod *entity.UserAuthMethod) error
}
type UserSecurityStateRepository interface {
	Create(ctx context.Context, securityState *entity.UserSecurityState) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserSecurityState, error)
	Update(ctx context.Context, securityState *entity.UserSecurityState) error
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
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.Role, error)
}
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.RefreshToken, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	SetReplacedBy(ctx context.Context, id uuid.UUID, replacedByID uuid.UUID) error
	Revoke(ctx context.Context, id uuid.UUID, reason string) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID, reason string) error
	RevokeByFamily(ctx context.Context, tokenFamily uuid.UUID, reason string) error
}
type UserRoleRepository interface {
	Create(ctx context.Context, userRole *entity.UserRole) error
	ListActiveByUserID(ctx context.Context, userID uuid.UUID, productID *uuid.UUID) ([]entity.UserRole, error)
}

type RolePermissionRepository interface {
	Create(ctx context.Context, rolePermission *entity.RolePermission) error
}

type ProductRepository interface {
	GetByCodeAndTenant(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Product, error)
	GetByIDAndTenant(ctx context.Context, productID, tenantID uuid.UUID) (*entity.Product, error)
}

type PermissionRepository interface {
	GetCodesByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]string, error)
}

type LoginSessionStore interface {
	CreateLoginSession(ctx context.Context, session *entity.LoginSession, ttl time.Duration) error
	GetLoginSession(ctx context.Context, sessionID uuid.UUID) (*entity.LoginSession, error)
	UpdateLoginSession(ctx context.Context, session *entity.LoginSession, ttl time.Duration) error
	DeleteLoginSession(ctx context.Context, sessionID uuid.UUID) error

	IncrementLoginAttempts(ctx context.Context, sessionID uuid.UUID) (int, error)
	UpdateLoginOTP(ctx context.Context, sessionID uuid.UUID, otpHash string, expiresAt time.Time) error
	MarkLoginVerified(ctx context.Context, sessionID uuid.UUID) error

	IncrementLoginRateLimit(ctx context.Context, email string, ttl time.Duration) (int64, error)
	GetLoginRateLimitCount(ctx context.Context, email string) (int64, error)
}

type UserSessionRepository interface {
	Create(ctx context.Context, session *entity.UserSession) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.UserSession, error)
	GetByRefreshTokenID(ctx context.Context, refreshTokenID uuid.UUID) (*entity.UserSession, error)
	UpdateLastActive(ctx context.Context, id uuid.UUID) error
	UpdateRefreshTokenID(ctx context.Context, sessionID uuid.UUID, refreshTokenID uuid.UUID) error
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
}

type UserTenantRegistrationRepository interface {
	ListActiveByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserTenantRegistration, error)
	ListByUserIDForClaims(ctx context.Context, userID uuid.UUID) ([]entity.UserTenantRegistration, error)
}

type ProductsByTenantRepository interface {
	ListActiveByTenantID(ctx context.Context, tenantID uuid.UUID) ([]entity.Product, error)
}

type RegistrationSessionStore interface {
	CreateRegistrationSession(ctx context.Context, session *entity.RegistrationSession, ttl time.Duration) error
	GetRegistrationSession(ctx context.Context, sessionID uuid.UUID) (*entity.RegistrationSession, error)
	UpdateRegistrationSession(ctx context.Context, session *entity.RegistrationSession, ttl time.Duration) error
	DeleteRegistrationSession(ctx context.Context, sessionID uuid.UUID) error

	IncrementRegistrationAttempts(ctx context.Context, sessionID uuid.UUID) (int, error)
	UpdateRegistrationOTP(ctx context.Context, sessionID uuid.UUID, otpHash string, expiresAt time.Time) error
	MarkRegistrationVerified(ctx context.Context, sessionID uuid.UUID, tokenHash string) error
	MarkRegistrationPasswordSet(ctx context.Context, sessionID uuid.UUID, passwordHash string, tokenHash string) error
	GetRegistrationPasswordHash(ctx context.Context, sessionID uuid.UUID) (string, error)

	LockRegistrationEmail(ctx context.Context, email string, ttl time.Duration) (bool, error)
	UnlockRegistrationEmail(ctx context.Context, email string) error
	IsRegistrationEmailLocked(ctx context.Context, email string) (bool, error)

	IncrementRegistrationRateLimit(ctx context.Context, email string, ttl time.Duration) (int64, error)
	GetRegistrationRateLimitCount(ctx context.Context, email string) (int64, error)
}

type TokenBlacklistStore interface {
	BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error
	IsTokenBlacklisted(ctx context.Context, jti string) (bool, error)
	BlacklistUser(ctx context.Context, userID uuid.UUID, timestamp time.Time, ttl time.Duration) error
	GetUserBlacklistTimestamp(ctx context.Context, userID uuid.UUID) (*time.Time, error)
}

type OAuthStateStore interface {
	StoreOAuthState(ctx context.Context, state string, ttl time.Duration) error
	GetAndDeleteOAuthState(ctx context.Context, state string) (bool, error)
}

type InMemoryStore interface {
	RegistrationSessionStore
	LoginSessionStore
	TokenBlacklistStore
	OAuthStateStore
}
