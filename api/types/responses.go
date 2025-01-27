package types

// Pagination represents the generic pagination format used in GET requests
// @Description pagination information
type Pagination struct {
	Total    int64 `json:"total" example:"2"`
	Page     int32 `json:"page" example:"1"`
	PageSize int32 `json:"page_size" example:"50"`
}

// ErrorResponse represents a 500 error on the server
// @Description an unknown error
type ErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}

// MessageResponse represents a basic success message
// @Description a generic success message
type MessageResponse struct {
	Message string `json:"message" example:"success"`
}
