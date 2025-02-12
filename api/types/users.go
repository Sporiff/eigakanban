package types

// UserResponse represents the response returned by the API
//
//	@Description	JSON representation of a user in the system
type UserResponse struct {
	UUID     string `json:"uuid" example:"77b62cff-0020-43d9-a90c-5d35bff89f7a"`
	Username string `json:"username" example:"username"`
	FullName string `json:"full_name" example:"Tim Test"`
	Bio      string `json:"bio" example:"This is a bio"`
}

// PaginatedUsersResponse represents a response containing a list of users
//
//	@Description	a response containing a list of users and a pagination object
type PaginatedUsersResponse struct {
	Pagination Pagination     `json:"pagination"`
	Users      []UserResponse `json:"users"`
}

// UpdateUserRequest represents the request body for updating a user
//
//	@Description	a request body for updating a user
type UpdateUserRequest struct {
	NewUsername *string `json:"username" example:"new_username"`
	NewName     *string `json:"full_name" example:"Tim Test"`
	NewBio      *string `json:"bio" example:"This is a bio"`
}

// UserDeletedResponse represents a success message for a user deletion
//
//	@Description	A success message confirming the user was deleted
type UserDeletedResponse struct {
	Message string `json:"success" example:"user deleted: 77b62cff-0020-43d9-a90c-5d35bff89f7a"`
}

// UserNotFoundResponse represents an error message for when a user isn't found
//
// @Description User not found
type UserNotFoundResponse struct {
	Error string `json:"error" example:"user not found"`
}
