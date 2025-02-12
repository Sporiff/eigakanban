package handlers

import (
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/services"
	"codeberg.org/sporiff/eigakanban/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterUser adds a new user to the system
//
//	@Summary		Register a new user account
//	@Description	Register a new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		types.RegisterUserRequest	true	"User details"
//	@Success		200		{object}	types.UserResponse			"User registered successfully"
//	@Failure		400		{object}	types.MissingFieldResponse	"Missing mandatory fields"
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/auth/register [post]
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req types.RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	user, err := h.authService.RegisterUser(c.Request.Context(), req)
	if err != nil {
		switch err.Error() {
		case "user already exists":
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		case "error creating default data":
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// LoginUser logs in a user
//
//	@Summary		Log in
//	@Description	Log in to user account using email or username
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		types.LoginUserRequest		true	"Login details"
//	@Success		200		{object}	types.TokenResponse			"Successful login"
//	@Failure		400		{object}	types.MissingFieldResponse	"Missing mandatory fields"
//	@Failure		404		{object}	types.UserNotFoundResponse	"User not found"
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/auth/login [post]
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var req types.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	user, err := h.authService.LoginUser(c.Request.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case "invalid password":
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		case "email and username are required":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.Set("user_uuid", user.Uuid)

	// Return both tokens in the response
	c.JSON(http.StatusOK, gin.H{
		"access_token":  user.AccessToken,
		"refresh_token": user.RefreshToken,
	})
}

// LogoutUser logs out the current user by clearing their refresh token
//
//	@Summary		Log out
//	@Description	Log out of the app
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.LogoutSuccessResponse			"Logout successful"
//	@Failure		400	{object}	types.AlreadyLoggedOutResponse		"Already logged out"
//	@Failure		400	{object}	types.RefreshTokenMissingResponse	"Missing refresh token"
//	@Failure		500	{object}	types.ErrorResponse
//	@Router			/auth/logout [post]
func (h *AuthHandler) LogoutUser(c *gin.Context) {
	// TODO: Investigate how to handle this in both headers and cookies
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already logged out"})
		return
	}

	err := h.authService.LogoutUser(c.Request.Context(), refreshToken)
	if err != nil {
		switch err.Error() {
		case "already logged out":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
