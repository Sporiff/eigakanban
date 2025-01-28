-- name: AddBoard :one
INSERT INTO
    boards (name, description, user_id)
VALUES
    (
        @board_name,
        @board_description,
        (
            SELECT
                user_id
            FROM
                users
            WHERE
                users.uuid = @user_uuid
        )
    )
RETURNING
    uuid,
    name,
    description,
    created_date;

-- name: GetBoardByUuid :one
SELECT
    uuid,
    name,
    description
FROM
    boards
WHERE
    uuid = @board_uuid
LIMIT
    1;

-- name: GetBoardsForUser :many
SELECT
    b.uuid,
    b.name,
    b.description
FROM
    boards b
        JOIN users u ON b.user_id = u.user_id
WHERE
    u.uuid = @user_uuid
ORDER BY
    b.created_date
LIMIT
    @page_size
    OFFSET
    @page;

-- name: GetBoardsCount :one
SELECT COUNT (*)
FROM boards;

-- name: GetBoardsCountForUser :one
SELECT COUNT (*)
FROM boards
WHERE boards.user_id = (SELECT user_id FROM users WHERE users.uuid = @user_uuid);

-- name: GetAllBoards :many
SELECT
    uuid,
    name,
    description
FROM
    boards
ORDER BY
    created_date
LIMIT
    @page_size
    OFFSET
    @page;

-- name: UpdateBoard :one
UPDATE boards
SET
    name = @board_name,
    description = COALESCE(@board_description, description)
WHERE
    uuid = @board_uuid
RETURNING
    uuid,
    name,
    description;

-- name: DeleteBoard :exec
DELETE FROM boards
WHERE
    uuid = @board_uuid;