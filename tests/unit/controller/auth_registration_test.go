package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"erp-service/config"
	"erp-service/delivery/http/controller"
	"erp-service/iam/auth"
	"erp-service/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) Logout(ctx context.Context, req *auth.LogoutRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAuthUsecase) LogoutAll(ctx context.Context, req *auth.LogoutAllRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAuthUsecase) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.RefreshTokenResponse), args.Error(1)
}

func (m *MockAuthUsecase) InitiateRegistration(ctx context.Context, req *auth.InitiateRegistrationRequest) (*auth.InitiateRegistrationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.InitiateRegistrationResponse), args.Error(1)
}

func (m *MockAuthUsecase) VerifyRegistrationOTP(ctx context.Context, req *auth.VerifyRegistrationOTPRequest) (*auth.VerifyRegistrationOTPResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.VerifyRegistrationOTPResponse), args.Error(1)
}

func (m *MockAuthUsecase) ResendRegistrationOTP(ctx context.Context, req *auth.ResendRegistrationOTPRequest) (*auth.ResendRegistrationOTPResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.ResendRegistrationOTPResponse), args.Error(1)
}

func (m *MockAuthUsecase) SetPassword(ctx context.Context, req *auth.SetPasswordRequest) (*auth.SetPasswordResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.SetPasswordResponse), args.Error(1)
}

func (m *MockAuthUsecase) CompleteProfileRegistration(ctx context.Context, req *auth.CompleteProfileRegistrationRequest) (*auth.CompleteProfileRegistrationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.CompleteProfileRegistrationResponse), args.Error(1)
}

func (m *MockAuthUsecase) GetRegistrationStatus(ctx context.Context, registrationID uuid.UUID, email string) (*auth.RegistrationStatusResponse, error) {
	args := m.Called(ctx, registrationID, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.RegistrationStatusResponse), args.Error(1)
}

func (m *MockAuthUsecase) InitiateLogin(ctx context.Context, req *auth.InitiateLoginRequest) (*auth.UnifiedLoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.UnifiedLoginResponse), args.Error(1)
}

func (m *MockAuthUsecase) VerifyLoginOTP(ctx context.Context, req *auth.VerifyLoginOTPRequest) (*auth.VerifyLoginOTPResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.VerifyLoginOTPResponse), args.Error(1)
}

func (m *MockAuthUsecase) ResendLoginOTP(ctx context.Context, req *auth.ResendLoginOTPRequest) (*auth.ResendLoginOTPResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.ResendLoginOTPResponse), args.Error(1)
}

func (m *MockAuthUsecase) GetLoginStatus(ctx context.Context, req *auth.GetLoginStatusRequest) (*auth.LoginStatusResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.LoginStatusResponse), args.Error(1)
}

func (m *MockAuthUsecase) GetGoogleAuthURL(ctx context.Context) (*auth.GoogleAuthURLResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.GoogleAuthURLResponse), args.Error(1)
}

func (m *MockAuthUsecase) HandleGoogleCallback(ctx context.Context, req *auth.GoogleCallbackRequest) (*auth.GoogleCallbackResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.GoogleCallbackResponse), args.Error(1)
}

func (m *MockAuthUsecase) CreateTransferToken(ctx context.Context, req *auth.CreateTransferTokenRequest) (*auth.CreateTransferTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.CreateTransferTokenResponse), args.Error(1)
}

func (m *MockAuthUsecase) ExchangeTransferToken(ctx context.Context, req *auth.ExchangeTransferTokenRequest) (*auth.ExchangeTransferTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.ExchangeTransferTokenResponse), args.Error(1)
}

func (m *MockAuthUsecase) LogoutTree(ctx context.Context, req *auth.LogoutTreeRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func setupTestApp() *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var appErr *errors.AppError
			if errors.As(err, &appErr) {
				return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
					"success": false,
					"error":   appErr.Code,
					"message": appErr.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "INTERNAL",
				"message": err.Error(),
			})
		},
	})
}

