package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/pkg/logger"
)

// SuccessResponse sends a successful API response
func SuccessResponse(c *gin.Context, data interface{}, message string) {
	response := models.SuccessResponse(data, message)
	c.JSON(http.StatusOK, response)
}

// CreatedResponse sends a created response
func CreatedResponse(c *gin.Context, data interface{}, message string) {
	response := models.SuccessResponse(data, message)
	c.JSON(http.StatusCreated, response)
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	response := models.ErrorResponseWithMessage("error", message, statusCode)
	c.JSON(statusCode, response)
}

// SendError sends an error response based on the error type
func SendError(c *gin.Context, err error) {
	var statusCode int
	var errorType, message string

	switch e := err.(type) {
	case *AppError:
		statusCode = e.StatusCode
		errorType = e.Err.Error()
		message = e.Message
	default:
		statusCode = http.StatusInternalServerError
		errorType = "internal_error"
		message = "An unexpected error occurred"
		logger.WithComponent("error-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Error("Unhandled error")
	}

	response := models.ErrorResponseWithMessage(errorType, message, statusCode)
	c.JSON(statusCode, response)
}

// SendValidationError sends a validation error response
func SendValidationError(c *gin.Context, message string) {
	response := models.ErrorResponseWithMessage("validation_error", message, http.StatusBadRequest)
	c.JSON(http.StatusBadRequest, response)
}
