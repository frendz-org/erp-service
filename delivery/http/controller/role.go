package controller

import (
	"iam-service/config"
	"iam-service/delivery/http/dto/response"
	"iam-service/iam/role"
	"iam-service/iam/role/roledto"
	"iam-service/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func convertRoleValidationErrors(errs validator.ValidationErrors) []errors.FieldError {
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

type RoleController struct {
	config      *config.Config
	roleUsecase role.Usecase
	validate    *validator.Validate
}

func NewRoleController(cfg *config.Config, roleUsecase role.Usecase) *RoleController {
	return &RoleController{
		config:      cfg,
		roleUsecase: roleUsecase,
		validate:    validate,
	}
}

func (rc *RoleController) Create(c *fiber.Ctx) error {
	var req roledto.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertRoleValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := rc.roleUsecase.Create(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		"Role created successfully",
		resp,
	))
}
