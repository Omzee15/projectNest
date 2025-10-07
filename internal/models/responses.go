package models

// Standard API response structure
type APIResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
}

type ErrorResponse struct {
	Error      string `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

// Helper functions for creating responses
func SuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Data:    data,
		Success: true,
		Message: message,
	}
}

func SuccessResponseWithoutMessage(data interface{}) APIResponse {
	return APIResponse{
		Data:    data,
		Success: true,
	}
}

func ErrorResponseWithMessage(error, message string, statusCode int) ErrorResponse {
	return ErrorResponse{
		Error:      error,
		Message:    message,
		StatusCode: statusCode,
	}
}
