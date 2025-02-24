package types

// ErrorResponse represents a 500 error on the server
// @Description an unknown error
type ErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}

// MissingFieldResponse represents an error response for a missing required field
//
//	@Description	an example of a missing field response
//	@Description	an example of a missing field response
type MissingFieldResponse struct {
	Error struct {
		Username string `json:"username" example:"This field is required"`
	} `json:"error"`
}

type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}
