package controller

import (
	"strings"

	"erp-service/delivery/http/dto/response"
	"erp-service/delivery/http/presenter"
	"erp-service/iam/auth"
	"erp-service/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (rc *AuthController) GoogleLogin(c *fiber.Ctx) error {
	resp, err := rc.authUsecase.GetGoogleAuthURL(c.Context())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"Google OAuth URL generated",
		presenter.ToGoogleAuthURLResponse(resp),
	))
}

func (rc *AuthController) GoogleCallback(c *fiber.Ctx) error {
	var req auth.GoogleCallbackRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.IPAddress = getClientIP(c).String()
	req.UserAgent = getUserAgent(c)

	resp, err := rc.authUsecase.HandleGoogleCallback(c.Context(), &req)
	if err != nil {
		return err
	}

	message := "Login successful"
	if resp.IsNewUser {
		message = "Profile completion required"
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		message,
		presenter.ToGoogleCallbackResponse(resp),
	))
}

func (rc *AuthController) CompleteGoogleProfile(c *fiber.Ctx) error {
	registrationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest("Invalid registration ID format")
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return errors.ErrUnauthorized("Authorization header is required")
	}

	registrationToken := strings.TrimPrefix(authHeader, "Bearer ")
	if registrationToken == authHeader {
		return errors.ErrUnauthorized("Invalid authorization format. Use: Bearer <token>")
	}

	var req auth.CompleteGoogleProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.RegistrationID = registrationID
	req.RegistrationToken = registrationToken
	req.IPAddress = getClientIP(c).String()
	req.UserAgent = getUserAgent(c)

	resp, err := rc.authUsecase.CompleteGoogleProfile(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		resp.Message,
		presenter.ToCompleteGoogleProfileResponse(resp),
	))
}
