package services

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/types"
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"time"
)

type AuthService struct {
	q *queries.Queries
}

func NewAuthService(q *queries.Queries) *AuthService {
	return &AuthService{q: q}
}

// RegisterUser creates a new user and populates default information
func (s *AuthService) RegisterUser(ctx context.Context, user types.RegisterUserRequest) (*queries.AddUserRow, error) {
	// Check for a user with a matching email/username
	userCount, err := s.q.CheckForUser(ctx, queries.CheckForUserParams{
		Email:    user.Email,
		Username: user.Username,
	})
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "failed to check for user: "+err.Error())
	}

	if userCount != 0 {
		return nil, types.NewAPIError(http.StatusBadRequest, "user already exists")
	}

	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return nil, types.NewAPIError(http.StatusBadRequest, "invalid password")
	}

	// Add the user to the database
	registeredUser, err := s.q.AddUser(ctx, queries.AddUserParams{
		Username:       user.Username,
		HashedPassword: hashedPassword,
		Email:          user.Email,
	})
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "failed to add user: "+err.Error())
	}

	// Create default data for the user
	err = s.createDefaultData(ctx, registeredUser)
	if err != nil {
		return nil, err
	}

	return &registeredUser, nil
}

// LoginUser logs in the user and sets up authentication
func (s *AuthService) LoginUser(ctx context.Context, email, username, password string) (*types.AuthenticatedUserResponse, error) {
	err := s.validateDetails(email, username)
	if err != nil {
		return nil, types.NewAPIError(http.StatusBadRequest, "no credentials were passed")
	}

	existingUserCount, err := s.q.CheckForUser(ctx, queries.CheckForUserParams{
		Email:    email,
		Username: username,
	})
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, err.Error())
	}

	if existingUserCount == 0 {
		return nil, types.NewAPIError(http.StatusNotFound, "user not found")
	}

	existingUser, err := s.q.GetExistingUser(ctx, queries.GetExistingUserParams{
		Email:    email,
		Username: username,
	})
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error fetching user: "+err.Error())
	}

	authenticated := helpers.CheckPasswordHash(password, existingUser.HashedPassword)

	if !authenticated {
		return nil, types.NewAPIError(http.StatusBadRequest, "incorrect password")
	}

	// Generate an access token
	accessToken, expiryDate, err := helpers.GenerateAccessToken(existingUser)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, err.Error())
	}

	// Generate a refresh token
	refreshToken, err := s.generateAndStoreRefreshToken(existingUser, ctx)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, err.Error())
	}

	userResponse := types.NewAuthenticateUserResponse(existingUser.Uuid.String(), accessToken, *refreshToken, expiryDate, existingUser.Superuser)

	return userResponse, nil
}

// LogoutUser logs out the user
func (s *AuthService) LogoutUser(ctx context.Context, refreshToken string) error {
	_, err := s.q.GetRefreshTokenByToken(ctx, refreshToken)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return types.NewAPIError(http.StatusUnauthorized, "already logged out")
	}

	err = s.q.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "failed to log out: "+err.Error())
	}

	return nil
}

// CreateNewAccessToken creates a new auth token for the logged-in user
func (s *AuthService) CreateNewAccessToken(c *gin.Context, refreshToken *string) (*types.AccessTokenResponse, error) {
	userUuid, err := helpers.ValidateUserUuidFromClaims(c)
	if err != nil {
		return nil, err
	}

	if userUuid == nil {
		return nil, types.NewAPIError(http.StatusBadRequest, "missing user uuid")
	}

	pgUuid, err := helpers.ValidateAndConvertUUID(*userUuid)
	if err != nil {
		return nil, err
	}

	existingUser, err := s.q.GetUserByUuid(c.Request.Context(), *pgUuid)
	if err != nil {
		return nil, types.NewAPIError(http.StatusNotFound, "user not found")
	}

	existingToken, err := s.q.GetRefreshTokenByToken(c.Request.Context(), *refreshToken)
	if err != nil {
		return nil, types.NewAPIError(http.StatusBadRequest, "missing user uuid")
	}

	if time.Now().After(existingToken.ExpiresAt.Time) {
		return nil, types.NewAPIError(http.StatusBadRequest, "refresh token expired")
	}

	accessToken, expiryDate, err := helpers.GenerateAccessToken(existingUser)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error generating access token")
	}

	c.Set("user_uuid", userUuid)
	c.Set("access_token", accessToken)

	tokenResponse := types.NewAccessTokenResponse(accessToken, expiryDate)

	return tokenResponse, nil
}

// createDefaultData adds default data for a user
func (s *AuthService) createDefaultData(ctx context.Context, user queries.AddUserRow) error {

	// Create a default list
	list, err := s.q.AddList(ctx, queries.AddListParams{
		Name:     "Watchlist",
		UserUuid: user.Uuid,
	})
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, err.Error())
	}

	// Create a default status
	status, err := s.q.AddStatus(ctx, queries.AddStatusParams{
		StatusLabel: helpers.MakePgString("backlog"),
		UserUuid:    user.Uuid,
	})
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, err.Error())
	}

	// Assign the default status to the default list
	_, err = s.q.AddListStatus(ctx, queries.AddListStatusParams{
		ListUuid:   list.Uuid,
		StatusUuid: status.Uuid,
	})
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

// validateDetails validates that both email and username are populated
func (s *AuthService) validateDetails(email, username string) error {
	if email == "" && username == "" {
		return types.NewAPIError(http.StatusBadRequest, "email or username required")
	}
	return nil
}

// generateAndStoreRefreshToken creates and stores a refresh token for the user
func (s *AuthService) generateAndStoreRefreshToken(user queries.GetExistingUserRow, ctx context.Context) (*string, error) {
	// Generate a refresh token
	refreshToken, err := helpers.GenerateRefreshToken(64)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error generating refresh token: "+err.Error())
	}

	// Set the refresh token to expire in 7 days
	refreshTokenExpiry := time.Now().Add(time.Hour * 24 * 7)

	RefreshTokenParams := queries.AddRefreshTokenParams{
		UserID:    user.UserID.Int64,
		Token:     *refreshToken,
		ExpiresAt: pgtype.Timestamptz{Time: refreshTokenExpiry, Valid: true},
	}

	// Store the refresh token in the database
	refreshTokenRow, err := s.q.AddRefreshToken(ctx, RefreshTokenParams)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error generating refresh token: "+err.Error())
	}

	return &refreshTokenRow.Token, nil
}
