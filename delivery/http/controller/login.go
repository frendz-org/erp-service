package controller

import (
	"erp-service/delivery/http/dto/response"
	"erp-service/delivery/http/presenter"
	"erp-service/iam/auth/authdto"
	"erp-service/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (rc *AuthController) RefreshToken(c *fiber.Ctx) error {
	var req authdto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.IPAddress = getClientIP(c).String()
	req.UserAgent = getUserAgent(c)

	resp, err := rc.authUsecase.RefreshToken(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"Token refreshed successfully",
		presenter.ToRefreshTokenResponse(resp),
	))
}

func (rc *AuthController) InitiateLogin(c *fiber.Ctx) error {
	var req authdto.InitiateLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.IPAddress = getClientIP(c).String()
	req.UserAgent = getUserAgent(c)

	resp, err := rc.authUsecase.InitiateLogin(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"OTP sent to your email",
		presenter.ToUnifiedLoginResponse(resp),
	))
}

func (rc *AuthController) VerifyLoginOTP(c *fiber.Ctx) error {
	loginSessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest("Invalid login session ID format")
	}

	var req authdto.VerifyLoginOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.LoginSessionID = loginSessionID
	req.IPAddress = getClientIP(c).String()
	req.UserAgent = getUserAgent(c)

	resp, err := rc.authUsecase.VerifyLoginOTP(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"Login successful",
		presenter.ToVerifyLoginOTPResponse(resp),
	))
}

func (rc *AuthController) ResendLoginOTP(c *fiber.Ctx) error {
	loginSessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest("Invalid login session ID format")
	}

	var req authdto.ResendLoginOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.LoginSessionID = loginSessionID

	resp, err := rc.authUsecase.ResendLoginOTP(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"OTP resent successfully",
		presenter.ToResendLoginOTPResponse2(resp),
	))
}

func (rc *AuthController) GetLoginStatus(c *fiber.Ctx) error {
	loginSessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest("Invalid login session ID format")
	}

	email := c.Query("email")
	if email == "" {
		return errors.ErrBadRequest("email query parameter is required")
	}

	resp, err := rc.authUsecase.GetLoginStatus(c.Context(), &authdto.GetLoginStatusRequest{
		LoginSessionID: loginSessionID,
		Email:          email,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"Login status retrieved",
		presenter.ToLoginStatusResponse(resp),
	))
}
