package handlers

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/types"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
)

type AuthHandler struct {
	db *pgxpool.Pool
	q  *queries.Queries
}

func NewAuthHandler(db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{
		db: db,
		q:  queries.New(db),
	}
}

// MissingFieldResponse represents an error response for a missing required field
//
//	@Description	an example of a missing field respo
//	@Description	an example of a missing field response
type MissingFieldResponse struct {
	Error struct {
		Username string `json:"username" example:"This field is required"`
	} `json:"error"`
}

// RegisterUser adds a new user to the system
//
//	@Summary		Register a new user account
//	@Description	Register a new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		handlers.RegisterUser.RegisterUserRequest	true	"User details"
//	@Success		200		{object}	handlers.UserResponse						"User registered successfully"
//	@Failure		400		{object}	handlers.MissingFieldResponse				"Missing mandatory fields"
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/auth/register [post]
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	// RegisterUserRequest represents the request body for registering a user
	//	@Description	A request body for registering a new user
	type RegisterUserRequest struct {
		Username string `json:"username" example:"test" binding:"required"`
		Email    string `json:"email" example:"test@test.com" binding:"required,email"`
		Password string `json:"password" example:"password" binding:"required"`
	}

	var req RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	params := queries.CheckForUserParams{
		Username: req.Username,
		Email:    req.Email,
	}

	existingUserCount, err := h.q.CheckForUser(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingUserCount > 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user already exists"})
		return
	}

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.q.AddUser(c.Request.Context(), queries.AddUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
//	@Param			body	body		handlers.LoginUser.LoginUserRequest		true	"Login details"
//	@Success		200		{object}	handlers.LoginUser.TokenResponse		"Successful login"
//	@Failure		400		{object}	handlers.MissingFieldResponse			"Missing mandatory fields"
//	@Failure		404		{object}	handlers.LoginUser.NoUserFoundResponse	"User not found"
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/auth/login [post]
func (h *AuthHandler) LoginUser(c *gin.Context) {
	// LoginUserRequest represents the request body for logging in a user
	//	@Description	request body for a login request.
	//	@Description	either email or username must be provided
	type LoginUserRequest struct {
		Email    string `json:"email" example:"test@test.com"`
		Username string `json:"username" example:"test"`
		Password string `json:"password" example:"password" binding:"required"`
	}

	// TokenResponse represents the response containing the JWT token
	//	@Description	a response containing a JWT for authentication
	type TokenResponse struct {
		Token string `json:"token" example:"jwt-token-string"`
	}

	type MissingFieldResponse struct {
		Error string `json:"error" example:"This field is required"`
	}

	type NoUserFoundResponse struct {
		Error string `json:"error" example:"User not found"`
	}

	var req LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	// Ensure either email or username is provided
	if req.Email == "" && req.Username == "" {
		helpers.HandleValidationError(c, errors.New("email and username are required"))
		return
	}

	// Prepare query parameters
	params := queries.GetExistingUserParams{
		Email:    req.Email,
		Username: req.Username,
	}

	// Fetch the user from the database
	existingUser, err := h.q.GetExistingUser(c.Request.Context(), params)
	if err != nil {
		response := types.ErrorResponse{Error: err.Error()}
		c.JSON(http.StatusInternalServerError, gin.H{"error": response.Error})
		return
	}

	// Check if the user exists
	if (queries.GetExistingUserRow{}) == existingUser {
		c.JSON(http.StatusNotFound, NoUserFoundResponse{Error: "User not found"})
		return
	}

	// Verify the password matches the stored password
	if !helpers.CheckPasswordHash(req.Password, existingUser.HashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate a JWT token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_uuid": existingUser.Uuid,
		"exp":       time.Now().Add(time.Hour * 1).Unix(),
	})

	// TODO: Look into getting users to set their own key
	accessTokenString, err := accessToken.SignedString([]byte("your-secret-key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	// Generate a refresh token
	refreshToken, err := helpers.GenerateRefreshToken(64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
	}

	// Set the refresh token to expire in 7 days
	refreshTokenExpiry := time.Now().Add(time.Hour * 24 * 7)

	RefreshTokenParams := queries.AddRefreshTokenParams{
		UserID:    existingUser.UserID.Int64,
		Token:     refreshToken,
		ExpiresAt: pgtype.Timestamptz{Time: refreshTokenExpiry, Valid: true},
	}

	// Store the refresh token in the database
	_, err = h.q.AddRefreshToken(c.Request.Context(), RefreshTokenParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store refresh token"})
		return
	}

	c.Set("user_uuid", existingUser.Uuid.String())

	// Return both tokens in the response
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessTokenString,
		"refresh_token": refreshToken,
	})
}

// LogoutUser logs out the current user by clearing their refresh token
//
//	@Summary		Log out
//	@Description	Log out of the app
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handlers.LogoutUser.LogoutSuccessResponse		"Logout successful"
//	@Failure		400	{object}	handlers.LogoutUser.AlreadyLoggedOutResponse	"Already logged out"
//	@Failure		400	{object}	handlers.LogoutUser.RefreshTokenMissingResponse	"Missing refresh token"
//	@Failure		500	{object}	types.ErrorResponse
//	@Router			/auth/logout [post]
func (h *AuthHandler) LogoutUser(c *gin.Context) {
	type RefreshTokenMissingResponse struct {
		Error string `json:"error" example:"Refresh token is required"`
	}

	type AlreadyLoggedOutResponse struct {
		Error string `json:"error" example:"Already logged out"`
	}

	type LogoutSuccessResponse struct {
		Message string `json:"message" example:"Logged out successfully"`
	}

	// TODO: Investigate how to handle this in both headers and cookies
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		response := RefreshTokenMissingResponse{
			Error: "Refresh token missing",
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
		return
	}

	_, err := h.q.GetRefreshTokenByToken(c.Request.Context(), refreshToken)
	if err != nil {
		response := AlreadyLoggedOutResponse{
			Error: "Already logged out",
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
		return
	}

	err = h.q.DeleteRefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		response := types.ErrorResponse{
			Error: "Failed to logout",
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
