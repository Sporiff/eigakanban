package helpers

import (
	"codeberg.org/sporiff/eigakanban/types"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

// FormatValidationError formats validation errors into a more meaningful message
func FormatValidationError(err error) map[string]string {
	errorMessages := make(map[string]string)

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			// Format the field name to be more user-friendly
			fieldName := strings.ToLower(fieldError.Field())

			// Return a custom message based on the field tag
			switch fieldError.Tag() {
			case "required":
				errorMessages[fieldName] = "This field is required"
			case "email":
				errorMessages[fieldName] = "Invalid email format"
			case "min":
				errorMessages[fieldName] = "This field must be at least " + fieldError.Param() + " characters"
			case "max":
				errorMessages[fieldName] = "This field must be at most " + fieldError.Param() + " characters"
			default:
				errorMessages[fieldName] = "Invalid value for " + fieldName
			}
		}
	}

	return errorMessages
}

// HandleValidationError returns a custom error response for validation errors
func HandleValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": FormatValidationError(err),
	})
}

func HandleAPIError(c *gin.Context, err error) {
	var apiErr *types.APIError
	if errors.As(err, &apiErr) {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error occurred"})
	}
}
