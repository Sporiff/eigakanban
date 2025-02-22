package types

import "github.com/golang-jwt/jwt/v5"

// AuthenticatedUserResponse represents the data passed back when logging in
type AuthenticatedUserResponse struct {
	Uuid         string
	AccessToken  string
	RefreshToken string
	SuperUser    bool
}

func (r *AuthenticatedUserResponse) Init(uuid, accessToken, refreshToken string, superUser bool) {
	r.Uuid = uuid
	r.AccessToken = accessToken
	r.RefreshToken = refreshToken
	r.SuperUser = superUser
}

// TokenClaims represents the claims stored in the JWT token
type TokenClaims struct {
	UserUuid  string // The user's UUID
	SuperUser bool   // Whether the user is a superuser
	ExpiresAt int64  // The expiry date of the token as a UNIX timestamp
	jwt.RegisteredClaims
}

// RegisterUserRequest represents the request body for registering a user
//
//	@Description	A request body for registering a new user
type RegisterUserRequest struct {
	Username string `json:"username" example:"test" binding:"required"`
	Email    string `json:"email" example:"test@test.com" binding:"required,email"`
	Password string `json:"password" example:"password" binding:"required"`
}

// LoginUserRequest represents the request body for logging in a user
//
//	@Description	request body for a login request.
//	@Description	either email or username must be provided
type LoginUserRequest struct {
	Email    string `json:"email" example:"test@test.com"`
	Username string `json:"username" example:"test"`
	Password string `json:"password" example:"password" binding:"required"`
}

// TokenResponse represents the request body for receiving access and refresh tokens
//
//	@Description	a response containing a JWT for authentication and a refresh token
type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"00000000-0000-0000-0000-000000000000"`
	ExpiryDate   string `json:"expiry_date" example:"2025-02-15 11:59:01.837871 +0100 CET m=+3603.614509085"`
	RefreshToken string `json:"refresh_token" example:"00000000-0000-0000-0000-000000000000"`
}

// LogoutSuccessResponse the response for a user logging out successfully
// @Description user logged out successfully
type LogoutSuccessResponse struct {
	Message string `json:"message" example:"logged out successfully"`
}

// AlreadyLoggedOutResponse the response for when a user has already been logged out
// @Description user already logged out
type AlreadyLoggedOutResponse struct {
	Error string `json:"error" example:"already logged out"`
}

// RefreshTokenMissingResponse the response for when a request is sent without a refresh token
// @Description refresh token missing
type RefreshTokenMissingResponse struct {
	Error string `json:"error" example:"refresh token missing"`
}
