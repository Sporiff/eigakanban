package services

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/types"
	"context"
	"errors"
)

type UsersService struct {
	q *queries.Queries
}

func NewUsersService(q *queries.Queries) *UsersService {
	return &UsersService{q: q}
}

// GetAllUsers returns the total number of users and a list of all users
func (s *UsersService) GetAllUsers(ctx context.Context, pagination *types.Pagination) ([]queries.GetAllUsersRow, *types.Pagination, error) {
	// Get the total number of users
	total, err := s.q.GetUserCount(ctx)
	if err != nil {
		return nil, nil, errors.New("error getting user count: " + err.Error())
	}

	// Update the pagination with the total count
	pagination.Total = total

	// Fetch the users for the current page
	users, err := s.q.GetAllUsers(ctx, queries.GetAllUsersParams{
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	})
	if err != nil {
		return nil, nil, errors.New("error getting all users: " + err.Error())
	}

	return users, pagination, nil
}

// GetUserByUuid fetches a user by UUID
func (s *UsersService) GetUserByUuid(ctx context.Context, uuid string) (queries.GetUserByUuidRow, error) {
	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return queries.GetUserByUuidRow{}, errors.New("error validating uuid: " + err.Error())
	}

	user, err := s.q.GetUserByUuid(ctx, pgUuid)
	if err != nil {
		return queries.GetUserByUuidRow{}, errors.New("error getting user by uuid: " + err.Error())
	}

	return user, nil
}

// UpdateUser updates user details
func (s *UsersService) UpdateUser(ctx context.Context, uuid string, user types.UpdateUserRequest) (queries.UpdateUserDetailsRow, error) {
	emptyUser := queries.UpdateUserDetailsRow{}

	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return emptyUser, errors.New("error validating uuid: " + err.Error())
	}

	params := queries.UpdateUserDetailsParams{
		UserUuid: pgUuid,
	}

	// If no username is provided, fetch the current value
	if user.NewUsername != nil {
		params.NewUsername = *user.NewUsername
	} else {
		currentUser, err := s.q.GetUserByUuid(ctx, pgUuid)
		if err != nil {
			return emptyUser, errors.New("error getting user by uuid: " + err.Error())
		}
		params.NewUsername = currentUser.Username
	}

	// Handle NewName and NewBio using helper functions
	helpers.AssignPgtypeText(&params.NewName, user.NewName)
	helpers.AssignPgtypeText(&params.NewBio, user.NewBio)

	// Update the user details
	userRow, err := s.q.UpdateUserDetails(ctx, params)
	if err != nil {
		return emptyUser, errors.New("error updating user: " + err.Error())
	}

	return userRow, nil
}

// DeleteUser deletes a user from the database
func (s *UsersService) DeleteUser(ctx context.Context, uuid string) error {
	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return errors.New("error validating uuidL " + err.Error())
	}

	err = s.q.DeleteUser(ctx, pgUuid)
	if err != nil {
		return errors.New("error deleting user: " + err.Error())
	}

	return nil
}

// GetUserCount returns a count of all users
func (s *UsersService) GetUserCount(ctx context.Context) (int64, error) {
	return s.q.GetUserCount(ctx)
}
