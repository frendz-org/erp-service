package controller

import (
	"strconv"

	"iam-service/delivery/http/middleware"
	"iam-service/delivery/http/presenter"
	"iam-service/pkg/errors"
	"iam-service/saving/member"
	"iam-service/saving/member/memberdto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MemberController struct {
	usecase member.Usecase
}

func NewMemberController(uc member.Usecase) *MemberController {
	return &MemberController{
		usecase: uc,
	}
}

func (ctrl *MemberController) Register(c *fiber.Ctx) error {
	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	req := &memberdto.RegisterRequest{
		TenantID:      tenantID,
		ProductID: productID,
		UserID:        userClaims.UserID,
	}

	result, err := ctrl.usecase.RegisterMember(c.UserContext(), req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapMemberRegisterResponse(result),
	})
}

func (ctrl *MemberController) List(c *fiber.Ctx) error {
	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	req := &memberdto.ListRequest{
		TenantID:      tenantID,
		ProductID: productID,
		Search:        c.Query("search"),
		Page:          page,
		PerPage:       perPage,
		SortBy:        c.Query("sort_by", "created_at"),
		SortOrder:     c.Query("sort_order", "desc"),
	}

	if status := c.Query("status"); status != "" {
		req.Status = &status
	}

	result, err := ctrl.usecase.ListMembers(c.UserContext(), req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapMemberListResponse(result),
	})
}

func (ctrl *MemberController) Get(c *fiber.Ctx) error {
	memberID, err := uuid.Parse(c.Params("memberId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid member ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	req := &memberdto.GetMemberRequest{
		MemberID:      memberID,
		TenantID:      tenantID,
		ProductID: productID,
	}

	result, err := ctrl.usecase.GetMember(c.UserContext(), req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapMemberDetailResponse(result),
	})
}

func (ctrl *MemberController) Approve(c *fiber.Ctx) error {
	memberID, err := uuid.Parse(c.Params("memberId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid member ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	var body struct {
		RoleCode string `json:"role_code" validate:"required"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := validate.Struct(&body); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req := &memberdto.ApproveRequest{
		MemberID:      memberID,
		TenantID:      tenantID,
		ProductID: productID,
		ApproverID:    userClaims.UserID,
		RoleCode:      body.RoleCode,
	}

	result, err := ctrl.usecase.ApproveMember(c.UserContext(), req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapMemberDetailResponse(result),
	})
}

func (ctrl *MemberController) Reject(c *fiber.Ctx) error {
	memberID, err := uuid.Parse(c.Params("memberId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid member ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	var body struct {
		Reason string `json:"reason" validate:"required,min=10,max=500"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := validate.Struct(&body); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req := &memberdto.RejectRequest{
		MemberID:      memberID,
		TenantID:      tenantID,
		ProductID: productID,
		ApproverID:    userClaims.UserID,
		Reason:        body.Reason,
	}

	result, err := ctrl.usecase.RejectMember(c.UserContext(), req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapMemberDetailResponse(result),
	})
}

func (ctrl *MemberController) ChangeRole(c *fiber.Ctx) error {
	memberID, err := uuid.Parse(c.Params("memberId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid member ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	var body struct {
		RoleCode string `json:"role_code" validate:"required"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := validate.Struct(&body); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req := &memberdto.ChangeRoleRequest{
		MemberID:      memberID,
		TenantID:      tenantID,
		ProductID: productID,
		ActorID:       userClaims.UserID,
		RoleCode:      body.RoleCode,
	}

	result, err := ctrl.usecase.ChangeRole(c.UserContext(), req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapMemberDetailResponse(result),
	})
}

func (ctrl *MemberController) Deactivate(c *fiber.Ctx) error {
	memberID, err := uuid.Parse(c.Params("memberId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid member ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	req := &memberdto.DeactivateRequest{
		MemberID:      memberID,
		TenantID:      tenantID,
		ProductID: productID,
		ActorID:       userClaims.UserID,
	}

	result, err := ctrl.usecase.DeactivateMember(c.UserContext(), req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapMemberDetailResponse(result),
	})
}
