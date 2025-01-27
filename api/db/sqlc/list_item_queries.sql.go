// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: list_item_queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addItemToListAtPosition = `-- name: AddItemToListAtPosition :one
WITH new_item AS (
    INSERT INTO list_items (list_id, item_id, position, prev_item_id, next_item_id)
        VALUES (
                   (SELECT list_id FROM lists WHERE lists.uuid = $1),
                   (SELECT item_id FROM items WHERE items.uuid = $2),
                   $3,
                   (SELECT list_item_id FROM list_items WHERE list_id = (SELECT list_id FROM lists WHERE uuid = $1) AND list_items.position = $4),
                   (SELECT list_item_id FROM list_items WHERE list_id = (SELECT list_id FROM lists WHERE uuid = $1) AND list_items.position = $5)
               )
        RETURNING list_item_id, list_id, position, prev_item_id, next_item_id
),
     update_prev_item AS (
         UPDATE list_items
             SET next_item_id = new_item.list_item_id
             FROM new_item
             WHERE list_items.list_item_id = new_item.prev_item_id
                 AND new_item.prev_item_id IS NOT NULL
     ),
     update_next_item AS (
         UPDATE list_items
             SET prev_item_id = new_item.list_item_id
             FROM new_item
             WHERE list_items.list_item_id = new_item.next_item_id
                 AND new_item.next_item_id IS NOT NULL
     )
SELECT list_items.uuid, list_items.position
FROM list_items
         JOIN items ON list_items.item_id = items.item_id
WHERE list_items.list_item_id = (SELECT list_item_id FROM new_item)
`

type AddItemToListAtPositionParams struct {
	ListUuid     pgtype.UUID `json:"list_uuid"`
	ItemUuid     pgtype.UUID `json:"item_uuid"`
	NewPosition  int32       `json:"new_position"`
	PrevPosition int32       `json:"prev_position"`
	NextPosition int32       `json:"next_position"`
}

type AddItemToListAtPositionRow struct {
	Uuid     pgtype.UUID `json:"uuid"`
	Position int32       `json:"position"`
}

// Arguments: list_uuid, item_uuid, new_position
func (q *Queries) AddItemToListAtPosition(ctx context.Context, arg AddItemToListAtPositionParams) (AddItemToListAtPositionRow, error) {
	row := q.db.QueryRow(ctx, addItemToListAtPosition,
		arg.ListUuid,
		arg.ItemUuid,
		arg.NewPosition,
		arg.PrevPosition,
		arg.NextPosition,
	)
	var i AddItemToListAtPositionRow
	err := row.Scan(&i.Uuid, &i.Position)
	return i, err
}

const deleteItemFromList = `-- name: DeleteItemFromList :exec
WITH deleted_item AS (
    DELETE FROM list_items
        WHERE list_items.uuid = $1
        RETURNING list_item_id, prev_item_id, next_item_id
),
     update_prev_item AS (
         UPDATE list_items
             SET next_item_id = deleted_item.next_item_id
             FROM deleted_item
             WHERE list_items.list_item_id = deleted_item.prev_item_id
                 AND deleted_item.prev_item_id IS NOT NULL
     ),
     update_next_item AS (
         UPDATE list_items
             SET prev_item_id = deleted_item.prev_item_id
             FROM deleted_item
             WHERE list_items.list_item_id = deleted_item.next_item_id
                 AND deleted_item.next_item_id IS NOT NULL
     )
SELECT 1
`

// Arguments: list_item_uuid
func (q *Queries) DeleteItemFromList(ctx context.Context, listItemUuid pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteItemFromList, listItemUuid)
	return err
}

const moveItemInList = `-- name: MoveItemInList :one
WITH current_item AS (
    SELECT list_item_id, prev_item_id, next_item_id, position, list_id
    FROM list_items
    WHERE list_items.uuid = $1
),
     updated_positions AS (
         UPDATE list_items
             SET position = position + CASE
                                           WHEN list_items.position >= $2 THEN 1
                                           ELSE -1
                 END
             WHERE list_id = (SELECT list_id FROM current_item)
                 AND position BETWEEN LEAST($2, (SELECT position FROM current_item))
                       AND GREATEST($2, (SELECT position FROM current_item))
             RETURNING list_item_id, position
     ),
     update_current_item AS (
         UPDATE list_items
             SET
                 position = $2,
                 prev_item_id = (SELECT list_item_id FROM list_items WHERE list_id = (SELECT list_id FROM current_item) AND list_items.position = $3),
                 next_item_id = (SELECT list_item_id FROM list_items WHERE list_id = (SELECT list_id FROM current_item) AND list_items.position = $4)
             WHERE list_item_id = (SELECT list_item_id FROM current_item)
             RETURNING uuid, position, prev_item_id, next_item_id
     )
