package controller

import (
	"iam-service/config"
	"iam-service/delivery/http/dto/response"
	"iam-service/delivery/http/presenter"
	"iam-service/masterdata"
	"iam-service/masterdata/masterdatadto"
	"iam-service/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MasterdataController struct {
	config   *config.Config
	usecase  masterdata.Usecase
	validate *validator.Validate
}

func NewMasterdataController(cfg *config.Config, usecase masterdata.Usecase) *MasterdataController {
	return &MasterdataController{
		config:   cfg,
		usecase:  usecase,
		validate: validate,
	}
}

func (mc *MasterdataController) ListCategories(c *fiber.Ctx) error {
	var req masterdatadto.ListCategoriesRequest
	if err := c.QueryParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid query parameters")
	}

	resp, err := mc.usecase.ListCategories(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(response.APIResponse{
		Success: true,
		Message: "Categories retrieved successfully",
		Data:    presenter.ToCategoryListResponse(resp.Categories),
		Pagination: &response.Pagination{
			Total:      resp.Pagination.Total,
			Page:       resp.Pagination.Page,
			Limit:      resp.Pagination.PerPage,
			TotalPages: resp.Pagination.TotalPages,
		},
	})
}

func (mc *MasterdataController) GetCategoryByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest("Invalid category ID format")
	}

	resp, err := mc.usecase.GetCategoryByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Category retrieved successfully", presenter.ToCategoryResponse(resp)))
}

func (mc *MasterdataController) GetCategoryByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return errors.ErrBadRequest("Category code is required")
	}

	resp, err := mc.usecase.GetCategoryByCode(c.Context(), code)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Category retrieved successfully", presenter.ToCategoryResponse(resp)))
}

func (mc *MasterdataController) CreateCategory(c *fiber.Ctx) error {
	var req masterdatadto.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := mc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertMasterdataValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := mc.usecase.CreateCategory(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		"Category created successfully",
		presenter.ToCategoryResponse(resp),
	))
}

func (mc *MasterdataController) UpdateCategory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest("Invalid category ID format")
	}

	var req masterdatadto.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := mc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertMasterdataValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := mc.usecase.UpdateCategory(c.Context(), id, &req)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Category updated successfully", presenter.ToCategoryResponse(resp)))
}

func (mc *MasterdataController) DeleteCategory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest("Invalid category ID format")
	}

	if err := mc.usecase.DeleteCategory(c.Context(), id); err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Category deleted successfully", nil))
}

func (mc *MasterdataController) GetCategoryChildren(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest("Invalid category ID format")
	}

	resp, err := mc.usecase.GetCategoryChildren(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Child categories retrieved successfully", presenter.ToCategoryListResponse(resp)))
}

func (mc *MasterdataController) ListItems(c *fiber.Ctx) error {
	var req masterdatadto.ListItemsRequest
	if err := c.QueryParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid query parameters")
	}

	resp, err := mc.usecase.ListItems(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(response.APIResponse{
		Success: true,
		Message: "Items retrieved successfully",
		Data:    presenter.ToItemListResponse(resp.Items),
		Pagination: &response.Pagination{
			Total:      resp.Pagination.Total,
			Page:       resp.Pagination.Page,
			Limit:      resp.Pagination.PerPage,
			TotalPages: resp.Pagination.TotalPages,
		},
	})
}

func (mc *MasterdataController) GetItemByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest("Invalid item ID format")
	}

	resp, err := mc.usecase.GetItemByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Item retrieved successfully", presenter.ToItemResponse(resp)))
}

func (mc *MasterdataController) GetItemByCode(c *fiber.Ctx) error {
	categoryCode := c.Params("categoryCode")
	if categoryCode == "" {
		return errors.ErrBadRequest("Category code is required")
	}

	itemCode := c.Params("itemCode")
	if itemCode == "" {
		return errors.ErrBadRequest("Item code is required")
	}

	var tenantID *uuid.UUID
	if tenantIDStr := c.Query("tenant_id"); tenantIDStr != "" {
		id, err := uuid.Parse(tenantIDStr)
		if err != nil {
			return errors.ErrBadRequest("Invalid tenant ID format")
		}
		tenantID = &id
	}

	resp, err := mc.usecase.GetItemByCode(c.Context(), categoryCode, tenantID, itemCode)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Item retrieved successfully", presenter.ToItemResponse(resp)))
}

