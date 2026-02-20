package response

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
type DebugInfo struct {
	Cause string `json:"cause,omitempty"`
	Stack string `json:"stack,omitempty"`
}
type APIResponse struct {
	Success    bool         `json:"success"`
	Message    string       `json:"message,omitempty"`
	Data       interface{}  `json:"data,omitempty"`
	Error      string       `json:"error,omitempty"`
	Errors     []FieldError `json:"errors,omitempty"`
	Pagination *Pagination  `json:"pagination,omitempty"`
	RequestID  string       `json:"request_id,omitempty"`
	Debug      *DebugInfo   `json:"debug,omitempty"`
}
