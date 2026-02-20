package internal

import (
	"context"
	"time"

	"iam-service/entity"
	usercontract "iam-service/iam/user/contract"
	"iam-service/masterdata/masterdatadto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockMasterdataUsecase struct {
	mock.Mock
}

func (m *MockMasterdataUsecase) ValidateItemCode(ctx context.Context, req *masterdatadto.ValidateCodeRequest) (*masterdatadto.ValidateCodeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*masterdatadto.ValidateCodeResponse), args.Error(1)
}

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) == nil {
		return fn(ctx)
	}
	return args.Error(0)
}

func NewMockTransactionManager() *MockTransactionManager {
	m := &MockTransactionManager{}
	m.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
	return m
}

type MockTenantRepository struct {
	mock.Mock
}

func (m *MockTenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tenant), args.Error(1)
}

func (m *MockTenantRepository) GetByCode(ctx context.Context, code string) (*entity.Tenant, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tenant), args.Error(1)
}

func (m *MockTenantRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, filter *usercontract.UserListFilter) ([]*entity.User, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entity.User), args.Get(1).(int64), args.Error(2)
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendRegistrationOTP(ctx context.Context, email, otp string, expiryMinutes int) error {
	args := m.Called(ctx, email, otp, expiryMinutes)
	return args.Error(0)
}

func (m *MockEmailService) SendLoginOTP(ctx context.Context, email, otp string, expiryMinutes int) error {
	args := m.Called(ctx, email, otp, expiryMinutes)
	return args.Error(0)
}

func (m *MockEmailService) SendWelcome(ctx context.Context, email, firstName string) error {
	args := m.Called(ctx, email, firstName)
	return args.Error(0)
}

func (m *MockEmailService) SendPasswordReset(ctx context.Context, email, token string, expiryMinutes int) error {
	args := m.Called(ctx, email, token, expiryMinutes)
	return args.Error(0)
}

func (m *MockEmailService) SendPINReset(ctx context.Context, email, otp string, expiryMinutes int) error {
	args := m.Called(ctx, email, otp, expiryMinutes)
	return args.Error(0)
}

func (m *MockEmailService) SendAdminInvitation(ctx context.Context, email, token string, expiryMinutes int) error {
	args := m.Called(ctx, email, token, expiryMinutes)
	return args.Error(0)
}

type MockUserProfileRepository struct {
	mock.Mock
}

func (m *MockUserProfileRepository) Create(ctx context.Context, profile *entity.UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockUserProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserProfile), args.Error(1)
}

func (m *MockUserProfileRepository) Update(ctx context.Context, profile *entity.UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

type MockUserAuthMethodRepository struct {
	mock.Mock
}

func (m *MockUserAuthMethodRepository) Create(ctx context.Context, authMethod *entity.UserAuthMethod) error {
	args := m.Called(ctx, authMethod)

	if authMethod.ID == uuid.Nil {
		authMethod.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockUserAuthMethodRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserAuthMethod, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserAuthMethod), args.Error(1)
}

func (m *MockUserAuthMethodRepository) Update(ctx context.Context, authMethod *entity.UserAuthMethod) error {
	args := m.Called(ctx, authMethod)
	return args.Error(0)
}

type MockUserSecurityStateRepository struct {
	mock.Mock
}

func (m *MockUserSecurityStateRepository) Create(ctx context.Context, securityState *entity.UserSecurityState) error {
	args := m.Called(ctx, securityState)
	return args.Error(0)
}

func (m *MockUserSecurityStateRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserSecurityState, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSecurityState), args.Error(1)
}

