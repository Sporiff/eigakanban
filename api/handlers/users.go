package handlers

import (
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/services"
	"codeberg.org/sporiff/eigakanban/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UsersHandler struct {
	usersService *services.UsersService
}

func NewUsersHandler(usersService *services.UsersService) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
	}
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
//	@Success		200			{object}	types.PaginatedUsersResponse
//	@Failure		500			{object}	types.ErrorResponse
//	@Router			/users [get]
func (h *UsersHandler) GetAllUsers(c *gin.Context) {
	pagination, err := helpers.ValidatePagination(c)
	if err != nil {
		helpers.HandleAPIError(c, err)
		return
	}

	result, err := h.usersService.GetAllUsers(c.Request.Context(), pagination)
	if err != nil {
		helpers.HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
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
//	@Success		200		{object}	types.UserResponse
//	@Failure		400		{object}	types.ErrorResponse
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/users/{uuid} [get]
func (h *UsersHandler) GetUserByUuid(c *gin.Context) {
	userUuid := c.Param("uuid")
	user, err := h.usersService.GetUserByUuid(c.Request.Context(), userUuid)
	if err != nil {
		helpers.HandleAPIError(c, err)
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
//	@Param			uuid	path		string					true	"User UUID"
//	@Param			body	body		types.UpdateUserRequest	true	"User details to update"
//	@Success		200		{object}	types.UserResponse
//	@Failure		400		{object}	types.ErrorResponse
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/users/{uuid} [patch]
func (h *UsersHandler) UpdateUser(c *gin.Context) {
	userUuid := c.Param("uuid")
	var req types.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	user, err := h.usersService.UpdateUser(c.Request.Context(), userUuid, req)
	if err != nil {
		helpers.HandleAPIError(c, err)
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
//	@Param			uuid	path		string						true	"User UUID"
//	@Success		200		{object}	types.UserDeletedResponse	"User deleted successfully"
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/users/{uuid} [delete]
func (h *UsersHandler) DeleteUser(c *gin.Context) {
	userUuid := c.Param("uuid")
	err := h.usersService.DeleteUser(c.Request.Context(), userUuid)
	if err != nil {
		helpers.HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "user deleted: " + userUuid})
}
