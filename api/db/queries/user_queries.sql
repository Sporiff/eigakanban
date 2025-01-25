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