func TestSetPasswordController(t *testing.T) {
	registrationID := uuid.New()
	registrationToken := "valid-registration-token"

	tests := []struct {
		name           string
		body           map[string]any
		registrationID string
		authHeader     string
		setupMock      func(*MockAuthUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]any)
	}{
		{
			name:           "success - password set",
			registrationID: registrationID.String(),
			authHeader:     "Bearer " + registrationToken,
			body: map[string]any{
				"password":              "SecureP@ssw0rd!",
				"confirmation_password": "SecureP@ssw0rd!",
			},
			setupMock: func(m *MockAuthUsecase) {
				m.On("SetPassword", mock.Anything, mock.AnythingOfType("*auth.SetPasswordRequest")).
					Return(&auth.SetPasswordResponse{
						RegistrationID:    registrationID.String(),
						Status:            "PASSWORD_SET",
						Message:           "Password set successfully. Please proceed to complete your profile.",
						RegistrationToken: "new-token",
						NextStep: auth.NextStep{
							Action:         "set-profile",
							Endpoint:       "/api/v1/auth/registration/complete-profile",
							RequiredFields: []string{"full_name", "gender", "date_of_birth"},
						},
					}, nil)
			},
			expectedStatus: fiber.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]any) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]any)
				assert.Equal(t, "PASSWORD_SET", data["status"])
				assert.Equal(t, registrationID.String(), data["registration_id"])
				nextStep := data["next_step"].(map[string]any)
				assert.Equal(t, "set-profile", nextStep["action"])

				requiredFields := nextStep["required_fields"].([]any)
				assert.Equal(t, []any{"full_name", "gender", "date_of_birth"}, requiredFields)
			},
		},
		{
			name:           "error - missing Authorization header",
			registrationID: registrationID.String(),
			authHeader:     "",
			body: map[string]any{
				"password":              "SecureP@ssw0rd!",
				"confirmation_password": "SecureP@ssw0rd!",
			},
			setupMock:      func(m *MockAuthUsecase) {},
			expectedStatus: fiber.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]any) {
				assert.False(t, resp["success"].(bool))
			},
		},
		{
			name:           "error - invalid registration ID",
			registrationID: "not-a-uuid",
			authHeader:     "Bearer " + registrationToken,
			body: map[string]any{
				"password":              "SecureP@ssw0rd!",
				"confirmation_password": "SecureP@ssw0rd!",
			},
			setupMock:      func(m *MockAuthUsecase) {},
			expectedStatus: fiber.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]any) {
				assert.False(t, resp["success"].(bool))
			},
		},
		{
			name:           "error - usecase returns forbidden",
			registrationID: registrationID.String(),
			authHeader:     "Bearer " + registrationToken,
			body: map[string]any{
				"password":              "SecureP@ssw0rd!",
				"confirmation_password": "SecureP@ssw0rd!",
			},
			setupMock: func(m *MockAuthUsecase) {
				m.On("SetPassword", mock.Anything, mock.AnythingOfType("*auth.SetPasswordRequest")).
					Return(nil, errors.ErrForbidden("Email has not been verified"))
			},
			expectedStatus: fiber.StatusForbidden,
			checkResponse: func(t *testing.T, resp map[string]any) {
				assert.False(t, resp["success"].(bool))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := new(MockAuthUsecase)
			tt.setupMock(mockUC)

			app := setupTestApp()
			ctrl := controller.NewRegistrationController(&config.Config{}, mockUC)
			app.Post("/registrations/:id/set-password", ctrl.SetPassword)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/registrations/"+tt.registrationID+"/set-password", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var respBody map[string]any
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
			tt.checkResponse(t, respBody)

			mockUC.AssertExpectations(t)
		})
	}
}

