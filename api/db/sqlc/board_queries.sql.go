// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: board_queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addBoard = `-- name: AddBoard :one
INSERT INTO
    boards (name, description, user_id)
VALUES
    (
        $1,
        $2,
        (
            SELECT
                user_id
            FROM
                users
            WHERE
                users.uuid = $3
        )
    )
RETURNING
    uuid,
    name,
    description,
    created_date
`

type AddBoardParams struct {
	BoardName        string      `json:"board_name"`
	BoardDescription pgtype.Text `json:"board_description"`
	UserUuid         pgtype.UUID `json:"user_uuid"`
}

type AddBoardRow struct {
	Uuid        pgtype.UUID        `json:"uuid"`
	Name        string             `json:"name"`
	Description pgtype.Text        `json:"description"`
	CreatedDate pgtype.Timestamptz `json:"created_date"`
}

func (q *Queries) AddBoard(ctx context.Context, arg AddBoardParams) (AddBoardRow, error) {
	row := q.db.QueryRow(ctx, addBoard, arg.BoardName, arg.BoardDescription, arg.UserUuid)
	var i AddBoardRow
	err := row.Scan(
		&i.Uuid,
		&i.Name,
		&i.Description,
		&i.CreatedDate,
	)
	return i, err
}

const deleteBoard = `-- name: DeleteBoard :exec
DELETE FROM boards
WHERE
    uuid = $1
`

func (q *Queries) DeleteBoard(ctx context.Context, boardUuid pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteBoard, boardUuid)
	return err
}

const getAllBoards = `-- name: GetAllBoards :many
SELECT
    name,
    description
FROM
    boards
ORDER BY
    created_date
LIMIT
    $2
    OFFSET
    $1
`

type GetAllBoardsParams struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}

type GetAllBoardsRow struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
}

func (q *Queries) GetAllBoards(ctx context.Context, arg GetAllBoardsParams) ([]GetAllBoardsRow, error) {
	rows, err := q.db.Query(ctx, getAllBoards, arg.Page, arg.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllBoardsRow
	for rows.Next() {
		var i GetAllBoardsRow
		if err := rows.Scan(&i.Name, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBoardByUuid = `-- name: GetBoardByUuid :one
SELECT
    name,
    description
FROM
    boards
WHERE
    uuid = $1
LIMIT
    1
`

type GetBoardByUuidRow struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
}

func (q *Queries) GetBoardByUuid(ctx context.Context, boardUuid pgtype.UUID) (GetBoardByUuidRow, error) {
	row := q.db.QueryRow(ctx, getBoardByUuid, boardUuid)
	var i GetBoardByUuidRow
	err := row.Scan(&i.Name, &i.Description)
	return i, err
}

const getBoardsForUser = `-- name: GetBoardsForUser :many
SELECT
    b.name,
    b.description
FROM
    boards b
        JOIN users u ON b.user_id = u.user_id
WHERE
    u.uuid = $1
ORDER BY
    b.created_date
LIMIT
    $3
    OFFSET
    $2
`

type GetBoardsForUserParams struct {
	UserUuid pgtype.UUID `json:"user_uuid"`
	Page     int32       `json:"page"`
	PageSize int32       `json:"page_size"`
}

type GetBoardsForUserRow struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
}

func (q *Queries) GetBoardsForUser(ctx context.Context, arg GetBoardsForUserParams) ([]GetBoardsForUserRow, error) {
	rows, err := q.db.Query(ctx, getBoardsForUser, arg.UserUuid, arg.Page, arg.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBoardsForUserRow
	for rows.Next() {
		var i GetBoardsForUserRow
		if err := rows.Scan(&i.Name, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBoard = `-- name: UpdateBoard :one
UPDATE boards
SET
    name = $1,
    description = $2
WHERE
    uuid = $3
RETURNING
    uuid,
    name,
    description
`

type UpdateBoardParams struct {
	BoardName        string      `json:"board_name"`
	BoardDescription pgtype.Text `json:"board_description"`
	BoardUuid        pgtype.UUID `json:"board_uuid"`
}

type UpdateBoardRow struct {
	Uuid        pgtype.UUID `json:"uuid"`
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
}

func (q *Queries) UpdateBoard(ctx context.Context, arg UpdateBoardParams) (UpdateBoardRow, error) {
	row := q.db.QueryRow(ctx, updateBoard, arg.BoardName, arg.BoardDescription, arg.BoardUuid)
	var i UpdateBoardRow
	err := row.Scan(&i.Uuid, &i.Name, &i.Description)
	return i, err
}
