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
		helpers.HandleAPIError(c, err)
		return
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
		helpers.HandleAPIError(c, err)
		return
	}

	c.Set("user_uuid", user.Uuid)
	c.Set("superuser", user.SuperUser)

	response := types.TokenResponse{
		AccessToken:  user.AccessToken,
		ExpiryDate:   user.ExpiryDate,
		RefreshToken: user.RefreshToken,
	}

	// Return both tokens in the response
	c.JSON(http.StatusOK, response)
}

// RefreshToken refreshes a user's access token by passing in a valid refresh token
//
//	@Summary		Log out
//	@Description	Log out of the app
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			Refresh-Token	header		string								true	"Refresh token"
//	@Success		200				{object}	types.AccessTokenResponse			"New access token"
//	@Failure		400				{object}	types.RefreshTokenMissingResponse	"Missing refresh token"
//	@Failure		500				{object}	types.ErrorResponse
//	@Router			/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := helpers.GetRefreshToken(c)
	if err != nil {
		helpers.HandleAPIError(c, err)
		return
	}

	tokenResponse, err := h.authService.CreateNewAccessToken(c, refreshToken)
	if err != nil {
		helpers.HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

// LogoutUser logs out the current user by clearing their refresh token
//
//	@Summary		Log out
//	@Description	Log out of the app
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			Refresh-Token	header		string								true	"Refresh token"
//	@Success		200				{object}	types.LogoutSuccessResponse			"Logout successful"
//	@Failure		400				{object}	types.AlreadyLoggedOutResponse		"Already logged out"
//	@Failure		400				{object}	types.RefreshTokenMissingResponse	"Missing refresh token"
//	@Failure		500				{object}	types.ErrorResponse
//	@Router			/auth/logout [post]
func (h *AuthHandler) LogoutUser(c *gin.Context) {
	refreshToken, err := helpers.GetRefreshToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing refresh token"})
		return
	}

	err = h.authService.LogoutUser(c.Request.Context(), *refreshToken)
	if err != nil {
		helpers.HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
