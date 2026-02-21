package controller

import (
	stderrors "errors"
	"io"
	"net/http"
	"strconv"

	"erp-service/delivery/http/middleware"
	"erp-service/delivery/http/presenter"
	"erp-service/pkg/errors"
	"erp-service/saving/participant"
	"erp-service/saving/participant/participantdto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var allowedUploadContentTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"application/pdf": true,
}

type ParticipantController struct {
	usecase participant.Usecase
}

func NewParticipantController(uc participant.Usecase) *ParticipantController {
	return &ParticipantController{
		usecase: uc,
	}
}

func participantError(c *fiber.Ctx, err error) error {
	appErr := errors.GetAppError(err)
	if appErr == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "internal server error",
			"code":    errors.CodeInternal,
		})
	}
	resp := fiber.Map{
		"success": false,
		"error":   appErr.Message,
		"code":    appErr.Code,
	}
	if len(appErr.Details) > 0 {
		resp["details"] = appErr.Details
	}
	return c.Status(appErr.HTTPStatus).JSON(resp)
}

func (ctrl *ParticipantController) Create(c *fiber.Ctx) error {
	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		return participantError(c, err)
	}

	var req participantdto.CreateParticipantRequest
	if err := c.BodyParser(&req); err != nil {
		return participantError(c, errors.ErrBadRequest("invalid request body"))
	}

	if err := validate.Struct(&req); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req.TenantID = tenantID
	req.ProductID = productID
	req.UserID = userClaims.UserID

	result, err := ctrl.usecase.CreateParticipant(c.UserContext(), &req)
	if err != nil {
		return participantError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapParticipantResponse(result),
	})
}

func (ctrl *ParticipantController) Get(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	result, err := ctrl.usecase.GetParticipant(c.UserContext(), &participantdto.GetParticipantRequest{
		ParticipantID: pID,
		TenantID:      tenantID,
		ProductID:     productID,
	})
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapParticipantResponse(result),
	})
}

func (ctrl *ParticipantController) List(c *fiber.Ctx) error {
	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	perPage, err := strconv.Atoi(c.Query("per_page", "10"))
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req := &participantdto.ListParticipantsRequest{
		TenantID:  tenantID,
		ProductID: productID,
		Search:    c.Query("search"),
		Status:    nil,
		Page:      page,
		PerPage:   perPage,
		SortBy:    c.Query("sort_by", "created_at"),
		SortOrder: c.Query("sort_order", "desc"),
	}

	if status := c.Query("status"); status != "" {
		req.Status = &status
	}

	if err := validate.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	result, err := ctrl.usecase.ListParticipants(c.UserContext(), req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) UpdatePersonalData(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		return participantError(c, err)
	}

	var req participantdto.UpdatePersonalDataRequest
	if err := c.BodyParser(&req); err != nil {
		return participantError(c, errors.ErrBadRequest("invalid request body"))
	}

	if err := validate.Struct(&req); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req.TenantID = tenantID
	req.ParticipantID = pID
	req.UserID = userClaims.UserID
	req.ProductID = productID

	result, err := ctrl.usecase.UpdatePersonalData(c.UserContext(), &req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapParticipantResponse(result),
	})
}

