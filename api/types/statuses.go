package types

// AddStatusRequest represents the request body for creating a status
// @Description A request body for adding a new status
type AddStatusRequest struct {
	StatusLabel string `json:"label" example:"test" binding:"required"`
}