SELECT uuid, position, prev_item_id, next_item_id
FROM update_current_item
`

type MoveItemInListParams struct {
	ListItemUuid pgtype.UUID `json:"list_item_uuid"`
	NewPosition  int32       `json:"new_position"`
	PrevPosition int32       `json:"prev_position"`
	NextPosition int32       `json:"next_position"`
}

type MoveItemInListRow struct {
	Uuid       pgtype.UUID `json:"uuid"`
	Position   int32       `json:"position"`
	PrevItemID pgtype.Int8 `json:"prev_item_id"`
	NextItemID pgtype.Int8 `json:"next_item_id"`
}

// Arguments: list_item_uuid, new_position
func (q *Queries) MoveItemInList(ctx context.Context, arg MoveItemInListParams) (MoveItemInListRow, error) {
	row := q.db.QueryRow(ctx, moveItemInList,
		arg.ListItemUuid,
		arg.NewPosition,
		arg.PrevPosition,
		arg.NextPosition,
	)
	var i MoveItemInListRow
	err := row.Scan(
		&i.Uuid,
		&i.Position,
		&i.PrevItemID,
		&i.NextItemID,
	)
	return i, err
}

const moveItemToAnotherList = `-- name: MoveItemToAnotherList :one
WITH current_item AS (
    SELECT list_item_id, item_id, prev_item_id, next_item_id, position, list_id
    FROM list_items
    WHERE list_items.uuid = $1
),
     update_source_list AS (
         UPDATE list_items
             SET
                 next_item_id = current_item.next_item_id
             FROM current_item
             WHERE list_items.list_item_id = current_item.prev_item_id
                 AND current_item.prev_item_id IS NOT NULL
     ),
     update_source_list_next AS (
         UPDATE list_items
             SET
                 prev_item_id = current_item.prev_item_id
             FROM current_item
             WHERE list_items.list_item_id = current_item.next_item_id
                 AND current_item.next_item_id IS NOT NULL
     ),
     new_item AS (
         INSERT INTO list_items (list_id, item_id, position, prev_item_id, next_item_id)
             VALUES (
                        (SELECT list_id FROM lists WHERE lists.uuid = $2),
                        (SELECT item_id FROM current_item),
                        $3,
                        (SELECT list_item_id FROM list_items WHERE list_id = (SELECT list_id FROM lists WHERE lists.uuid = $2) AND list_items.position = $4),
                        (SELECT list_item_id FROM list_items WHERE list_id = (SELECT list_id FROM lists WHERE lists.uuid = $2) AND list_items.position = $5)
                    )
             RETURNING list_item_id, list_id, position, prev_item_id, next_item_id
     ),
     update_destination_list AS (
         UPDATE list_items
             SET
                 next_item_id = new_item.list_item_id
             FROM new_item
             WHERE list_items.list_item_id = new_item.prev_item_id
                 AND new_item.prev_item_id IS NOT NULL
     ),
     update_destination_list_next AS (
         UPDATE list_items
             SET
                 prev_item_id = new_item.list_item_id
             FROM new_item
             WHERE list_items.list_item_id = new_item.next_item_id
                 AND new_item.next_item_id IS NOT NULL
     )
SELECT uuid, position, prev_item_id, next_item_id
FROM list_items
WHERE list_items.list_item_id = (SELECT list_item_id FROM new_item)
`

type MoveItemToAnotherListParams struct {
	ListItemUuid   pgtype.UUID `json:"list_item_uuid"`
	TargetListUuid pgtype.UUID `json:"target_list_uuid"`
	NewPosition    int32       `json:"new_position"`
	PrevPosition   int32       `json:"prev_position"`
	NextPosition   int32       `json:"next_position"`
}

type MoveItemToAnotherListRow struct {
	Uuid       pgtype.UUID `json:"uuid"`
	Position   int32       `json:"position"`
	PrevItemID pgtype.Int8 `json:"prev_item_id"`
	NextItemID pgtype.Int8 `json:"next_item_id"`
}

// Arguments: list_item_uuid, target_list_uuid, new_position, prev_position, next_position
func (q *Queries) MoveItemToAnotherList(ctx context.Context, arg MoveItemToAnotherListParams) (MoveItemToAnotherListRow, error) {
	row := q.db.QueryRow(ctx, moveItemToAnotherList,
		arg.ListItemUuid,
		arg.TargetListUuid,
		arg.NewPosition,
		arg.PrevPosition,
		arg.NextPosition,
	)
	var i MoveItemToAnotherListRow
	err := row.Scan(
		&i.Uuid,
		&i.Position,
		&i.PrevItemID,
		&i.NextItemID,
	)
	return i, err
}
