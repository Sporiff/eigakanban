-- name: AddListStatus :one
INSERT INTO
    list_statuses (list_id, status_id)
VALUES
    (
        (
            SELECT
                l.list_id
            FROM
                lists l
            WHERE
                l.uuid = @list_uuid
        ),
        (
            SELECT
                s.status_id
            FROM
                statuses s
            WHERE
                s.uuid = @status_uuid
        )
    )
RETURNING
    uuid,
    created_date;

-- name: GetListStatus :one
SELECT
    uuid,
    created_date
FROM
    list_statuses
WHERE
    uuid = @list_status_uuid
LIMIT
    1;

-- name: GetStatusesForList :many
SELECT
    ls.uuid,
    s.uuid,
    s.label,
    ls.created_date
FROM
    list_statuses ls
        JOIN statuses s ON s.status_id = ls.status_id
        JOIN lists l ON l.list_id = ls.list_id
WHERE
    l.uuid = @list_uuid
ORDER BY
    l.uuid,
    ls.created_date
LIMIT
    @page_size
    OFFSET
    @page;

-- name: DeleteListStatus :exec
DELETE FROM list_statuses
WHERE
    uuid = @list_status_uuid;