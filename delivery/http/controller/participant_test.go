package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	jwtpkg "erp-service/pkg/jwt"
	"erp-service/saving/participant/participantdto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockParticipantUsecase struct {
	mock.Mock
}

func (m *MockParticipantUsecase) CreateParticipant(ctx context.Context, req *participantdto.CreateParticipantRequest) (*participantdto.ParticipantResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.ParticipantResponse), args.Error(1)
}

func (m *MockParticipantUsecase) UpdatePersonalData(ctx context.Context, req *participantdto.UpdatePersonalDataRequest) (*participantdto.ParticipantResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.ParticipantResponse), args.Error(1)
}

func (m *MockParticipantUsecase) GetParticipant(ctx context.Context, req *participantdto.GetParticipantRequest) (*participantdto.ParticipantResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.ParticipantResponse), args.Error(1)
}

func (m *MockParticipantUsecase) ListParticipants(ctx context.Context, req *participantdto.ListParticipantsRequest) (*participantdto.ListParticipantsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.ListParticipantsResponse), args.Error(1)
}

func (m *MockParticipantUsecase) DeleteParticipant(ctx context.Context, req *participantdto.DeleteParticipantRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockParticipantUsecase) SaveIdentity(ctx context.Context, req *participantdto.SaveIdentityRequest) (*participantdto.IdentityResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.IdentityResponse), args.Error(1)
}

func (m *MockParticipantUsecase) DeleteIdentity(ctx context.Context, req *participantdto.DeleteChildEntityRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockParticipantUsecase) SaveAddress(ctx context.Context, req *participantdto.SaveAddressRequest) (*participantdto.AddressResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.AddressResponse), args.Error(1)
}

func (m *MockParticipantUsecase) DeleteAddress(ctx context.Context, req *participantdto.DeleteChildEntityRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockParticipantUsecase) SaveBankAccount(ctx context.Context, req *participantdto.SaveBankAccountRequest) (*participantdto.BankAccountResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.BankAccountResponse), args.Error(1)
}

func (m *MockParticipantUsecase) DeleteBankAccount(ctx context.Context, req *participantdto.DeleteChildEntityRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockParticipantUsecase) SaveFamilyMember(ctx context.Context, req *participantdto.SaveFamilyMemberRequest) (*participantdto.FamilyMemberResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.FamilyMemberResponse), args.Error(1)
}

func (m *MockParticipantUsecase) DeleteFamilyMember(ctx context.Context, req *participantdto.DeleteChildEntityRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockParticipantUsecase) SaveEmployment(ctx context.Context, req *participantdto.SaveEmploymentRequest) (*participantdto.EmploymentResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.EmploymentResponse), args.Error(1)
}

func (m *MockParticipantUsecase) SavePension(ctx context.Context, req *participantdto.SavePensionRequest) (*participantdto.PensionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.PensionResponse), args.Error(1)
}

func (m *MockParticipantUsecase) SaveBeneficiary(ctx context.Context, req *participantdto.SaveBeneficiaryRequest) (*participantdto.BeneficiaryResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.BeneficiaryResponse), args.Error(1)
}

func (m *MockParticipantUsecase) DeleteBeneficiary(ctx context.Context, req *participantdto.DeleteChildEntityRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockParticipantUsecase) SaveAddresses(ctx context.Context, req *participantdto.SaveAddressesRequest) ([]participantdto.AddressResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]participantdto.AddressResponse), args.Error(1)
}

func (m *MockParticipantUsecase) SaveFamilyMembers(ctx context.Context, req *participantdto.SaveFamilyMembersRequest) ([]participantdto.FamilyMemberResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]participantdto.FamilyMemberResponse), args.Error(1)
}

func (m *MockParticipantUsecase) SaveBeneficiaries(ctx context.Context, req *participantdto.SaveBeneficiariesRequest) ([]participantdto.BeneficiaryResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]participantdto.BeneficiaryResponse), args.Error(1)
}

func (m *MockParticipantUsecase) UploadFile(ctx context.Context, req *participantdto.UploadFileRequest) (*participantdto.FileUploadResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.FileUploadResponse), args.Error(1)
}

func (m *MockParticipantUsecase) SubmitParticipant(ctx context.Context, req *participantdto.SubmitParticipantRequest) (*participantdto.ParticipantResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.ParticipantResponse), args.Error(1)
}

