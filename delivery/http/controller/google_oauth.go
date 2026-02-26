package controller

import (
	"erp-service/delivery/http/dto/response"
	"erp-service/delivery/http/presenter"
	"erp-service/iam/auth"
	"erp-service/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
		message = "Account created and logged in successfully"
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		message,
		presenter.ToGoogleCallbackResponse(resp),
	))
}
