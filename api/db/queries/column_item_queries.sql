-- name: AddColumnItemAtPosition :one
WITH moved_items AS (
    UPDATE column_items
        SET position = position + 1
        WHERE column_id = $1 AND position >= $2
        RETURNING column_item_id, column_id, item_id, position
)
INSERT INTO column_items (column_id, item_id, user_id, position)
VALUES ($1, $2, $3, $2)
RETURNING column_item_id, column_id, item_id, position;

-- name: MoveItemToColumn :one
WITH moved_items AS (
    UPDATE column_items AS ci
        SET position = position + 1
        WHERE ci.column_id = $2 AND ci.position >= (SELECT ci2.position FROM column_items AS ci2 WHERE ci2.column_item_id = $1)
        RETURNING ci.column_item_id, ci.column_id, ci.item_id, ci.position
)
UPDATE column_items AS ci
SET column_id = $2
WHERE ci.column_item_id = $1
RETURNING ci.column_item_id, ci.column_id, ci.item_id, ci.position;

-- name: MoveColumnItemUp :one
WITH current_position AS (
    SELECT position
    FROM column_items AS ci
    WHERE ci.column_item_id = $1
),
     items_to_shift AS (
         SELECT ci.column_item_id, ci.position
         FROM column_items AS ci
         WHERE ci.column_id = $2
           AND ci.position = (SELECT cp.position FROM current_position AS cp) - 1
     )
UPDATE column_items AS ci
SET position = CASE
                   WHEN ci.column_item_id = $1 THEN ci.position - 1
                   WHEN ci.column_item_id IN (SELECT its.column_item_id FROM items_to_shift AS its) THEN ci.position + 1
                   ELSE ci.position
    END
WHERE ci.column_item_id IN (SELECT its.column_item_id FROM items_to_shift AS its)
   OR ci.column_item_id = $1
RETURNING ci.column_item_id, ci.column_id, ci.item_id, ci.position;

-- name: MoveColumnItemDown :one
WITH current_position AS (
    SELECT position
    FROM column_items AS ci
    WHERE ci.column_item_id = $1
),
     items_to_shift AS (
         SELECT ci.column_item_id, ci.position
         FROM column_items AS ci
         WHERE ci.column_id = $2
           AND ci.position = (SELECT cp.position FROM current_position AS cp) + 1
     )
UPDATE column_items AS ci
SET position = CASE
                   WHEN ci.column_item_id = $1 THEN ci.position + 1
                   WHEN ci.column_item_id IN (SELECT its.column_item_id FROM items_to_shift AS its) THEN ci.position - 1
                   ELSE ci.position
    END
WHERE ci.column_item_id IN (SELECT its.column_item_id FROM items_to_shift AS its)
   OR ci.column_item_id = $1
RETURNING ci.column_item_id, ci.column_id, ci.item_id, ci.position;

-- name: DeleteColumnItem :exec
WITH deleted_item AS (
    DELETE FROM column_items as ci
        WHERE ci.column_item_id = $1
        RETURNING ci.position, ci.column_id
)
UPDATE column_items
SET position = position - 1
WHERE column_id = (SELECT di.column_id FROM deleted_item AS di)
  AND position > (SELECT di.position FROM deleted_item AS di);