func (ctrl *ParticipantController) SaveIdentity(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	var req participantdto.SaveIdentityRequest
	if err := c.BodyParser(&req); err != nil {
		return participantError(c, errors.ErrBadRequest("invalid request body"))
	}

	if err := validate.Struct(&req); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req.TenantID = tenantID
	req.ParticipantID = pID
	req.ProductID = productID

	result, err := ctrl.usecase.SaveIdentity(c.UserContext(), &req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) DeleteIdentity(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	iID, err := uuid.Parse(c.Params("identityId"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid identity ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	if err := ctrl.usecase.DeleteIdentity(c.UserContext(), &participantdto.DeleteChildEntityRequest{
		ChildID:       iID,
		ParticipantID: pID,
		TenantID:      tenantID,
		ProductID:     productID,
	}); err != nil {
		return participantError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ctrl *ParticipantController) SaveAddress(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	var req participantdto.SaveAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return participantError(c, errors.ErrBadRequest("invalid request body"))
	}

	if err := validate.Struct(&req); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req.TenantID = tenantID
	req.ParticipantID = pID
	req.ProductID = productID

	result, err := ctrl.usecase.SaveAddress(c.UserContext(), &req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) DeleteAddress(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	aID, err := uuid.Parse(c.Params("addressId"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid address ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	if err := ctrl.usecase.DeleteAddress(c.UserContext(), &participantdto.DeleteChildEntityRequest{
		ChildID:       aID,
		ParticipantID: pID,
		TenantID:      tenantID,
		ProductID:     productID,
	}); err != nil {
		return participantError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ctrl *ParticipantController) SaveBankAccount(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	var req participantdto.SaveBankAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return participantError(c, errors.ErrBadRequest("invalid request body"))
	}

	if err := validate.Struct(&req); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req.TenantID = tenantID
	req.ParticipantID = pID
	req.ProductID = productID

	result, err := ctrl.usecase.SaveBankAccount(c.UserContext(), &req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) DeleteBankAccount(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	aID, err := uuid.Parse(c.Params("accountId"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid bank account ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	if err := ctrl.usecase.DeleteBankAccount(c.UserContext(), &participantdto.DeleteChildEntityRequest{
		ChildID:       aID,
		ParticipantID: pID,
		TenantID:      tenantID,
		ProductID:     productID,
	}); err != nil {
		return participantError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ctrl *ParticipantController) SaveFamilyMember(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	var req participantdto.SaveFamilyMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return participantError(c, errors.ErrBadRequest("invalid request body"))
	}

	if err := validate.Struct(&req); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req.TenantID = tenantID
	req.ParticipantID = pID
	req.ProductID = productID

	result, err := ctrl.usecase.SaveFamilyMember(c.UserContext(), &req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) DeleteFamilyMember(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	mID, err := uuid.Parse(c.Params("memberId"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid family member ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	if err := ctrl.usecase.DeleteFamilyMember(c.UserContext(), &participantdto.DeleteChildEntityRequest{
		ChildID:       mID,
		ParticipantID: pID,
		TenantID:      tenantID,
		ProductID:     productID,
	}); err != nil {
		return participantError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ctrl *ParticipantController) SaveEmployment(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	var req participantdto.SaveEmploymentRequest
	if err := c.BodyParser(&req); err != nil {
		return participantError(c, errors.ErrBadRequest("invalid request body"))
	}

	if err := validate.Struct(&req); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req.TenantID = tenantID
	req.ParticipantID = pID
	req.ProductID = productID

	result, err := ctrl.usecase.SaveEmployment(c.UserContext(), &req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapEmploymentResponse(result),
	})
}

func (ctrl *ParticipantController) SavePension(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	var req participantdto.SavePensionRequest
	if err := c.BodyParser(&req); err != nil {
		return participantError(c, errors.ErrBadRequest("invalid request body"))
	}

	if err := validate.Struct(&req); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req.TenantID = tenantID
	req.ParticipantID = pID
	req.ProductID = productID

	result, err := ctrl.usecase.SavePension(c.UserContext(), &req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapPensionResponse(result),
	})
}

func (ctrl *ParticipantController) SaveBeneficiary(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	var req participantdto.SaveBeneficiaryRequest
	if err := c.BodyParser(&req); err != nil {
		return participantError(c, errors.ErrBadRequest("invalid request body"))
	}

	if err := validate.Struct(&req); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req.TenantID = tenantID
	req.ParticipantID = pID
	req.ProductID = productID

	result, err := ctrl.usecase.SaveBeneficiary(c.UserContext(), &req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) DeleteBeneficiary(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	bID, err := uuid.Parse(c.Params("beneficiaryId"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid beneficiary ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	if err := ctrl.usecase.DeleteBeneficiary(c.UserContext(), &participantdto.DeleteChildEntityRequest{
		ChildID:       bID,
		ParticipantID: pID,
		TenantID:      tenantID,
		ProductID:     productID,
	}); err != nil {
		return participantError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ctrl *ParticipantController) UploadFile(c *fiber.Ctx) error {
	const maxFileSize = 5 * 1024 * 1024

	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	fieldName := c.FormValue("field_name")
	if fieldName == "" {
		return participantError(c, errors.ErrBadRequest("field_name is required"))
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return participantError(c, errors.ErrBadRequest("file is required"))
	}

	if fileHeader.Size > maxFileSize {
		return participantError(c, errors.ErrBadRequest("file size exceeds 5MB limit"))
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedUploadContentTypes[contentType] {
		return participantError(c, errors.ErrBadRequest("unsupported file type; allowed: jpeg, png, gif, pdf"))
	}

	file, err := fileHeader.Open()
	if err != nil {
		return participantError(c, errors.ErrInternal("failed to open uploaded file"))
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return participantError(c, errors.ErrInternal("failed to read uploaded file"))
	}
	detectedType := http.DetectContentType(buf[:n])
	if !allowedUploadContentTypes[detectedType] {
		return participantError(c, errors.ErrBadRequest("file content does not match an allowed type; allowed: jpeg, png, gif, pdf"))
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return participantError(c, errors.ErrInternal("failed to process uploaded file"))
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	uploaderID, err := middleware.GetUserID(c)
	if err != nil {
		return participantError(c, err)
	}

	req := &participantdto.UploadFileRequest{
		TenantID:      tenantID,
		ParticipantID: pID,
		ProductID:     productID,
		UploadedBy:    uploaderID,
		FieldName:     fieldName,
		FileName:      fileHeader.Filename,
		ContentType:   detectedType,
		Reader:        file,
		Size:          fileHeader.Size,
	}

	result, err := ctrl.usecase.UploadFile(c.UserContext(), req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) GetStatusHistory(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	result, err := ctrl.usecase.GetStatusHistory(c.UserContext(), &participantdto.GetParticipantRequest{
		ParticipantID: pID,
		TenantID:      tenantID,
		ProductID:     productID,
	})
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) Submit(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		return participantError(c, err)
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req := &participantdto.SubmitParticipantRequest{
		TenantID:      tenantID,
		ParticipantID: pID,
		ProductID:     productID,
		UserID:        userClaims.UserID,
	}

	result, err := ctrl.usecase.SubmitParticipant(c.UserContext(), req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapParticipantResponse(result),
	})
}

func (ctrl *ParticipantController) Approve(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		return participantError(c, err)
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req := &participantdto.ApproveParticipantRequest{
		TenantID:      tenantID,
		ParticipantID: pID,
		ProductID:     productID,
		UserID:        userClaims.UserID,
	}

	result, err := ctrl.usecase.ApproveParticipant(c.UserContext(), req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapParticipantResponse(result),
	})
}

func (ctrl *ParticipantController) Reject(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		return participantError(c, err)
	}

	var body struct {
		Reason string `json:"reason" validate:"required,min=10,max=500"`
	}
	if err := c.BodyParser(&body); err != nil {
		return participantError(c, errors.ErrBadRequest("invalid request body"))
	}

	if err := validate.Struct(&body); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req := &participantdto.RejectParticipantRequest{
		TenantID:      tenantID,
		ParticipantID: pID,
		ProductID:     productID,
		UserID:        userClaims.UserID,
		Reason:        body.Reason,
	}

	result, err := ctrl.usecase.RejectParticipant(c.UserContext(), req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapParticipantResponse(result),
	})
}

func (ctrl *ParticipantController) SelfRegister(c *fiber.Ctx) error {
	claims, err := middleware.GetUserClaims(c)
	if err != nil {
		return participantError(c, err)
	}

	var req participantdto.SelfRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return participantError(c, errors.ErrBadRequest("invalid request body"))
	}

	if err := validate.Struct(&req); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	req.UserID = claims.UserID

	result, err := ctrl.usecase.SelfRegister(c.UserContext(), &req)
	if err != nil {
		return participantError(c, err)
	}

	status := fiber.StatusCreated
	if result.IsLinked {
		status = fiber.StatusOK
	}
	return c.Status(status).JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapSelfRegisterResponse(result),
	})
}

func (ctrl *ParticipantController) Delete(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		return participantError(c, err)
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	err = ctrl.usecase.DeleteParticipant(c.UserContext(), &participantdto.DeleteParticipantRequest{
		ParticipantID: pID,
		TenantID:      tenantID,
		ProductID:     productID,
		UserID:        userClaims.UserID,
	})
	if err != nil {
		return participantError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ctrl *ParticipantController) SaveAddresses(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		return participantError(c, err)
	}

	var req participantdto.SaveAddressesRequest
	if err := c.BodyParser(&req); err != nil {
		return participantError(c, errors.ErrBadRequest("invalid request body"))
	}

	if err := validate.Struct(&req); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req.TenantID = tenantID
	req.ProductID = productID
	req.ParticipantID = pID
	req.UserID = userClaims.UserID

	result, err := ctrl.usecase.SaveAddresses(c.UserContext(), &req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) SaveFamilyMembers(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		return participantError(c, err)
	}

	var req participantdto.SaveFamilyMembersRequest
	if err := c.BodyParser(&req); err != nil {
		return participantError(c, errors.ErrBadRequest("invalid request body"))
	}

	if err := validate.Struct(&req); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req.TenantID = tenantID
	req.ProductID = productID
	req.ParticipantID = pID
	req.UserID = userClaims.UserID

	result, err := ctrl.usecase.SaveFamilyMembers(c.UserContext(), &req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) SaveBeneficiaries(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return participantError(c, errors.ErrBadRequest("invalid participant ID"))
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		return participantError(c, err)
	}

	var req participantdto.SaveBeneficiariesRequest
	if err := c.BodyParser(&req); err != nil {
		return participantError(c, errors.ErrBadRequest("invalid request body"))
	}

	if err := validate.Struct(&req); err != nil {
		var ve validator.ValidationErrors
		if stderrors.As(err, &ve) {
			return participantError(c, errors.ErrValidationWithFields(convertValidationErrors(ve)))
		}
		return participantError(c, errors.ErrBadRequest("invalid request"))
	}

	productID, err := middleware.GetProductIDFromContext(c)
	if err != nil {
		return participantError(c, err)
	}

	req.TenantID = tenantID
	req.ProductID = productID
	req.ParticipantID = pID
	req.UserID = userClaims.UserID

	result, err := ctrl.usecase.SaveBeneficiaries(c.UserContext(), &req)
	if err != nil {
		return participantError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}