func TestCompleteProfileRegistrationController(t *testing.T) {
	registrationID := uuid.New()
	userID := uuid.New()
	registrationToken := "valid-registration-token"

	tests := []struct {
		name           string
		body           map[string]any
		registrationID string
		authHeader     string
		setupMock      func(*MockAuthUsecase)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]any)
	}{
		{
			name:           "success - profile completed and user created",
			registrationID: registrationID.String(),
			authHeader:     "Bearer " + registrationToken,
			body: map[string]any{
				"full_name":     "John Michael Smith",
				"date_of_birth": "1990-01-15",
				"gender":        "GENDER_001",
			},
			setupMock: func(m *MockAuthUsecase) {
				m.On("CompleteProfileRegistration", mock.Anything, mock.AnythingOfType("*auth.CompleteProfileRegistrationRequest")).
					Return(&auth.CompleteProfileRegistrationResponse{
						UserID:  userID,
						Email:   "john@example.com",
						Status:  "active",
						Message: "Registration completed successfully. You are now logged in.",
						Profile: auth.RegistrationUserProfile{
							FirstName: "John Michael",
							LastName:  "Smith",
						},
						AccessToken:  "access-token-value",
						RefreshToken: "refresh-token-value",
						TokenType:    "Bearer",
						ExpiresIn:    3600,
					}, nil)
			},
			expectedStatus: fiber.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]any) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]any)
				assert.Equal(t, "active", data["status"])
				assert.Equal(t, "john@example.com", data["email"])
				assert.Equal(t, "access-token-value", data["access_token"])
				assert.Equal(t, "refresh-token-value", data["refresh_token"])
				assert.Equal(t, "Bearer", data["token_type"])
				profile := data["profile"].(map[string]any)
				assert.Equal(t, "John Michael", profile["first_name"])
				assert.Equal(t, "Smith", profile["last_name"])
			},
		},
		{
			name:           "error - missing Authorization header",
			registrationID: registrationID.String(),
			authHeader:     "",
			body: map[string]any{
				"full_name": "John Smith",
			},
			setupMock:      func(m *MockAuthUsecase) {},
			expectedStatus: fiber.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]any) {
				assert.False(t, resp["success"].(bool))
			},
		},
		{
			name:           "error - usecase returns validation error",
			registrationID: registrationID.String(),
			authHeader:     "Bearer " + registrationToken,
			body: map[string]any{
				"full_name":     "John Smith",
				"date_of_birth": "2020-01-01",
				"gender":        "GENDER_001",
			},
			setupMock: func(m *MockAuthUsecase) {
				m.On("CompleteProfileRegistration", mock.Anything, mock.AnythingOfType("*auth.CompleteProfileRegistrationRequest")).
					Return(nil, errors.ErrValidation("You must be at least 18 years old to register"))
			},
			expectedStatus: fiber.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]any) {
				assert.False(t, resp["success"].(bool))
				assert.Contains(t, resp["message"], "18 years old")
			},
		},
		{
			name:           "error - usecase returns conflict (email already registered)",
			registrationID: registrationID.String(),
			authHeader:     "Bearer " + registrationToken,
			body: map[string]any{
				"full_name":     "John Smith",
				"date_of_birth": "1990-01-15",
				"gender":        "GENDER_001",
			},
			setupMock: func(m *MockAuthUsecase) {
				m.On("CompleteProfileRegistration", mock.Anything, mock.AnythingOfType("*auth.CompleteProfileRegistrationRequest")).
					Return(nil, errors.ErrConflict("This email has already been registered"))
			},
			expectedStatus: fiber.StatusConflict,
			checkResponse: func(t *testing.T, resp map[string]any) {
				assert.False(t, resp["success"].(bool))
				assert.Contains(t, resp["message"], "already been registered")
			},
		},
		{
			name:           "error - invalid date_of_birth format rejected at DTO validation layer",
			registrationID: registrationID.String(),
			authHeader:     "Bearer " + registrationToken,
			body: map[string]any{
				"full_name":     "John Smith",
				"date_of_birth": "not-a-date",
				"gender":        "GENDER_001",
			},

			setupMock:      func(m *MockAuthUsecase) {},
			expectedStatus: fiber.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]any) {
				assert.False(t, resp["success"].(bool))
			},
		},
		{
			name:           "error - date_of_birth with wrong separator rejected at DTO validation layer",
			registrationID: registrationID.String(),
			authHeader:     "Bearer " + registrationToken,
			body: map[string]any{
				"full_name":     "John Smith",
				"date_of_birth": "1990/01/15",
				"gender":        "GENDER_001",
			},
			setupMock:      func(m *MockAuthUsecase) {},
			expectedStatus: fiber.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]any) {
				assert.False(t, resp["success"].(bool))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := new(MockAuthUsecase)
			tt.setupMock(mockUC)

			app := setupTestApp()
			ctrl := controller.NewRegistrationController(&config.Config{}, mockUC)
			app.Post("/registrations/:id/complete-profile", ctrl.CompleteProfileRegistration)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/registrations/"+tt.registrationID+"/complete-profile", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var respBody map[string]any
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
			tt.checkResponse(t, respBody)

			mockUC.AssertExpectations(t)
		})
	}
}
