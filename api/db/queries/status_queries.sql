-- name: AddStatus :one
INSERT INTO
    statuses (user_id, label)
VALUES
    (
        (
            SELECT
                user_id
            FROM
                users
            WHERE
                users.uuid = @user_uuid
        ),
        @status_label
    )
RETURNING
    uuid,
    label,
    created_date;

-- name: GetStatus :one
SELECT
    uuid,
    label,
    created_date
FROM
    statuses
WHERE
    uuid = @status_uuid
LIMIT
    1;

-- name: GetStatusesForUser :many
SELECT
    s.uuid,
    s.label,
    s.created_date
FROM
    statuses s
        JOIN users u ON u.user_id = s.user_id
WHERE
    u.uuid = @user_uuid
ORDER BY
    s.created_date
LIMIT
    @page_size
    OFFSET
    @page;

-- name: GetAllStatuses :many
SELECT
    uuid,
    label,
    created_date
FROM
    statuses
ORDER BY
    created_date
LIMIT
    @page_size
    OFFSET
    @page;