package errors

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

type Kind int

const (
	KindUnexpected Kind = iota
	KindNotFound
	KindDuplicate
	KindValidation
	KindConnection
	KindTimeout
	KindPermission
	KindBadRequest
	KindUnauthorized
	KindForbidden
)

func (k Kind) String() string {
	switch k {
	case KindNotFound:
		return "not_found"
	case KindDuplicate:
		return "duplicate"
	case KindValidation:
		return "validation"
	case KindConnection:
		return "connection"
	case KindTimeout:
		return "timeout"
	case KindPermission:
		return "permission"
	case KindBadRequest:
		return "bad_request"
	case KindUnauthorized:
		return "unauthorized"
	case KindForbidden:
		return "forbidden"
	default:
		return "unexpected"
	}
}

const (
	CodeInternal           = "ERR_INTERNAL"
	CodeValidation         = "ERR_VALIDATION"
	CodeNotFound           = "ERR_NOT_FOUND"
	CodeConflict           = "ERR_CONFLICT"
	CodeBadRequest         = "ERR_BAD_REQUEST"
	CodeUnauthorized       = "ERR_UNAUTHORIZED"
	CodeForbidden          = "ERR_FORBIDDEN"
	CodeTooManyRequests    = "ERR_TOO_MANY_REQUESTS"
	CodeServiceUnavailable = "ERR_SERVICE_UNAVAILABLE"
	CodeUnprocessable      = "ERR_UNPROCESSABLE"

	CodeInvalidCredentials = "ERR_INVALID_CREDENTIALS"
	CodeTokenExpired       = "ERR_TOKEN_EXPIRED"
	CodeTokenInvalid       = "ERR_TOKEN_INVALID"
	CodeSessionExpired     = "ERR_SESSION_EXPIRED"
	CodeOTPInvalid         = "ERR_OTP_INVALID"
	CodeOTPExpired         = "ERR_OTP_EXPIRED"
	CodePINInvalid         = "ERR_PIN_INVALID"
	CodePINLocked          = "ERR_PIN_LOCKED"
	CodePINRequired        = "ERR_PIN_REQUIRED"

	CodeUserNotFound      = "ERR_USER_NOT_FOUND"
	CodeUserAlreadyExists = "ERR_USER_ALREADY_EXISTS"
	CodeUserNotApproved   = "ERR_USER_NOT_APPROVED"
	CodeUserSuspended     = "ERR_USER_SUSPENDED"
	CodeUserInactive      = "ERR_USER_INACTIVE"
	CodeEmailNotVerified  = "ERR_EMAIL_NOT_VERIFIED"
	CodeProfileIncomplete = "ERR_PROFILE_INCOMPLETE"

	CodeTenantNotFound  = "ERR_TENANT_NOT_FOUND"
	CodeTenantInactive  = "ERR_TENANT_INACTIVE"
	CodeTenantSuspended = "ERR_TENANT_SUSPENDED"
	CodeInvalidTenant   = "ERR_INVALID_TENANT"

	CodePermissionDenied      = "ERR_PERMISSION_DENIED"
	CodeRoleNotFound          = "ERR_ROLE_NOT_FOUND"
	CodeInvalidPermission     = "ERR_INVALID_PERMISSION"
	CodeAccessForbidden       = "ERR_ACCESS_FORBIDDEN"
	CodePlatformAdminRequired = "ERR_PLATFORM_ADMIN_REQUIRED"

	CodeEmployeeNotFound = "ERR_EMPLOYEE_NOT_FOUND"
	CodeEmployeeExists   = "ERR_EMPLOYEE_EXISTS"
	CodeInvalidNIK       = "ERR_INVALID_NIK"

	CodeContributionNotFound    = "ERR_CONTRIBUTION_NOT_FOUND"
	CodeContributionInvalid     = "ERR_CONTRIBUTION_INVALID"
	CodeContributionDuplicate   = "ERR_CONTRIBUTION_DUPLICATE"
	CodeInvalidContributionData = "ERR_INVALID_CONTRIBUTION_DATA"

	CodeAllocationNotFound = "ERR_ALLOCATION_NOT_FOUND"
	CodeAllocationInvalid  = "ERR_ALLOCATION_INVALID"
	CodeAllocationMismatch = "ERR_ALLOCATION_MISMATCH"
	CodeInvalidProportion  = "ERR_INVALID_PROPORTION"

	CodeFileNotFound      = "ERR_FILE_NOT_FOUND"
	CodeFileInvalid       = "ERR_FILE_INVALID"
	CodeFileTooLarge      = "ERR_FILE_TOO_LARGE"
	CodeUnsupportedFormat = "ERR_UNSUPPORTED_FORMAT"

	CodeDatabaseError     = "ERR_DATABASE"
	CodeTransactionFailed = "ERR_TRANSACTION_FAILED"
	CodeDuplicateEntry    = "ERR_DUPLICATE_ENTRY"
)

