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

-- name: AddBoardItem :one
INSERT INTO
    board_items (item_id, board_id, user_id)
VALUES
    ($1, $2, $3)
RETURNING
    board_item_id,
    item_id,
    board_id;

-- name: GetBoardItem :one
SELECT
    board_item_id,
    item_id,
    board_item_id
FROM
    board_items
WHERE
    board_item_id = $1
LIMIT
    1;

-- name: GetBoardItemsForUser :many
SELECT
    board_item_id,
    item_id,
    board_item_id
FROM
    board_items
WHERE
    user_id = $1
ORDER BY
    created_date
LIMIT
    $2
    OFFSET
    $3;

-- name: GetBoardItemsForBoard :many
SELECT
    board_item_id,
    item_id,
    board_item_id
FROM
    board_items
WHERE
    board_id = $1
ORDER BY
    created_date
LIMIT
    $2
    OFFSET
    $3;

-- name: DeleteBoardItem :exec
DELETE FROM board_items
WHERE
    board_item_id = $1;

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