func (m *MockParticipantUsecase) ApproveParticipant(ctx context.Context, req *participantdto.ApproveParticipantRequest) (*participantdto.ParticipantResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.ParticipantResponse), args.Error(1)
}

func (m *MockParticipantUsecase) RejectParticipant(ctx context.Context, req *participantdto.RejectParticipantRequest) (*participantdto.ParticipantResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.ParticipantResponse), args.Error(1)
}

func (m *MockParticipantUsecase) GetStatusHistory(ctx context.Context, req *participantdto.GetParticipantRequest) ([]participantdto.StatusHistoryResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]participantdto.StatusHistoryResponse), args.Error(1)
}

func (m *MockParticipantUsecase) SelfRegister(ctx context.Context, req *participantdto.SelfRegisterRequest) (*participantdto.SelfRegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*participantdto.SelfRegisterResponse), args.Error(1)
}

func setupParticipantApp(uc *MockParticipantUsecase, userID uuid.UUID) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	ctrl := NewParticipantController(uc)

	app.Post("/self-register", func(c *fiber.Ctx) error {
		c.Locals("user_claims", &jwtpkg.JWTClaims{
			UserID: userID,
		})
		return c.Next()
	}, ctrl.SelfRegister)

	return app
}

func TestSelfRegisterController_HR001_ValidationErrors(t *testing.T) {
	userID := uuid.New()
	mockUC := new(MockParticipantUsecase)
	app := setupParticipantApp(mockUC, userID)

	body := map[string]any{
		"organization": "AB",
	}
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/self-register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var respBody map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	assert.False(t, respBody["success"].(bool))

	mockUC.AssertNotCalled(t, "SelfRegister")
}

func TestSelfRegisterController_HR001_BodyParseError(t *testing.T) {
	userID := uuid.New()
	mockUC := new(MockParticipantUsecase)
	app := setupParticipantApp(mockUC, userID)

	req := httptest.NewRequest("POST", "/self-register", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var respBody map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	assert.False(t, respBody["success"].(bool))
	mockUC.AssertNotCalled(t, "SelfRegister")
}

func TestSelfRegisterController_HR002_PresenterUsed(t *testing.T) {
	userID := uuid.New()
	mockUC := new(MockParticipantUsecase)
	app := setupParticipantApp(mockUC, userID)

	expectedResp := &participantdto.SelfRegisterResponse{
		IsLinked:           false,
		RegistrationStatus: "PENDING_APPROVAL",
		Data: &participantdto.SelfRegisterParticipantData{
			ParticipantNumber: "ABC12345",
			Status:            "DRAFT",
		},
	}

	mockUC.On("SelfRegister", mock.Anything, mock.AnythingOfType("*participantdto.SelfRegisterRequest")).
		Return(expectedResp, nil)

	body := map[string]any{
		"organization":       "TENANT_001",
		"identity_number":    "3201011501900001",
		"participant_number": "ABC12345",
		"phone_number":       "+6281234567890",
	}
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/self-register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var respBody map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	assert.True(t, respBody["success"].(bool))

	data := respBody["data"].(map[string]any)
	assert.Equal(t, false, data["is_linked"])
	assert.Equal(t, "PENDING_APPROVAL", data["registration_status"])

	mockUC.AssertExpectations(t)
}

func TestSelfRegisterController_HR002_IsLinked_Returns200(t *testing.T) {
	userID := uuid.New()
	mockUC := new(MockParticipantUsecase)
	app := setupParticipantApp(mockUC, userID)

	expectedResp := &participantdto.SelfRegisterResponse{
		IsLinked:           true,
		RegistrationStatus: "PENDING_APPROVAL",
		Data: &participantdto.SelfRegisterParticipantData{
			ParticipantNumber: "ABC12345",
			Status:            "DRAFT",
		},
	}

	mockUC.On("SelfRegister", mock.Anything, mock.AnythingOfType("*participantdto.SelfRegisterRequest")).
		Return(expectedResp, nil)

	body := map[string]any{
		"organization":       "TENANT_001",
		"identity_number":    "3201011501900001",
		"participant_number": "ABC12345",
		"phone_number":       "+6281234567890",
	}
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/self-register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockUC.AssertExpectations(t)
}
