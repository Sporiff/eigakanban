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
