package auth_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"erp-service/config"
	"erp-service/entity"
	"erp-service/iam/auth"
	"erp-service/masterdata"
	pkgerrors "erp-service/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCompleteGoogleProfile(t *testing.T) {
	registrationID := uuid.New()
	email := "google.user@example.com"
	jwtSecret := "test-secret-key-for-testing-purposes"

	generateValidToken := func() (string, string) {
		claims := jwt.MapClaims{
			"registration_id": registrationID.String(),
			"email":           email,
			"purpose":         auth.GoogleRegistrationCompleteTokenPurpose,
			"exp":             time.Now().Add(15 * time.Minute).Unix(),
			"iat":             time.Now().Unix(),
			"jti":             uuid.New().String(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(jwtSecret))
		hash := sha256.Sum256([]byte(tokenString))
		tokenHash := hex.EncodeToString(hash[:])
		return tokenString, tokenHash
	}

	validReq := &auth.CompleteGoogleProfileRequest{
		FullName:    "Ahmad Fauzi",
		DateOfBirth: "1990-05-20",
		Gender:      "GENDER_001",
	}

	type setupMockFn func(
		txManager *MockTransactionManager,
		redis *MockInMemoryStore,
		userRepo *MockUserRepository,
		profileRepo *MockUserProfileRepository,
		authMethodRepo *MockUserAuthMethodRepository,
		securityStateRepo *MockUserSecurityStateRepository,
		emailSvc *MockEmailService,
		refreshTokenRepo *MockRefreshTokenRepository,
		masterdataUsecase *MockMasterdataUsecase,
		tokenHash string,
	)

	tests := []struct {
		name          string
		req           *auth.CompleteGoogleProfileRequest
		setupToken    func() string
		setupMocks    setupMockFn
		expectedError string
		expectedCode  string
		validateResp  func(*testing.T, *auth.CompleteGoogleProfileResponse)
	}{
		{
			name: "success - happy path",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(
				txManager *MockTransactionManager,
				redis *MockInMemoryStore,
				userRepo *MockUserRepository,
				profileRepo *MockUserProfileRepository,
				authMethodRepo *MockUserAuthMethodRepository,
				securityStateRepo *MockUserSecurityStateRepository,
				emailSvc *MockEmailService,
				refreshTokenRepo *MockRefreshTokenRepository,
				mdUsecase *MockMasterdataUsecase,
				tokenHash string,
			) {
				session := &entity.GoogleRegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					GoogleID:              "google-id-123",
					Name:                  "Ahmad Fauzi",
					Picture:               "https://picture.url/photo.jpg",
					GoogleEmailVerified:   true,
					Status:                entity.GoogleRegistrationSessionStatusPendingProfile,
					RegistrationTokenHash: tokenHash,
					CreatedAt:             time.Now().Add(-1 * time.Minute),
					ExpiresAt:             time.Now().Add(14 * time.Minute),
				}
				redis.On("GetAndDeleteGoogleRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				mdUsecase.On("ValidateItemCode", mock.Anything, &masterdata.ValidateCodeRequest{
					CategoryCode: "GENDER",
					ItemCode:     "GENDER_001",
				}).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(false, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
				profileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserProfile")).Return(nil)
				authMethodRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserAuthMethod")).Return(nil)
				securityStateRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserSecurityState")).Return(nil)
				refreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)
				// No DeleteGoogleRegistrationSession: deletion is now done atomically
				// by GetAndDeleteGoogleRegistrationSession at the start of the flow.
				emailSvc.On("SendWelcome", mock.Anything, email, "Ahmad").Return(nil)
			},
			validateResp: func(t *testing.T, resp *auth.CompleteGoogleProfileResponse) {
				assert.Equal(t, email, resp.Email)
				assert.NotEmpty(t, resp.AccessToken)
				assert.NotEmpty(t, resp.RefreshToken)
				assert.Equal(t, "Bearer", resp.TokenType)
				assert.Greater(t, resp.ExpiresIn, 0)
				assert.Equal(t, "Ahmad", resp.Profile.FirstName)
				assert.Equal(t, "Fauzi", resp.Profile.LastName)
				assert.Equal(t, string(entity.UserStatusActive), resp.Status)
			},
		},
		{
			name: "error - invalid registration token",
			req:  validReq,
			setupToken: func() string {
				return "invalid-token"
			},
			setupMocks: func(
				txManager *MockTransactionManager,
				redis *MockInMemoryStore,
				userRepo *MockUserRepository,
				profileRepo *MockUserProfileRepository,
				authMethodRepo *MockUserAuthMethodRepository,
				securityStateRepo *MockUserSecurityStateRepository,
				emailSvc *MockEmailService,
				refreshTokenRepo *MockRefreshTokenRepository,
				mdUsecase *MockMasterdataUsecase,
				tokenHash string,
			) {
			},
			expectedError: "invalid",
			expectedCode:  pkgerrors.CodeUnauthorized,
		},
		{
			name: "error - expired google registration session (410 Gone)",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(
				txManager *MockTransactionManager,
				redis *MockInMemoryStore,
				userRepo *MockUserRepository,
				profileRepo *MockUserProfileRepository,
				authMethodRepo *MockUserAuthMethodRepository,
				securityStateRepo *MockUserSecurityStateRepository,
				emailSvc *MockEmailService,
				refreshTokenRepo *MockRefreshTokenRepository,
				mdUsecase *MockMasterdataUsecase,
				tokenHash string,
			) {
				session := &entity.GoogleRegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					GoogleID:              "google-id-123",
					Status:                entity.GoogleRegistrationSessionStatusPendingProfile,
					RegistrationTokenHash: tokenHash,
					CreatedAt:             time.Now().Add(-20 * time.Minute),
					ExpiresAt:             time.Now().Add(-5 * time.Minute), // already expired
				}
				redis.On("GetAndDeleteGoogleRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "expired",
			expectedCode:  "GOOGLE_REGISTRATION_EXPIRED",
		},
		{
			name: "error - session not in PENDING_PROFILE state (already completed)",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(
				txManager *MockTransactionManager,
				redis *MockInMemoryStore,
				userRepo *MockUserRepository,
				profileRepo *MockUserProfileRepository,
				authMethodRepo *MockUserAuthMethodRepository,
				securityStateRepo *MockUserSecurityStateRepository,
				emailSvc *MockEmailService,
				refreshTokenRepo *MockRefreshTokenRepository,
				mdUsecase *MockMasterdataUsecase,
				tokenHash string,
			) {
				session := &entity.GoogleRegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					GoogleID:              "google-id-123",
					Status:                entity.GoogleRegistrationSessionStatusCompleted, // already completed
					RegistrationTokenHash: tokenHash,
					CreatedAt:             time.Now().Add(-1 * time.Minute),
					ExpiresAt:             time.Now().Add(14 * time.Minute),
				}
				redis.On("GetAndDeleteGoogleRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "not in the correct state",
			expectedCode:  pkgerrors.CodeConflict,
		},
		{
			name: "error - token hash mismatch",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(
				txManager *MockTransactionManager,
				redis *MockInMemoryStore,
				userRepo *MockUserRepository,
				profileRepo *MockUserProfileRepository,
				authMethodRepo *MockUserAuthMethodRepository,
				securityStateRepo *MockUserSecurityStateRepository,
				emailSvc *MockEmailService,
				refreshTokenRepo *MockRefreshTokenRepository,
				mdUsecase *MockMasterdataUsecase,
				tokenHash string,
			) {
				session := &entity.GoogleRegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					GoogleID:              "google-id-123",
					Status:                entity.GoogleRegistrationSessionStatusPendingProfile,
					RegistrationTokenHash: "different-hash-entirely",
					CreatedAt:             time.Now().Add(-1 * time.Minute),
					ExpiresAt:             time.Now().Add(14 * time.Minute),
				}
				redis.On("GetAndDeleteGoogleRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "token",
			expectedCode:  pkgerrors.CodeUnauthorized,
		},
		{
			name: "error - invalid gender format",
			req: &auth.CompleteGoogleProfileRequest{
				FullName:    "Ahmad Fauzi",
				DateOfBirth: "1990-05-20",
				Gender:      "male",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(
				txManager *MockTransactionManager,
				redis *MockInMemoryStore,
				userRepo *MockUserRepository,
				profileRepo *MockUserProfileRepository,
				authMethodRepo *MockUserAuthMethodRepository,
				securityStateRepo *MockUserSecurityStateRepository,
				emailSvc *MockEmailService,
				refreshTokenRepo *MockRefreshTokenRepository,
				mdUsecase *MockMasterdataUsecase,
				tokenHash string,
			) {
				session := &entity.GoogleRegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					GoogleID:              "google-id-123",
					Status:                entity.GoogleRegistrationSessionStatusPendingProfile,
					RegistrationTokenHash: tokenHash,
					CreatedAt:             time.Now().Add(-1 * time.Minute),
					ExpiresAt:             time.Now().Add(14 * time.Minute),
				}
				redis.On("GetAndDeleteGoogleRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "GENDER_NNN",
			expectedCode:  pkgerrors.CodeValidation,
		},
		{
			name: "error - age under 18",
			req: &auth.CompleteGoogleProfileRequest{
				FullName:    "Young Person",
				DateOfBirth: time.Now().AddDate(-17, 0, 0).Format("2006-01-02"),
				Gender:      "GENDER_001",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(
				txManager *MockTransactionManager,
				redis *MockInMemoryStore,
				userRepo *MockUserRepository,
				profileRepo *MockUserProfileRepository,
				authMethodRepo *MockUserAuthMethodRepository,
				securityStateRepo *MockUserSecurityStateRepository,
				emailSvc *MockEmailService,
				refreshTokenRepo *MockRefreshTokenRepository,
				mdUsecase *MockMasterdataUsecase,
				tokenHash string,
			) {
				session := &entity.GoogleRegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					GoogleID:              "google-id-123",
					Status:                entity.GoogleRegistrationSessionStatusPendingProfile,
					RegistrationTokenHash: tokenHash,
					CreatedAt:             time.Now().Add(-1 * time.Minute),
					ExpiresAt:             time.Now().Add(14 * time.Minute),
				}
				redis.On("GetAndDeleteGoogleRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				mdUsecase.On("ValidateItemCode", mock.Anything, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
			},
			expectedError: "18 years",
			expectedCode:  pkgerrors.CodeValidation,
		},
		{
			name: "error - email already taken (race condition guard)",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(
				txManager *MockTransactionManager,
				redis *MockInMemoryStore,
				userRepo *MockUserRepository,
				profileRepo *MockUserProfileRepository,
				authMethodRepo *MockUserAuthMethodRepository,
				securityStateRepo *MockUserSecurityStateRepository,
				emailSvc *MockEmailService,
				refreshTokenRepo *MockRefreshTokenRepository,
				mdUsecase *MockMasterdataUsecase,
				tokenHash string,
			) {
				session := &entity.GoogleRegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					GoogleID:              "google-id-123",
					Status:                entity.GoogleRegistrationSessionStatusPendingProfile,
					RegistrationTokenHash: tokenHash,
					CreatedAt:             time.Now().Add(-1 * time.Minute),
					ExpiresAt:             time.Now().Add(14 * time.Minute),
				}
				redis.On("GetAndDeleteGoogleRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				mdUsecase.On("ValidateItemCode", mock.Anything, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(true, nil)
			},
			expectedError: "already been registered",
			expectedCode:  pkgerrors.CodeConflict,
		},
		{
			name: "success - single name splits correctly",
			req: &auth.CompleteGoogleProfileRequest{
				FullName:    "Soekarno",
				DateOfBirth: "1990-05-20",
				Gender:      "GENDER_001",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(
				txManager *MockTransactionManager,
				redis *MockInMemoryStore,
				userRepo *MockUserRepository,
				profileRepo *MockUserProfileRepository,
				authMethodRepo *MockUserAuthMethodRepository,
				securityStateRepo *MockUserSecurityStateRepository,
				emailSvc *MockEmailService,
				refreshTokenRepo *MockRefreshTokenRepository,
				mdUsecase *MockMasterdataUsecase,
				tokenHash string,
			) {
				session := &entity.GoogleRegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					GoogleID:              "google-id-123",
					Status:                entity.GoogleRegistrationSessionStatusPendingProfile,
					RegistrationTokenHash: tokenHash,
					CreatedAt:             time.Now().Add(-1 * time.Minute),
					ExpiresAt:             time.Now().Add(14 * time.Minute),
				}
				redis.On("GetAndDeleteGoogleRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				mdUsecase.On("ValidateItemCode", mock.Anything, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(false, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
				profileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserProfile")).Return(nil)
				authMethodRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserAuthMethod")).Return(nil)
				securityStateRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserSecurityState")).Return(nil)
				refreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)
				// No DeleteGoogleRegistrationSession: deletion is now done atomically
				// by GetAndDeleteGoogleRegistrationSession at the start of the flow.
				emailSvc.On("SendWelcome", mock.Anything, email, "Soekarno").Return(nil)
			},
			validateResp: func(t *testing.T, resp *auth.CompleteGoogleProfileResponse) {
				assert.Equal(t, "Soekarno", resp.Profile.FirstName)
				assert.Equal(t, "", resp.Profile.LastName)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txManager := NewMockTransactionManager()
			redis := &MockInMemoryStore{}
			userRepo := &MockUserRepository{}
			profileRepo := &MockUserProfileRepository{}
			authMethodRepo := &MockUserAuthMethodRepository{}
			securityStateRepo := &MockUserSecurityStateRepository{}
			emailSvc := &MockEmailService{}
			refreshTokenRepo := &MockRefreshTokenRepository{}
			mdUsecase := &MockMasterdataUsecase{}

			tokenString := tt.setupToken()
			_, tokenHash := generateValidToken()
			if tokenString != "invalid-token" {
				hash := sha256.Sum256([]byte(tokenString))
				tokenHash = hex.EncodeToString(hash[:])
			}
			tt.setupMocks(txManager, redis, userRepo, profileRepo, authMethodRepo, securityStateRepo, emailSvc, refreshTokenRepo, mdUsecase, tokenHash)

			cfg := &config.Config{
				JWT: config.JWTConfig{
					AccessSecret:  jwtSecret,
					RefreshSecret: "refresh-secret",
					SigningMethod: "HS256",
					AccessExpiry:  3600 * time.Second,
					RefreshExpiry: 86400 * time.Second,
					Issuer:        "erp-service",
					Audience:      []string{"erp-api"},
				},
			}

			uc := auth.NewUsecase(
				txManager, cfg,
				userRepo, profileRepo, authMethodRepo, securityStateRepo,
				nil, nil, refreshTokenRepo, nil, nil, nil,
				emailSvc, redis,
				nil, nil, nil, nil,
				mdUsecase,
			)

			req := &auth.CompleteGoogleProfileRequest{
				RegistrationID:    registrationID,
				RegistrationToken: tokenString,
				IPAddress:         "127.0.0.1",
				UserAgent:         "test-agent",
				FullName:          tt.req.FullName,
				DateOfBirth:       tt.req.DateOfBirth,
				Gender:            tt.req.Gender,
			}

			resp, err := uc.CompleteGoogleProfile(context.Background(), req)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				appErr, ok := err.(*pkgerrors.AppError)
				require.True(t, ok, "Error should be AppError, got: %T: %v", err, err)
				assert.Equal(t, tt.expectedCode, appErr.Code)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			if tt.validateResp != nil {
				tt.validateResp(t, resp)
			}
		})
	}
}
