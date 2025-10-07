package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates a struct using the validator package
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// GetValidationErrors extracts validation error messages
func GetValidationErrors(err error) []string {
	var errors []string
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			errors = append(errors, getErrorMessage(e))
		}
	}
	return errors
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters long"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters long"
	case "oneof":
		return fe.Field() + " must be one of: " + fe.Param()
	default:
		return fe.Field() + " is invalid"
	}
}

// BindAndValidate binds JSON request and validates it
func BindAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return NewBadRequestError("Invalid JSON format")
	}

	if err := ValidateStruct(obj); err != nil {
		validationErrors := GetValidationErrors(err)
		if len(validationErrors) > 0 {
			return NewBadRequestError(validationErrors[0])
		}
		return NewBadRequestError("Validation failed")
	}

	return nil
}
