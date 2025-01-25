-- name: AddBoard :one
INSERT INTO
    boards (name, description, user_id)
VALUES
    ($1, $2, $3)
RETURNING
    board_id,
    name,
    description,
    created_date;

-- name: GetBoardById :one
SELECT
    name,
    description
FROM
    boards
WHERE
    board_id = $1
LIMIT
    1;

-- name: GetBoardsForUser :many
SELECT
    name,
    description
FROM
    boards
WHERE
    user_id = $1
ORDER BY
    board_id
LIMIT
    $2
    OFFSET
    $3;

-- name: GetAllBoards :many
SELECT
    name,
    description
FROM
    boards
ORDER BY
    board_id
LIMIT
    $1
    OFFSET
    $2;

-- name: UpdateBoard :one
UPDATE boards
SET
    name = $2,
    description = $3
WHERE
    board_id = $1
RETURNING
    board_id,
    name,
    description;

-- name: DeleteBoard :exec
DELETE FROM boards
WHERE
    board_id = $1;