package services

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/types"
	"context"
	"database/sql"
	"errors"
	"net/http"
)

type UsersService struct {
	q *queries.Queries
}

func NewUsersService(q *queries.Queries) *UsersService {
	return &UsersService{q: q}
}

// GetAllUsers returns the total number of users and a list of all users
func (s *UsersService) GetAllUsers(ctx context.Context, pagination *types.Pagination) (*types.PaginatedUsersResponse, error) {
	// Get the total number of users
	total, err := s.GetUserCount(ctx)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error getting user count")
	}

	if total == 0 {
		response := &types.PaginatedUsersResponse{
			Pagination: *pagination,
			Users:      []types.UserResponse{},
		}
		return response, nil
	}

	// Update the pagination with the total count
	pagination.Total = total

	// Fetch the users for the current page
	items, err := s.q.GetAllUsers(ctx, queries.GetAllUsersParams{
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	})
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error getting users")
	}

	users := make([]types.UserResponse, len(items))

	for i, item := range items {
		users[i] = types.UserResponse{
			UUID:     item.Uuid.String(),
			Username: item.Username,
			FullName: item.FullName.String,
			Bio:      item.Bio.String,
		}
	}

	response := types.PaginatedUsersResponse{
		Pagination: *pagination,
		Users:      users,
	}

	return &response, nil
}

// GetUserByUuid fetches a user by UUID
func (s *UsersService) GetUserByUuid(ctx context.Context, uuid string) (*queries.GetUserByUuidRow, error) {
	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return nil, err
	}

	user, err := s.q.GetUserByUuid(ctx, *pgUuid)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error getting user")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewAPIError(http.StatusNotFound, "user not found")
	}

	return &user, nil
}

// UpdateUser updates user details
func (s *UsersService) UpdateUser(ctx context.Context, uuid string, user types.UpdateUserRequest) (*queries.UpdateUserDetailsRow, error) {
	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return nil, err
	}

	existingUser, err := s.q.GetUserByUuid(ctx, *pgUuid)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error getting user")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewAPIError(http.StatusNotFound, "user not found")
	}

	params := queries.UpdateUserDetailsParams{
		UserUuid: *pgUuid,
	}

	// If no username is provided, fetch the current value
	if user.NewUsername != nil {
		params.NewUsername = *user.NewUsername
	} else {
		params.NewUsername = existingUser.Username
	}

	// Handle NewName and NewBio using helper functions
	helpers.AssignPgtypeText(&params.NewName, user.NewName)
	helpers.AssignPgtypeText(&params.NewBio, user.NewBio)

	// Update the user details
	userRow, err := s.q.UpdateUserDetails(ctx, params)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error updating user")
	}

	return &userRow, nil
}

// DeleteUser deletes a user from the database
func (s *UsersService) DeleteUser(ctx context.Context, uuid string) error {
	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return err
	}

	err = s.q.DeleteUser(ctx, *pgUuid)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "error deleting user")
	}

	return nil
}

// GetUserCount returns a count of all users
func (s *UsersService) GetUserCount(ctx context.Context) (int64, error) {
	return s.q.GetUserCount(ctx)
}