func (m *MockUserSecurityStateRepository) Update(ctx context.Context, securityState *entity.UserSecurityState) error {
	args := m.Called(ctx, securityState)
	return args.Error(0)
}

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Create(ctx context.Context, role *entity.Role) error {
	args := m.Called(ctx, role)

	if role.ID == uuid.Nil {
		role.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockRoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByName(ctx context.Context, productID uuid.UUID, name string) (*entity.Role, error) {
	args := m.Called(ctx, productID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByCode(ctx context.Context, productID uuid.UUID, code string) (*entity.Role, error) {
	args := m.Called(ctx, productID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.Role, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Role), args.Error(1)
}

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.RefreshToken, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) SetReplacedBy(ctx context.Context, id uuid.UUID, replacedByID uuid.UUID) error {
	args := m.Called(ctx, id, replacedByID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	args := m.Called(ctx, token)

	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	if token.TokenFamily == uuid.Nil {
		token.TokenFamily = uuid.New()
	}
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID, reason string) error {
	args := m.Called(ctx, id, reason)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID, reason string) error {
	args := m.Called(ctx, userID, reason)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeByFamily(ctx context.Context, tokenFamily uuid.UUID, reason string) error {
	args := m.Called(ctx, tokenFamily, reason)
	return args.Error(0)
}

type MockUserRoleRepository struct {
	mock.Mock
}

func (m *MockUserRoleRepository) Create(ctx context.Context, userRole *entity.UserRole) error {
	args := m.Called(ctx, userRole)

	if userRole.ID == uuid.Nil {
		userRole.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockUserRoleRepository) ListActiveByUserID(ctx context.Context, userID uuid.UUID, productID *uuid.UUID) ([]entity.UserRole, error) {
	args := m.Called(ctx, userID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserRole), args.Error(1)
}

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetByCodeAndTenant(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Product, error) {
	args := m.Called(ctx, tenantID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductRepository) GetByIDAndTenant(ctx context.Context, productID, tenantID uuid.UUID) (*entity.Product, error) {
	args := m.Called(ctx, productID, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) GetCodesByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]string, error) {
	args := m.Called(ctx, roleIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

type MockUserSessionRepository struct {
	mock.Mock
}

func (m *MockUserSessionRepository) Create(ctx context.Context, session *entity.UserSession) error {
	args := m.Called(ctx, session)
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockUserSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.UserSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSession), args.Error(1)
}

func (m *MockUserSessionRepository) GetByRefreshTokenID(ctx context.Context, refreshTokenID uuid.UUID) (*entity.UserSession, error) {
	args := m.Called(ctx, refreshTokenID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSession), args.Error(1)
}

func (m *MockUserSessionRepository) UpdateLastActive(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserSessionRepository) UpdateRefreshTokenID(ctx context.Context, sessionID uuid.UUID, refreshTokenID uuid.UUID) error {
	args := m.Called(ctx, sessionID, refreshTokenID)
	return args.Error(0)
}

func (m *MockUserSessionRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserSessionRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type MockInMemoryStore struct {
	mock.Mock
}

func (m *MockInMemoryStore) CreateRegistrationSession(ctx context.Context, session *entity.RegistrationSession, ttl time.Duration) error {
	args := m.Called(ctx, session, ttl)
	return args.Error(0)
}

func (m *MockInMemoryStore) GetRegistrationSession(ctx context.Context, sessionID uuid.UUID) (*entity.RegistrationSession, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.RegistrationSession), args.Error(1)
}

func (m *MockInMemoryStore) UpdateRegistrationSession(ctx context.Context, session *entity.RegistrationSession, ttl time.Duration) error {
	args := m.Called(ctx, session, ttl)
	return args.Error(0)
}

func (m *MockInMemoryStore) DeleteRegistrationSession(ctx context.Context, sessionID uuid.UUID) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockInMemoryStore) IncrementRegistrationAttempts(ctx context.Context, sessionID uuid.UUID) (int, error) {
	args := m.Called(ctx, sessionID)
	return args.Int(0), args.Error(1)
}

func (m *MockInMemoryStore) UpdateRegistrationOTP(ctx context.Context, sessionID uuid.UUID, otpHash string, expiresAt time.Time) error {
	args := m.Called(ctx, sessionID, otpHash, expiresAt)
	return args.Error(0)
}

func (m *MockInMemoryStore) MarkRegistrationVerified(ctx context.Context, sessionID uuid.UUID, tokenHash string) error {
	args := m.Called(ctx, sessionID, tokenHash)
	return args.Error(0)
}

func (m *MockInMemoryStore) MarkRegistrationPasswordSet(ctx context.Context, sessionID uuid.UUID, passwordHash string, tokenHash string) error {
	args := m.Called(ctx, sessionID, passwordHash, tokenHash)
	return args.Error(0)
}

func (m *MockInMemoryStore) GetRegistrationPasswordHash(ctx context.Context, sessionID uuid.UUID) (string, error) {
	args := m.Called(ctx, sessionID)
	return args.String(0), args.Error(1)
}

func (m *MockInMemoryStore) LockRegistrationEmail(ctx context.Context, email string, ttl time.Duration) (bool, error) {
	args := m.Called(ctx, email, ttl)
	return args.Bool(0), args.Error(1)
}

func (m *MockInMemoryStore) UnlockRegistrationEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockInMemoryStore) IsRegistrationEmailLocked(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockInMemoryStore) IncrementRegistrationRateLimit(ctx context.Context, email string, ttl time.Duration) (int64, error) {
	args := m.Called(ctx, email, ttl)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockInMemoryStore) GetRegistrationRateLimitCount(ctx context.Context, email string) (int64, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockInMemoryStore) CreateLoginSession(ctx context.Context, session *entity.LoginSession, ttl time.Duration) error {
	args := m.Called(ctx, session, ttl)
	return args.Error(0)
}

func (m *MockInMemoryStore) GetLoginSession(ctx context.Context, sessionID uuid.UUID) (*entity.LoginSession, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.LoginSession), args.Error(1)
}

func (m *MockInMemoryStore) UpdateLoginSession(ctx context.Context, session *entity.LoginSession, ttl time.Duration) error {
	args := m.Called(ctx, session, ttl)
	return args.Error(0)
}

func (m *MockInMemoryStore) DeleteLoginSession(ctx context.Context, sessionID uuid.UUID) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockInMemoryStore) IncrementLoginAttempts(ctx context.Context, sessionID uuid.UUID) (int, error) {
	args := m.Called(ctx, sessionID)
	return args.Int(0), args.Error(1)
}

func (m *MockInMemoryStore) UpdateLoginOTP(ctx context.Context, sessionID uuid.UUID, otpHash string, expiresAt time.Time) error {
	args := m.Called(ctx, sessionID, otpHash, expiresAt)
	return args.Error(0)
}

func (m *MockInMemoryStore) MarkLoginVerified(ctx context.Context, sessionID uuid.UUID) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockInMemoryStore) IncrementLoginRateLimit(ctx context.Context, email string, ttl time.Duration) (int64, error) {
	args := m.Called(ctx, email, ttl)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockInMemoryStore) GetLoginRateLimitCount(ctx context.Context, email string) (int64, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockInMemoryStore) BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error {
	args := m.Called(ctx, jti, ttl)
	return args.Error(0)
}

func (m *MockInMemoryStore) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	args := m.Called(ctx, jti)
	return args.Bool(0), args.Error(1)
}

func (m *MockInMemoryStore) BlacklistUser(ctx context.Context, userID uuid.UUID, timestamp time.Time, ttl time.Duration) error {
	args := m.Called(ctx, userID, timestamp, ttl)
	return args.Error(0)
}

func (m *MockInMemoryStore) GetUserBlacklistTimestamp(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*time.Time), args.Error(1)
}

type MockUserTenantRegistrationRepository struct {
	mock.Mock
}

func (m *MockUserTenantRegistrationRepository) ListActiveByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserTenantRegistration, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserTenantRegistration), args.Error(1)
}

type MockProductsByTenantRepository struct {
	mock.Mock
}

func (m *MockProductsByTenantRepository) ListActiveByTenantID(ctx context.Context, tenantID uuid.UUID) ([]entity.Product, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Product), args.Error(1)
}
