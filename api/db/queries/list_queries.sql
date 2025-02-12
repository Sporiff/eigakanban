-- name: AddList :one
INSERT INTO
    lists (name, user_id)
VALUES
    (
        @name,
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
    created_date;

-- name: GetListByUuid :one
SELECT
    uuid,
    name
FROM
    lists
WHERE
    uuid = @list_uuid;

-- name: GetListsByUser :many
SELECT
    l.uuid,
    l.name
FROM
    lists l
        JOIN users u ON u.user_id = l.user_id
WHERE
    u.uuid = @user_uuid
ORDER BY
    l.created_date
LIMIT
    @page_size
    OFFSET
    @page;

-- name: UpdateList :one
UPDATE lists
SET
    name = @list_name
WHERE
    uuid = @list_uuid
RETURNING
    uuid,
    name;

-- name: DeleteList :exec
DELETE FROM lists
WHERE
    uuid = @list_uuid;