func (mc *MasterdataController) CreateItem(c *fiber.Ctx) error {
	var req masterdatadto.CreateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := mc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertMasterdataValidationErrors(err.(validator.ValidationErrors)))
	}

	if userID, err := getUserID(c); err == nil {
		req.CreatedBy = &userID
	}

	resp, err := mc.usecase.CreateItem(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		"Item created successfully",
		presenter.ToItemResponse(resp),
	))
}

func (mc *MasterdataController) UpdateItem(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest("Invalid item ID format")
	}

	var req masterdatadto.UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := mc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertMasterdataValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := mc.usecase.UpdateItem(c.Context(), id, &req)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Item updated successfully", presenter.ToItemResponse(resp)))
}

func (mc *MasterdataController) DeleteItem(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest("Invalid item ID format")
	}

	if err := mc.usecase.DeleteItem(c.Context(), id); err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Item deleted successfully", nil))
}

func (mc *MasterdataController) GetItemChildren(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrBadRequest("Invalid item ID format")
	}

	resp, err := mc.usecase.GetItemChildren(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Child items retrieved successfully", presenter.ToItemListResponse(resp)))
}

func (mc *MasterdataController) GetItemTree(c *fiber.Ctx) error {
	categoryCode := c.Params("categoryCode")
	if categoryCode == "" {
		return errors.ErrBadRequest("Category code is required")
	}

	var tenantID *uuid.UUID
	if tenantIDStr := c.Query("tenant_id"); tenantIDStr != "" {
		id, err := uuid.Parse(tenantIDStr)
		if err != nil {
			return errors.ErrBadRequest("Invalid tenant ID format")
		}
		tenantID = &id
	}

	resp, err := mc.usecase.GetItemTree(c.Context(), categoryCode, tenantID)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Item tree retrieved successfully", presenter.ToItemTreeResponse(resp)))
}

func (mc *MasterdataController) ListItemsByParent(c *fiber.Ctx) error {
	categoryCode := c.Params("categoryCode")
	if categoryCode == "" {
		return errors.ErrBadRequest("Category code is required")
	}

	parentCode := c.Params("parentCode")
	if parentCode == "" {
		return errors.ErrBadRequest("Parent code is required")
	}

	var tenantID *uuid.UUID
	if tenantIDStr := c.Query("tenant_id"); tenantIDStr != "" {
		id, err := uuid.Parse(tenantIDStr)
		if err != nil {
			return errors.ErrBadRequest("Invalid tenant ID format")
		}
		tenantID = &id
	}

	resp, err := mc.usecase.ListItemsByParent(c.Context(), categoryCode, parentCode, tenantID)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Items retrieved successfully", presenter.ToItemListResponse(resp)))
}

func (mc *MasterdataController) GetDefaultItem(c *fiber.Ctx) error {
	categoryCode := c.Params("categoryCode")
	if categoryCode == "" {
		return errors.ErrBadRequest("Category code is required")
	}

	var tenantID *uuid.UUID
	if tenantIDStr := c.Query("tenant_id"); tenantIDStr != "" {
		id, err := uuid.Parse(tenantIDStr)
		if err != nil {
			return errors.ErrBadRequest("Invalid tenant ID format")
		}
		tenantID = &id
	}

	resp, err := mc.usecase.GetItemDefault(c.Context(), categoryCode, tenantID)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Default item retrieved successfully", presenter.ToItemResponse(resp)))
}

func (mc *MasterdataController) ValidateCode(c *fiber.Ctx) error {
	var req masterdatadto.ValidateCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := mc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertMasterdataValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := mc.usecase.ValidateItemCode(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Validation completed", presenter.ToValidateCodeResponse(resp)))
}

func (mc *MasterdataController) ValidateCodes(c *fiber.Ctx) error {
	var req masterdatadto.ValidateCodesRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := mc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertMasterdataValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := mc.usecase.ValidateItemCodes(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Batch validation completed", presenter.ToValidateCodesResponse(resp)))
}

func convertMasterdataValidationErrors(errs validator.ValidationErrors) []errors.FieldError {
	result := make([]errors.FieldError, len(errs))
	for i, err := range errs {
		field := err.Field()
		var message string
		switch err.Tag() {
		case "required":
			message = field + " is required"
		case "max":
			message = field + " must be at most " + err.Param() + " characters"
		case "min":
			message = field + " must be at least " + err.Param()
		case "oneof":
			message = field + " must be one of: " + err.Param()
		case "uuid":
			message = field + " must be a valid UUID"
		default:
			message = field + " is invalid"
		}
		result[i] = errors.FieldError{Field: field, Message: message}
	}
	return result
}
