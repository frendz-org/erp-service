package presenter

import (
	"erp-service/delivery/http/dto/response"
	"erp-service/iam/auth"
)

func ToInitiateRegistrationResponse(resp *auth.InitiateRegistrationResponse) *response.InitiateRegistrationResponse {
	if resp == nil {
		return nil
	}
	return &response.InitiateRegistrationResponse{
		RegistrationID: resp.RegistrationID,
		Email:          resp.Email,
		Status:         resp.Status,
		Message:        resp.Message,
		ExpiresAt:      resp.ExpiresAt,
		OTPConfig: response.InitiateRegistrationOTPConfig{
			ExpiresInMinutes:      resp.OTPConfig.ExpiresInMinutes,
			ResendCooldownSeconds: resp.OTPConfig.ResendCooldownSeconds,
		},
	}
}

func ToVerifyRegistrationOTPResponse(resp *auth.VerifyRegistrationOTPResponse) *response.VerifyRegistrationOTPResponse {
	if resp == nil {
		return nil
	}
	return &response.VerifyRegistrationOTPResponse{
		RegistrationID:    resp.RegistrationID,
		Status:            resp.Status,
		Message:           resp.Message,
		RegistrationToken: resp.RegistrationToken,
		TokenExpiresAt:    resp.TokenExpiresAt,
		NextStep: response.VerifyRegistrationOTPNextStep{
			Action:   resp.NextStep.Action,
			Endpoint: resp.NextStep.Endpoint,
		},
	}
}

func ToResendRegistrationOTPResponse(resp *auth.ResendRegistrationOTPResponse) *response.ResendRegistrationOTPResponse {
	if resp == nil {
		return nil
	}
	return &response.ResendRegistrationOTPResponse{
		RegistrationID:        resp.RegistrationID,
		Message:               resp.Message,
		ExpiresAt:             resp.ExpiresAt,
		ResendsRemaining:      resp.ResendsRemaining,
		NextResendAvailableAt: resp.NextResendAvailableAt,
	}
}

func ToRegistrationStatusResponse(resp *auth.RegistrationStatusResponse) *response.RegistrationStatusResponse {
	if resp == nil {
		return nil
	}
	return &response.RegistrationStatusResponse{
		RegistrationID:       resp.RegistrationID,
		Email:                resp.Email,
		Status:               resp.Status,
		ExpiresAt:            resp.ExpiresAt,
		OTPAttemptsRemaining: resp.OTPAttemptsRemaining,
		ResendsRemaining:     resp.ResendsRemaining,
	}
}

func ToSetPasswordResponse(resp *auth.SetPasswordResponse) *response.SetPasswordResponse {
	if resp == nil {
		return nil
	}
	return &response.SetPasswordResponse{
		RegistrationID:    resp.RegistrationID,
		Status:            resp.Status,
		Message:           resp.Message,
		RegistrationToken: resp.RegistrationToken,
		NextStep: response.SetPasswordNextStep{
			Action:         resp.NextStep.Action,
			Endpoint:       resp.NextStep.Endpoint,
			RequiredFields: resp.NextStep.RequiredFields,
		},
	}
}

func ToCompleteProfileRegistrationResponse(resp *auth.CompleteProfileRegistrationResponse) *response.CompleteProfileRegistrationResponse {
	if resp == nil {
		return nil
	}
	return &response.CompleteProfileRegistrationResponse{
		UserID:  resp.UserID,
		Email:   resp.Email,
		Status:  resp.Status,
		Message: resp.Message,
		Profile: response.CompleteProfileRegistrationProfile{
			FirstName: resp.Profile.FirstName,
			LastName:  resp.Profile.LastName,
		},
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		TokenType:    resp.TokenType,
		ExpiresIn:    resp.ExpiresIn,
	}
}

