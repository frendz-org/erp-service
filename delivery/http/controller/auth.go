package controller

import (
	"strings"

	"iam-service/config"
	"iam-service/delivery/http/dto/response"
	"iam-service/delivery/http/presenter"
	"iam-service/iam/auth"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func convertValidationErrors(errs validator.ValidationErrors) []errors.FieldError {
	result := make([]errors.FieldError, len(errs))
	for i, err := range errs {
		field := err.Field()
		var message string
		switch err.Tag() {
		case "required":
			message = field + " is required"
		case "email":
			message = field + " must be a valid email address"
		case "min":
			message = field + " must be at least " + err.Param() + " characters"
		case "max":
			message = field + " must be at most " + err.Param() + " characters"
		default:
			message = field + " is invalid"
		}
		result[i] = errors.FieldError{Field: field, Message: message}
	}
	return result
}

type AuthController struct {
	config      *config.Config
	authUsecase auth.Usecase
	validate    *validator.Validate
}

func NewRegistrationController(cfg *config.Config, authUsecase auth.Usecase) *AuthController {
	return &AuthController{
		config:      cfg,
		authUsecase: authUsecase,
		validate:    validate,
	}
}
func (rc *AuthController) Logout(c *fiber.Ctx) error {
	var req authdto.LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	userID, jti, tokenExp, err := extractClaimsForLogout(c)
	if err != nil {
		return err
	}

	req.UserID = userID
	req.AccessTokenJTI = jti
	req.AccessTokenExp = tokenExp
	req.IPAddress = getClientIP(c).String()
	req.UserAgent = getUserAgent(c)

	if err := rc.authUsecase.Logout(c.Context(), &req); err != nil {

		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"Logout successful",
		nil,
	))
}
func (rc *AuthController) LogoutAll(c *fiber.Ctx) error {
	userID, _, _, err := extractClaimsForLogout(c)
	if err != nil {
		return err
	}

	req := &authdto.LogoutAllRequest{
		UserID:    userID,
		IPAddress: getClientIP(c).String(),
		UserAgent: getUserAgent(c),
	}

	if err := rc.authUsecase.LogoutAll(c.Context(), req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"All sessions logged out successfully",
		nil,
	))
}

func (rc *AuthController) InitiateRegistration(c *fiber.Ctx) error {
	var req authdto.InitiateRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.IPAddress = getClientIP(c).String()
	req.UserAgent = getUserAgent(c)

	resp, err := rc.authUsecase.InitiateRegistration(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		resp.Message,
		presenter.ToInitiateRegistrationResponse(resp),
	))
}

func (rc *AuthController) VerifyRegistrationOTP(c *fiber.Ctx) error {
	registrationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest("Invalid registration ID format")
	}

	var req authdto.VerifyRegistrationOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.RegistrationID = registrationID

	resp, err := rc.authUsecase.VerifyRegistrationOTP(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		resp.Message,
		presenter.ToVerifyRegistrationOTPResponse(resp),
	))
}

func (rc *AuthController) ResendRegistrationOTP(c *fiber.Ctx) error {
	registrationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest("Invalid registration ID format")
	}

	var req authdto.ResendRegistrationOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.RegistrationID = registrationID

	resp, err := rc.authUsecase.ResendRegistrationOTP(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		resp.Message,
		presenter.ToResendRegistrationOTPResponse(resp),
	))
}

func (rc *AuthController) GetRegistrationStatus(c *fiber.Ctx) error {
	registrationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest("Invalid registration ID format")
	}

	email := c.Query("email")
	if email == "" {
		return errors.ErrBadRequest("email query parameter is required")
	}

	resp, err := rc.authUsecase.GetRegistrationStatus(c.Context(), registrationID, email)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"Registration status retrieved",
		presenter.ToRegistrationStatusResponse(resp),
	))
}

func (rc *AuthController) SetPassword(c *fiber.Ctx) error {
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

	var req authdto.SetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.RegistrationID = registrationID
	req.RegistrationToken = registrationToken

	resp, err := rc.authUsecase.SetPassword(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		resp.Message,
		presenter.ToSetPasswordResponse(resp),
	))
}

func (rc *AuthController) CompleteProfileRegistration(c *fiber.Ctx) error {
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

	var req authdto.CompleteProfileRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.RegistrationID = registrationID
	req.RegistrationToken = registrationToken

	resp, err := rc.authUsecase.CompleteProfileRegistration(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		resp.Message,
		presenter.ToCompleteProfileRegistrationResponse(resp),
	))
}

func (rc *AuthController) CompleteRegistration(c *fiber.Ctx) error {
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

	var req authdto.CompleteRegistrationRequest
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

	resp, err := rc.authUsecase.CompleteRegistration(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		resp.Message,
		presenter.ToCompleteRegistrationResponse(resp),
	))
}
