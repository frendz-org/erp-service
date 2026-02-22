package presenter

import (
	"erp-service/delivery/http/dto/response"
	"erp-service/iam/auth"
)

func ToUnifiedLoginResponse(resp *auth.UnifiedLoginResponse) *response.UnifiedLoginResponse {
	if resp == nil {
		return nil
	}

	result := &response.UnifiedLoginResponse{
		Status:          string(resp.Status),
		LoginSessionID:  resp.LoginSessionID,
		Email:           resp.Email,
		OTPExpiresAt:    resp.OTPExpiresAt,
		AttemptsAllowed: resp.AttemptsAllowed,
		ResendsAllowed:  resp.ResendsAllowed,
		SessionExpires:  resp.SessionExpires,
		AccessToken:     resp.AccessToken,
		RefreshToken:    resp.RefreshToken,
		ExpiresIn:       resp.ExpiresIn,
		TokenType:       resp.TokenType,
	}

	if resp.User != nil {
		result.User = toLoginUserResponse(resp.User)
	}

	return result
}

func ToVerifyLoginOTPResponse(resp *auth.VerifyLoginOTPResponse) *response.VerifyLoginOTPResponse {
	if resp == nil {
		return nil
	}
	return &response.VerifyLoginOTPResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		TokenType:    resp.TokenType,
		User:         *toLoginUserResponse(&resp.User),
	}
}

func ToResendLoginOTPResponse2(resp *auth.ResendLoginOTPResponse) *response.ResendLoginOTPResponse {
	if resp == nil {
		return nil
	}
	return &response.ResendLoginOTPResponse{
		Status:           resp.Status,
		LoginSessionID:   resp.LoginSessionID,
		Email:            resp.Email,
		OTPExpiresAt:     resp.OTPExpiresAt,
		ResendsRemaining: resp.ResendsRemaining,
		CooldownSeconds:  resp.CooldownSeconds,
	}
}

func ToLoginStatusResponse(resp *auth.LoginStatusResponse) *response.LoginStatusResponse {
	if resp == nil {
		return nil
	}
	return &response.LoginStatusResponse{
		Status:            resp.Status,
		LoginSessionID:    resp.LoginSessionID,
		Email:             resp.Email,
		AttemptsRemaining: resp.AttemptsRemaining,
		ResendsRemaining:  resp.ResendsRemaining,
		ExpiresAt:         resp.ExpiresAt,
		CooldownRemaining: resp.CooldownRemaining,
	}
}

func ToRefreshTokenResponse(resp *auth.RefreshTokenResponse) *response.RefreshTokenResponse {
	if resp == nil {
		return nil
	}
	return &response.RefreshTokenResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		TokenType:    resp.TokenType,
		User:         *toLoginUserResponse(&resp.User),
	}
}

func toLoginUserResponse(user *auth.LoginUserResponse) *response.LoginUserResponse {
	if user == nil {
		return nil
	}

	result := &response.LoginUserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
	}

	for _, t := range user.Tenants {
		tenant := response.LoginTenantResponse{
			TenantID: t.TenantID,
		}
		for _, p := range t.Products {
			tenant.Products = append(tenant.Products, response.LoginProductResponse{
				ProductID:   p.ProductID,
				ProductCode: p.ProductCode,
				Roles:       p.Roles,
				Permissions: p.Permissions,
			})
		}
		result.Tenants = append(result.Tenants, tenant)
	}

	return result
}
