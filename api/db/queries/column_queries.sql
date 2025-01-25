-- name: AddColumn :one
WITH max_position AS (
    SELECT COALESCE(MAX(position), 0) AS max_position
    FROM kbcolumns
    WHERE board_id = $2
)
INSERT INTO kbcolumns (name, board_id, user_id, position)
VALUES
    ($1, $2, $3, (SELECT max_position + 1 FROM max_position))
RETURNING *;

-- name: AddColumnAtPosition :one
WITH shift_columns AS (
    UPDATE kbcolumns
        SET position = position + 1
        WHERE board_id = $2 AND position >= $3
        RETURNING column_id, position
)
INSERT INTO kbcolumns (name, board_id, user_id, position)
VALUES ($1, $2, $3, $3)
RETURNING *;

-- name: GetColumn :one
SELECT * FROM kbcolumns
WHERE column_id = $1;

-- name: GetColumns :many
SELECT * FROM kbcolumns
ORDER BY
    created_date
LIMIT
    $1
    OFFSET
    $2;

-- name: GetColumnsForBoard :many
SELECT * FROM kbcolumns
WHERE board_id = $1
ORDER BY
    position
LIMIT
    $2
    OFFSET
    $3;

-- name: GetColumnsForUser :many
SELECT * FROM kbcolumns
WHERE user_id = $1
ORDER BY
    created_date
LIMIT
    $2
    OFFSET
    $3;

-- name: UpdateColumn :one
UPDATE kbcolumns
SET
    name = $2
WHERE
    column_id = $1
RETURNING
    *;

-- name: MoveColumn :one
WITH neighbors AS (
    SELECT column_id, position
    FROM kbcolumns
    WHERE board_id = $2 AND column_id != $1
    ORDER BY position
),
     updated_positions AS (
         UPDATE kbcolumns
             SET position = CASE
                                WHEN kbcolumns.column_id = $1 THEN $2  -- New position for the moved column
                                WHEN kbcolumns.column_id = $3 THEN $4  -- Update neighbor's position
                                ELSE kbcolumns.position
                 END
             WHERE kbcolumns.board_id = $5
             RETURNING kbcolumns.column_id, kbcolumns.position
     )
UPDATE kbcolumns
SET position = kbcolumns.position + CASE
                                        WHEN kbcolumns.column_id = $1 THEN -1  -- Move column left
                                        WHEN kbcolumns.column_id = $2 THEN +1  -- Move column right
                                        ELSE 0
    END
WHERE kbcolumns.column_id IN (SELECT updated_positions.column_id FROM updated_positions)
RETURNING kbcolumns.column_id, kbcolumns.position;

-- name: DeleteColumn :exec
WITH deleted_column AS (
    DELETE FROM kbcolumns as kc
        WHERE kc.column_id = $1
        RETURNING kc.position, kc.board_id
)
UPDATE kbcolumns
SET position = position - 1
WHERE board_id = (SELECT board_id FROM deleted_column)
  AND position > (SELECT position FROM deleted_column);