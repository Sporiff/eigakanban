package services

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/types"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
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
	existingUserCount, err := s.q.CheckForUser(ctx, queries.CheckForUserParams{
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		return nil, err
	}

	// If the user already exists, return an error
	if existingUserCount > 0 {
		return nil, errors.New("user already exists: " + user.Username)
	}

	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// Add the user to the database
	registeredUser, err := s.q.AddUser(ctx, queries.AddUserParams{
		Username:       user.Username,
		HashedPassword: hashedPassword,
		Email:          user.Email,
	})
	if err != nil {
		return nil, err
	}

	// Create default data for the user
	err = s.createDefaultData(ctx, registeredUser)
	if err != nil {
		return nil, errors.New("error creating default data: " + err.Error())
	}

	return &registeredUser, nil
}

// createDefaultData adds default data for a user
func (s *AuthService) createDefaultData(ctx context.Context, user queries.AddUserRow) error {

	// Create a default list
	list, err := s.q.AddList(ctx, queries.AddListParams{
		Name:     "Watchlist",
		UserUuid: user.Uuid,
	})
	if err != nil {
		return err
	}

	// Create a default status
	status, err := s.q.AddStatus(ctx, queries.AddStatusParams{
		StatusLabel: helpers.MakePgString("backlog"),
		UserUuid:    user.Uuid,
	})
	if err != nil {
		return err
	}

	_, err = s.q.AddListStatus(ctx, queries.AddListStatusParams{
		ListUuid:   list.Uuid,
		StatusUuid: status.Uuid,
	})
	if err != nil {
		return err
	}

	return nil
}

// LoginUser logs in the user and sets up authentication
func (s *AuthService) LoginUser(ctx context.Context, email, username, password string) (types.AuthenticatedUserResponse, error) {
	var userResponse = types.AuthenticatedUserResponse{}

	existingUser, err := s.checkForUser(ctx, email, username, password)
	if err != nil {
		return userResponse, err
	}

	// Generate an access token
	accessToken, err := helpers.GenerateAccessToken(existingUser)
	if err != nil {
		return userResponse, err
	}

	// Generate a refresh token
	refreshToken, err := s.generateAndStoreRefreshToken(existingUser, ctx)
	if err != nil {
		return userResponse, err
	}

	userResponse.Init(existingUser.Uuid.String(), accessToken, refreshToken)

	return userResponse, err
}

func (s *AuthService) checkForUser(ctx context.Context, email, username, password string) (queries.GetExistingUserRow, error) {
	// Create an empty user to return if user doesn't exist
	emptyUser := queries.GetExistingUserRow{}

	// Fetch the user from the database
	existingUser, err := s.q.GetExistingUser(ctx, queries.GetExistingUserParams{
		Email:    email,
		Username: username,
	})
	if err != nil {
		return emptyUser, err
	}

	// Check if the user exists
	if existingUser == emptyUser {
		return existingUser, errors.New("user not found")
	}

	// Verify the password matches the stored password
	if !helpers.CheckPasswordHash(password, existingUser.HashedPassword) {
		return existingUser, errors.New("invalid password")
	}

	return existingUser, nil
}

// validateDetails validates that both email and username are populated
func (s *AuthService) validateDetails(email, username string) error {
	// Ensure either email or username is provided
	if email == "" && username == "" {
		return errors.New("email and username are required")
	}
	return nil
}

// generateAndStoreRefreshToken creates and stores a refresh token for the user
func (s *AuthService) generateAndStoreRefreshToken(user queries.GetExistingUserRow, ctx context.Context) (string, error) {
	// Generate a refresh token
	refreshToken, err := helpers.GenerateRefreshToken(64)
	if err != nil {
		return "", errors.New("error generating refresh token: " + err.Error())
	}

	// Set the refresh token to expire in 7 days
	refreshTokenExpiry := time.Now().Add(time.Hour * 24 * 7)

	RefreshTokenParams := queries.AddRefreshTokenParams{
		UserID:    user.UserID.Int64,
		Token:     refreshToken,
		ExpiresAt: pgtype.Timestamptz{Time: refreshTokenExpiry, Valid: true},
	}

	// Store the refresh token in the database
	refreshTokenRow, err := s.q.AddRefreshToken(ctx, RefreshTokenParams)
	if err != nil {
		return "", errors.New("error generating refresh token: " + err.Error())
	}

	return refreshTokenRow.Token, nil
}

// LogoutUser logs out the user
func (s *AuthService) LogoutUser(ctx context.Context, refreshToken string) error {
	_, err := s.q.GetRefreshTokenByToken(ctx, refreshToken)
	if err != nil {
		return errors.New("already logged out")
	}

	err = s.q.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		return errors.New("failed to log out: " + err.Error())
	}

	return nil
}