type AppError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`

	HTTPStatus int `json:"-"`

	Op   string `json:"-"`
	Kind Kind   `json:"-"`
	File string `json:"-"`
	Line int    `json:"-"`
	Err  error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithOp(op string) *AppError {
	e.Op = op
	return e
}

func (e *AppError) WithKind(kind Kind) *AppError {
	e.Kind = kind
	return e
}

func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

func newWithCaller(code, message string, httpStatus int, kind Kind, skip int) *AppError {
	file, line := "unknown", 0
	if _, f, l, ok := runtime.Caller(skip); ok {
		if idx := strings.LastIndex(f, "/"); idx >= 0 {
			if idx2 := strings.LastIndex(f[:idx], "/"); idx2 >= 0 {
				f = f[idx2+1:]
			} else {
				f = f[idx+1:]
			}
		}
		file, line = f, l
	}
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Kind:       kind,
		File:       file,
		Line:       line,
	}
}

func New(code, message string, httpStatus int) *AppError {
	return newWithCaller(code, message, httpStatus, KindUnexpected, 2)
}

func Wrap(err error, code, message string, httpStatus int) *AppError {
	appErr := newWithCaller(code, message, httpStatus, KindUnexpected, 2)
	appErr.Err = err
	return appErr
}

func ErrInternal(message string) *AppError {
	return newWithCaller(CodeInternal, message, http.StatusInternalServerError, KindUnexpected, 2)
}

func ErrBadRequest(message string) *AppError {
	return newWithCaller(CodeBadRequest, message, http.StatusBadRequest, KindBadRequest, 2)
}

func ErrValidation(message string) *AppError {
	return newWithCaller(CodeValidation, message, http.StatusBadRequest, KindValidation, 2)
}

func ErrNotFound(message string) *AppError {
	return newWithCaller(CodeNotFound, message, http.StatusNotFound, KindNotFound, 2)
}

func ErrConflict(message string) *AppError {
	return newWithCaller(CodeConflict, message, http.StatusConflict, KindDuplicate, 2)
}

func ErrUnauthorized(message string) *AppError {
	return newWithCaller(CodeUnauthorized, message, http.StatusUnauthorized, KindUnauthorized, 2)
}

func ErrForbidden(message string) *AppError {
	return newWithCaller(CodeForbidden, message, http.StatusForbidden, KindForbidden, 2)
}

func ErrTooManyRequests(message string) *AppError {
	return newWithCaller(CodeTooManyRequests, message, http.StatusTooManyRequests, KindBadRequest, 2)
}

func ErrInvalidCredentials() *AppError {
	return newWithCaller(CodeInvalidCredentials, "Invalid email or password", http.StatusUnauthorized, KindUnauthorized, 2)
}

func ErrTokenExpired() *AppError {
	return newWithCaller(CodeTokenExpired, "Token has expired", http.StatusUnauthorized, KindUnauthorized, 2)
}

func ErrTokenInvalid() *AppError {
	return newWithCaller(CodeTokenInvalid, "Invalid token", http.StatusUnauthorized, KindUnauthorized, 2)
}

func ErrOTPInvalid() *AppError {
	return newWithCaller(CodeOTPInvalid, "Invalid OTP code", http.StatusBadRequest, KindValidation, 2)
}

func ErrOTPExpired() *AppError {
	return newWithCaller(CodeOTPExpired, "OTP code has expired", http.StatusBadRequest, KindValidation, 2)
}

func ErrPINInvalid() *AppError {
	return newWithCaller(CodePINInvalid, "Invalid PIN", http.StatusBadRequest, KindValidation, 2)
}

func ErrPINLocked() *AppError {
	return newWithCaller(CodePINLocked, "PIN is locked due to too many failed attempts", http.StatusForbidden, KindForbidden, 2)
}

func ErrPINRequired() *AppError {
	return newWithCaller(CodePINRequired, "PIN verification required for this operation", http.StatusForbidden, KindForbidden, 2)
}

func ErrUserNotFound() *AppError {
	return newWithCaller(CodeUserNotFound, "User not found", http.StatusNotFound, KindNotFound, 2)
}

func ErrUserAlreadyExists() *AppError {
	return newWithCaller(CodeUserAlreadyExists, "User with this email already exists", http.StatusConflict, KindDuplicate, 2)
}

func ErrUserNotApproved() *AppError {
	return newWithCaller(CodeUserNotApproved, "User registration is pending approval", http.StatusForbidden, KindForbidden, 2)
}

func ErrUserSuspended() *AppError {
	return newWithCaller(CodeUserSuspended, "User account is suspended", http.StatusForbidden, KindForbidden, 2)
}

func ErrProfileIncomplete() *AppError {
	return newWithCaller(CodeProfileIncomplete, "User profile is incomplete", http.StatusForbidden, KindForbidden, 2)
}

func ErrTenantNotFound() *AppError {
	return newWithCaller(CodeTenantNotFound, "Tenant not found", http.StatusNotFound, KindNotFound, 2)
}

func ErrTenantInactive() *AppError {
	return newWithCaller(CodeTenantInactive, "Tenant is inactive", http.StatusForbidden, KindForbidden, 2)
}

func ErrPermissionDenied() *AppError {
	return newWithCaller(CodePermissionDenied, "You do not have permission to perform this action", http.StatusForbidden, KindPermission, 2)
}

func ErrRoleNotFound() *AppError {
	return newWithCaller(CodeRoleNotFound, "Role not found", http.StatusNotFound, KindNotFound, 2)
}

func ErrAccessForbidden(message string) *AppError {
	return newWithCaller(CodeAccessForbidden, message, http.StatusForbidden, KindForbidden, 2)
}

func ErrPlatformAdminRequired() *AppError {
	return newWithCaller(CodePlatformAdminRequired, "This operation requires platform administrator privileges", http.StatusForbidden, KindPermission, 2)
}

func ErrEmployeeNotFound() *AppError {
	return newWithCaller(CodeEmployeeNotFound, "Employee not found", http.StatusNotFound, KindNotFound, 2)
}

func ErrEmployeeExists() *AppError {
	return newWithCaller(CodeEmployeeExists, "Employee with this NIK already exists", http.StatusConflict, KindDuplicate, 2)
}

func ErrContributionNotFound() *AppError {
	return newWithCaller(CodeContributionNotFound, "Contribution not found", http.StatusNotFound, KindNotFound, 2)
}

func ErrAllocationNotFound() *AppError {
	return newWithCaller(CodeAllocationNotFound, "Allocation not found", http.StatusNotFound, KindNotFound, 2)
}

func ErrInvalidProportion() *AppError {
	return newWithCaller(CodeInvalidProportion, "Allocation proportions must sum to 100%", http.StatusBadRequest, KindValidation, 2)
}

func ErrFileNotFound() *AppError {
	return newWithCaller(CodeFileNotFound, "File not found", http.StatusNotFound, KindNotFound, 2)
}

func ErrFileTooLarge(maxSize string) *AppError {
	return newWithCaller(CodeFileTooLarge, fmt.Sprintf("File exceeds maximum size of %s", maxSize), http.StatusBadRequest, KindValidation, 2)
}

func ErrUnsupportedFormat(format string) *AppError {
	return newWithCaller(CodeUnsupportedFormat, fmt.Sprintf("Unsupported file format: %s", format), http.StatusBadRequest, KindValidation, 2)
}

func ErrDatabase(message string) *AppError {
	return newWithCaller(CodeDatabaseError, message, http.StatusInternalServerError, KindConnection, 2)
}

func ErrUnprocessable(message string) *AppError {
	return newWithCaller(CodeUnprocessable, message, http.StatusUnprocessableEntity, KindValidation, 2)
}

func ErrDuplicateEntry(field string) *AppError {
	return newWithCaller(CodeDuplicateEntry, fmt.Sprintf("Duplicate entry for %s", field), http.StatusConflict, KindDuplicate, 2)
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ErrValidationWithFields(fields []FieldError) *AppError {
	err := newWithCaller(CodeValidation, "Validation failed", http.StatusBadRequest, KindValidation, 2)
	err.Details = map[string]interface{}{
		"fields": fields,
	}
	return err
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

func GetHTTPStatus(err error) int {
	if appErr := GetAppError(err); appErr != nil {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}

func GetCode(err error) string {
	if appErr := GetAppError(err); appErr != nil {
		return appErr.Code
	}
	return CodeInternal
}

func IsNotFound(err error) bool {
	if appErr := GetAppError(err); appErr != nil {
		return appErr.Kind == KindNotFound
	}
	return false
}

func IsConflict(err error) bool {
	if appErr := GetAppError(err); appErr != nil {
		return appErr.Kind == KindDuplicate
	}
	return false
}

func IsValidation(err error) bool {
	if appErr := GetAppError(err); appErr != nil {
		return appErr.Kind == KindValidation
	}
	return false
}

func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	if appErr := GetAppError(err); appErr != nil {
		return appErr.Kind == KindDuplicate
	}
	return false
}
