-- name: AddItem :one
INSERT INTO
    items (title)
VALUES
    (
        @title
    )
RETURNING
    uuid,
    title,
    created_date;

-- name: GetItemByUuid :one
SELECT
    uuid,
    title,
    created_date
FROM
    items
WHERE
    uuid = @item_uuid
LIMIT
    1;

-- name: GetItemsCount :one
SELECT COUNT (*)
FROM items;

-- name: GetAllItems :many
SELECT
    uuid,
    title
FROM
    items
ORDER BY
    created_date
LIMIT
    @page_size
    OFFSET
    @page;

-- name: UpdateItem :one
UPDATE items
SET title = COALESCE(@item_title, title)
WHERE uuid = @item_uuid
RETURNING
    uuid,
    title;

-- name: DeleteItem :exec
DELETE FROM items
WHERE
    uuid = @item_uuid;