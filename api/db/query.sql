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

-- name: AddUser :one
INSERT INTO
    users (username, hashed_password, email, full_name, bio)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING
    user_id,
    username,
    full_name,
    bio,
    created_date;

-- name: GetUserById :one
SELECT
    username,
    full_name,
    bio
FROM
    users
WHERE
    user_id = $1
LIMIT
    1;

-- name: GetAllUsers :many
SELECT
    username,
    full_name,
    bio
FROM
    users
ORDER BY
    user_id
LIMIT
    $1
    OFFSET
    $2;

-- name: UpdateUserDetails :one
UPDATE users
SET
    username = $2,
    full_name = $3,
    bio = $4
WHERE
    user_id = $1
RETURNING
    username,
    full_name,
    bio;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
    user_id = $1;

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

-- name: AddReview :one
INSERT INTO
    reviews (item_id, user_id, content)
VALUES
    ($1, $2, $3)
RETURNING
    review_id,
    user_id,
    item_id,
    content,
    created_date;

-- name: GetReview :one
SELECT
    *
FROM
    reviews
WHERE
    review_id = $1
LIMIT
    1;

-- name: GetReviewsForUser :many
SELECT
    *
FROM
    reviews
WHERE
    user_id = $1
ORDER BY
    created_date
LIMIT
    $2
    OFFSET
    $3;

-- name: GetReviewsForItem :many
SELECT
    *
FROM
    reviews
WHERE
    item_id = $1
ORDER BY
    created_date
LIMIT
    $2
    OFFSET
    $3;

-- name: UpdateReview :one
UPDATE reviews
SET
    content = $2
WHERE
    review_id = $1
RETURNING
    review_id,
    item_id,
    content;

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE
    review_id = $1;

-- name: AddStatus :one
INSERT INTO
    statuses (user_id, label)
VALUES
    ($1, $2)
RETURNING
    status_id,
    user_id,
    label,
    created_date;

-- name: GetStatus :one
SELECT * FROM statuses
WHERE
    status_id = $1
LIMIT
    1;

-- name: GetStatusesForUser :many
SELECT * FROM statuses
WHERE
    user_id = $1
ORDER BY
    created_date
LIMIT
    $2
OFFSET
    $3;

-- name: GetAllStatuses :many
SELECT * FROM statuses
ORDER BY
    created_date
LIMIT
    $1
    OFFSET
    $2;