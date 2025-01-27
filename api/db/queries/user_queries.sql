-- name: AddUser :one
INSERT INTO
    users (username, hashed_password, email, full_name, bio)
VALUES
    (@username, @hashed_password, @email, @full_name, @bio)
RETURNING
    uuid,
    username,
    full_name,
    bio,
    created_date;

-- name: GetUserByUuid :one
SELECT
    uuid,
    username,
    full_name,
    bio
FROM
    users
WHERE
    uuid = @user_uuid
LIMIT
    1;

-- name: GetAllUsers :many
SELECT
    uuid,
    username,
    full_name,
    bio
FROM
    users
ORDER BY
    created_date
LIMIT
    @page_size
    OFFSET
    @page;

-- name: UpdateUserDetails :one
UPDATE users
SET
    username = COALESCE(@new_username, username),
    full_name = COALESCE(@new_name, full_name),
    bio = COALESCE(@new_bio, bio)
WHERE
    uuid = @user_uuid
RETURNING
    uuid,
    username,
    full_name,
    bio;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
    uuid = @user_uuid;