// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user_queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addUser = `-- name: AddUser :one
INSERT INTO
    users (username, hashed_password, email, full_name, bio)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING
    uuid,
    username,
    full_name,
    bio,
    created_date
`

type AddUserParams struct {
	Username       string      `json:"username"`
	HashedPassword string      `json:"hashed_password"`
	Email          string      `json:"email"`
	FullName       pgtype.Text `json:"full_name"`
	Bio            pgtype.Text `json:"bio"`
}

type AddUserRow struct {
	Uuid        pgtype.UUID        `json:"uuid"`
	Username    string             `json:"username"`
	FullName    pgtype.Text        `json:"full_name"`
	Bio         pgtype.Text        `json:"bio"`
	CreatedDate pgtype.Timestamptz `json:"created_date"`
}

func (q *Queries) AddUser(ctx context.Context, arg AddUserParams) (AddUserRow, error) {
	row := q.db.QueryRow(ctx, addUser,
		arg.Username,
		arg.HashedPassword,
		arg.Email,
		arg.FullName,
		arg.Bio,
	)
	var i AddUserRow
	err := row.Scan(
		&i.Uuid,
		&i.Username,
		&i.FullName,
		&i.Bio,
		&i.CreatedDate,
	)
	return i, err
}

const checkForUser = `-- name: CheckForUser :one
SELECT COUNT(*)
FROM users
WHERE
    email = $1
   OR
    username = $2
`

type CheckForUserParams struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (q *Queries) CheckForUser(ctx context.Context, arg CheckForUserParams) (int64, error) {
	row := q.db.QueryRow(ctx, checkForUser, arg.Email, arg.Username)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE
    uuid = $1
`

func (q *Queries) DeleteUser(ctx context.Context, userUuid pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, userUuid)
	return err
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT
    uuid,
    username,
    full_name,
    bio
FROM
    users
ORDER BY
    created_date
LIMIT
    $2
    OFFSET
    $1
`

type GetAllUsersParams struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}

type GetAllUsersRow struct {
	Uuid     pgtype.UUID `json:"uuid"`
	Username string      `json:"username"`
	FullName pgtype.Text `json:"full_name"`
	Bio      pgtype.Text `json:"bio"`
}

func (q *Queries) GetAllUsers(ctx context.Context, arg GetAllUsersParams) ([]GetAllUsersRow, error) {
	rows, err := q.db.Query(ctx, getAllUsers, arg.Page, arg.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllUsersRow
	for rows.Next() {
		var i GetAllUsersRow
		if err := rows.Scan(
			&i.Uuid,
			&i.Username,
			&i.FullName,
			&i.Bio,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getExistingUser = `-- name: GetExistingUser :one
SELECT
    user_id,
    uuid,
    username,
    email,
    hashed_password
FROM
    users
WHERE
    email = $1
OR
    username = $2
LIMIT
    1
`

type GetExistingUserParams struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type GetExistingUserRow struct {
	UserID         pgtype.Int8 `json:"user_id"`
	Uuid           pgtype.UUID `json:"uuid"`
	Username       string      `json:"username"`
	Email          string      `json:"email"`
	HashedPassword string      `json:"hashed_password"`
}

func (q *Queries) GetExistingUser(ctx context.Context, arg GetExistingUserParams) (GetExistingUserRow, error) {
	row := q.db.QueryRow(ctx, getExistingUser, arg.Email, arg.Username)
	var i GetExistingUserRow
	err := row.Scan(
		&i.UserID,
		&i.Uuid,
		&i.Username,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT
    uuid,
    username,
    full_name,
    bio
FROM
    users
WHERE
    user_id = $1
LIMIT
    1
`

type GetUserByIdRow struct {
	Uuid     pgtype.UUID `json:"uuid"`
	Username string      `json:"username"`
	FullName pgtype.Text `json:"full_name"`
	Bio      pgtype.Text `json:"bio"`
}

func (q *Queries) GetUserById(ctx context.Context, userID pgtype.Int8) (GetUserByIdRow, error) {
	row := q.db.QueryRow(ctx, getUserById, userID)
	var i GetUserByIdRow
	err := row.Scan(
		&i.Uuid,
		&i.Username,
		&i.FullName,
		&i.Bio,
	)
	return i, err
}

const getUserByUuid = `-- name: GetUserByUuid :one
SELECT
    uuid,
    username,
    full_name,
    bio
FROM
    users
WHERE
    uuid = $1
LIMIT
    1
`

type GetUserByUuidRow struct {
	Uuid     pgtype.UUID `json:"uuid"`
	Username string      `json:"username"`
	FullName pgtype.Text `json:"full_name"`
	Bio      pgtype.Text `json:"bio"`
}

func (q *Queries) GetUserByUuid(ctx context.Context, userUuid pgtype.UUID) (GetUserByUuidRow, error) {
	row := q.db.QueryRow(ctx, getUserByUuid, userUuid)
	var i GetUserByUuidRow
	err := row.Scan(
		&i.Uuid,
		&i.Username,
		&i.FullName,
		&i.Bio,
	)
	return i, err
}

const getUserCount = `-- name: GetUserCount :one
SELECT COUNT (*)
FROM users
`

func (q *Queries) GetUserCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getUserCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const updateUserDetails = `-- name: UpdateUserDetails :one
UPDATE users
SET
    username = COALESCE($1, username),
    full_name = COALESCE($2, full_name),
    bio = COALESCE($3, bio)
WHERE
    uuid = $4
RETURNING
    uuid,
    username,
    full_name,
    bio
`

type UpdateUserDetailsParams struct {
	NewUsername string      `json:"new_username"`
	NewName     pgtype.Text `json:"new_name"`
	NewBio      pgtype.Text `json:"new_bio"`
	UserUuid    pgtype.UUID `json:"user_uuid"`
}

type UpdateUserDetailsRow struct {
	Uuid     pgtype.UUID `json:"uuid"`
	Username string      `json:"username"`
	FullName pgtype.Text `json:"full_name"`
	Bio      pgtype.Text `json:"bio"`
}

func (q *Queries) UpdateUserDetails(ctx context.Context, arg UpdateUserDetailsParams) (UpdateUserDetailsRow, error) {
	row := q.db.QueryRow(ctx, updateUserDetails,
		arg.NewUsername,
		arg.NewName,
		arg.NewBio,
		arg.UserUuid,
	)
	var i UpdateUserDetailsRow
	err := row.Scan(
		&i.Uuid,
		&i.Username,
		&i.FullName,
		&i.Bio,
	)
	return i, err
}
