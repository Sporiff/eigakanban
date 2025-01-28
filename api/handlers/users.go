package handlers

import (
	"codeberg.org/sporiff/eigakanban/helpers"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"strconv"

	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/types"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandler struct {
	db *pgxpool.Pool
	q  *queries.Queries
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{
		db: db,
		q:  queries.New(db),
	}
}

// UserResponse represents the response returned by the API
//
//	@Description	JSON representation of a user in the system
type UserResponse struct {
	UUID     string `json:"uuid" example:"77b62cff-0020-43d9-a90c-5d35bff89f7a"`
	Username string `json:"username" example:"username"`
	FullName string `json:"full_name" example:"Tim Test"`
	Bio      string `json:"bio" example:"This is a bio"`
}

// GetAllUsers returns a paginated array of users
//
//	@Summary		Get all users
//	@Description	Get all users in a paginated list
//	@Security		BearerAuth
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"Page"
//	@Param			page_size	query		int	false	"Page size"
//	@Success		200			{object}	handlers.GetAllUsers.PaginatedUsersResponse
//	@Failure		500			{object}	types.ErrorResponse
//	@Router			/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// PaginatedUsersResponse represents a response containing a list of users
	//	@Description	a response containing a list of users and a pagination object
	type PaginatedUsersResponse struct {
		Pagination types.Pagination `json:"pagination"`
		Users      []UserResponse   `json:"users"`
	}

	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page must be an integer"})
		return
	}

	pageSize, err := strconv.ParseInt(c.DefaultQuery("page_size", "50"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page_size must be an integer"})
		return
	}

	total, err := h.q.GetUserCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	users, err := h.q.GetAllUsers(c.Request.Context(), queries.GetAllUsersParams{
		Page:     int32(page - 1),
		PageSize: int32(pageSize),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	paginationValues := types.Pagination{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Total:    total,
	}

	c.JSON(http.StatusOK, gin.H{"pagination": paginationValues, "users": users})
}

// GetUserByUuid returns a user by UUID
//
//	@Summary		Get user by UUID
//	@Description	Get a user by UUID
//	@Security		BearerAuth
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string	true	"User UUID"
//	@Success		200		{object}	UserResponse
//	@Failure		400		{object}	types.ErrorResponse
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/users/{uuid} [get]
func (h *UserHandler) GetUserByUuid(c *gin.Context) {
	userUuid := c.Param("uuid")
	pgUuid, err := helpers.ValidateAndConvertUUID(userUuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.q.GetUserByUuid(c.Request.Context(), pgUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser updates user details
//
//	@Summary		Update user details
//	@Description	Update user details by UUID
//	@Security		BearerAuth
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string									true	"User UUID"
//	@Param			body	body		handlers.UpdateUser.UpdateUserRequest	true	"User details to update"
//	@Success		200		{object}	UserResponse
//	@Failure		400		{object}	types.ErrorResponse
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/users/{uuid} [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// UpdateUserRequest represents the request body for updating a user
	//	@Description	a request body for updating a user
	type UpdateUserRequest struct {
		NewUsername *string `json:"username" example:"new_username"`
		NewName     *string `json:"full_name" example:"Tim Test"`
		NewBio      *string `json:"bio" example:"This is a bio"`
	}

	userUuid := c.Param("uuid")
	pgUuid, err := helpers.ValidateAndConvertUUID(userUuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	params := queries.UpdateUserDetailsParams{
		UserUuid: pgUuid,
	}

	// If no username is provided, we need to fetch the current username so it doesn't get overwritten
	if req.NewUsername != nil {
		params.NewUsername = *req.NewUsername
	} else {
		currentUser, err := h.q.GetUserByUuid(c.Request.Context(), pgUuid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch current user details"})
			return
		}
		params.NewUsername = currentUser.Username
	}

	// Ensure that any missing fields don't overwrite the current field
	if req.NewName != nil {
		params.NewName = pgtype.Text{String: *req.NewName, Valid: true}
	} else {
		params.NewName = pgtype.Text{Valid: false}
	}

	if req.NewBio != nil {
		params.NewBio = pgtype.Text{String: *req.NewBio, Valid: true}
	} else {
		params.NewBio = pgtype.Text{Valid: false}
	}

	// Update the user details
	user, err := h.q.UpdateUserDetails(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// DeleteUser deletes a user from the database by UUID
//
//	@Summary		Delete user
//	@Description	Delete a user by UUID
//	@Security		BearerAuth
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string									true	"User UUID"
//	@Success		200		{object}	handlers.DeleteUser.UserDeletedResponse	"User deleted successfully"
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/users/{uuid} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// UserDeletedResponse represents a success message for a user deletion
	//	@Description	A success message confirming the user was deleted
	type UserDeletedResponse struct {
		Message string `json:"success" example:"user deleted: 77b62cff-0020-43d9-a90c-5d35bff89f7a"`
	}

	userUuid := c.Param("uuid")
	pgUuid, err := helpers.ValidateAndConvertUUID(userUuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.q.DeleteUser(c.Request.Context(), pgUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseString := UserDeletedResponse{
		Message: fmt.Sprintf("user deleted: %s", userUuid),
	}

	c.JSON(http.StatusOK, gin.H{"result": responseString.Message})
}
