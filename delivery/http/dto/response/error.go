package response

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(errorCode string, message string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   errorCode,
		Message: message,
	}
}

func ErrorResponseWithDetails(errorCode string, message string, details interface{}) APIResponse {
	return APIResponse{
		Success: false,
		Error:   errorCode,
		Message: message,
		Data:    details,
	}
}
