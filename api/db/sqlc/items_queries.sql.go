// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: items_queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addItem = `-- name: AddItem :one
INSERT INTO
    items (title)
VALUES
    (
        $1
    )
RETURNING
    uuid,
    title,
    created_date
`

type AddItemRow struct {
	Uuid        pgtype.UUID        `json:"uuid"`
	Title       string             `json:"title"`
	CreatedDate pgtype.Timestamptz `json:"created_date"`
}

func (q *Queries) AddItem(ctx context.Context, title string) (AddItemRow, error) {
	row := q.db.QueryRow(ctx, addItem, title)
	var i AddItemRow
	err := row.Scan(&i.Uuid, &i.Title, &i.CreatedDate)
	return i, err
}

const deleteItem = `-- name: DeleteItem :exec
DELETE FROM items
WHERE
    uuid = $1
`

func (q *Queries) DeleteItem(ctx context.Context, itemUuid pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteItem, itemUuid)
	return err
}

const getAllItems = `-- name: GetAllItems :many
SELECT
    uuid,
    title
FROM
    items
ORDER BY
    created_date
LIMIT
    $2
    OFFSET
    $1
`

type GetAllItemsParams struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}

type GetAllItemsRow struct {
	Uuid  pgtype.UUID `json:"uuid"`
	Title string      `json:"title"`
}

func (q *Queries) GetAllItems(ctx context.Context, arg GetAllItemsParams) ([]GetAllItemsRow, error) {
	rows, err := q.db.Query(ctx, getAllItems, arg.Page, arg.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllItemsRow
	for rows.Next() {
		var i GetAllItemsRow
		if err := rows.Scan(&i.Uuid, &i.Title); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getItemByUuid = `-- name: GetItemByUuid :one
SELECT
    uuid,
    title,
    created_date
FROM
    items
WHERE
    uuid = $1
LIMIT
    1
`

type GetItemByUuidRow struct {
	Uuid        pgtype.UUID        `json:"uuid"`
	Title       string             `json:"title"`
	CreatedDate pgtype.Timestamptz `json:"created_date"`
}

func (q *Queries) GetItemByUuid(ctx context.Context, itemUuid pgtype.UUID) (GetItemByUuidRow, error) {
	row := q.db.QueryRow(ctx, getItemByUuid, itemUuid)
	var i GetItemByUuidRow
	err := row.Scan(&i.Uuid, &i.Title, &i.CreatedDate)
	return i, err
}

const getItemsCount = `-- name: GetItemsCount :one
SELECT COUNT (*)
FROM items
`

func (q *Queries) GetItemsCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getItemsCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const updateItem = `-- name: UpdateItem :one
UPDATE items
SET title = COALESCE($1, title)
WHERE uuid = $2
RETURNING
    uuid,
    title
`

type UpdateItemParams struct {
	ItemTitle string      `json:"item_title"`
	ItemUuid  pgtype.UUID `json:"item_uuid"`
}

type UpdateItemRow struct {
	Uuid  pgtype.UUID `json:"uuid"`
	Title string      `json:"title"`
}

func (q *Queries) UpdateItem(ctx context.Context, arg UpdateItemParams) (UpdateItemRow, error) {
	row := q.db.QueryRow(ctx, updateItem, arg.ItemTitle, arg.ItemUuid)
	var i UpdateItemRow
	err := row.Scan(&i.Uuid, &i.Title)
	return i, err
}
