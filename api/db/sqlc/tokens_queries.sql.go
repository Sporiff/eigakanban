// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: tokens_queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addRefreshToken = `-- name: AddRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expires_at)
VALUES (
        $1, $2, $3
       )
RETURNING token, expires_at
`

type AddRefreshTokenParams struct {
	UserID    int64              `json:"user_id"`
	Token     string             `json:"token"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
}

type AddRefreshTokenRow struct {
	Token     string             `json:"token"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
}

func (q *Queries) AddRefreshToken(ctx context.Context, arg AddRefreshTokenParams) (AddRefreshTokenRow, error) {
	row := q.db.QueryRow(ctx, addRefreshToken, arg.UserID, arg.Token, arg.ExpiresAt)
	var i AddRefreshTokenRow
	err := row.Scan(&i.Token, &i.ExpiresAt)
	return i, err
}

const deleteRefreshToken = `-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = $1
`

func (q *Queries) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := q.db.Exec(ctx, deleteRefreshToken, token)
	return err
}

const getRefreshTokenByToken = `-- name: GetRefreshTokenByToken :one
SELECT
    user_id,
    token,
    expires_at
FROM refresh_tokens
WHERE token = $1
`

type GetRefreshTokenByTokenRow struct {
	UserID    int64              `json:"user_id"`
	Token     string             `json:"token"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
}

func (q *Queries) GetRefreshTokenByToken(ctx context.Context, token string) (GetRefreshTokenByTokenRow, error) {
	row := q.db.QueryRow(ctx, getRefreshTokenByToken, token)
	var i GetRefreshTokenByTokenRow
	err := row.Scan(&i.UserID, &i.Token, &i.ExpiresAt)
	return i, err
}
