package types

// AddStatusRequest represents the request body for creating a status
// @Description A request body for adding a new status
type AddStatusRequest struct {
	StatusLabel string `json:"label" example:"test" binding:"required"`
}

// StatusesResponse represents a status
// @Description status details
type StatusesResponse struct {
	UUID  string `json:"uuid" example:"00000000-0000-0000-0000-000000000000"`
	Label string `json:"label" example:"backlog"`
}

// PaginatedStatusesResponse represents a paginated list of statuses
// @Description a paginated list of statuses
type PaginatedStatusesResponse struct {
	Statuses   []StatusesResponse `json:"statuses"`
	Pagination Pagination         `json:"pagination"`
}
