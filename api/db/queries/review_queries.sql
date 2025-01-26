-- name: AddReview :one
INSERT INTO
    reviews (item_id, user_id, content)
VALUES
    (
        (
            SELECT
                item_id
            FROM
                items
            WHERE
                items.uuid = @item_uuid
        ),
        (
            SELECT
                user_id
            FROM
                users
            WHERE
                users.uuid = @user_uuid
        ),
        @content
    )
RETURNING
    uuid,
    content,
    created_date;

-- name: GetReview :one
SELECT
    uuid,
    content,
    created_date
FROM
    reviews
WHERE
    uuid = @review_uuid
LIMIT
    1;

-- name: GetReviewsForUser :many
SELECT
    r.uuid,
    r.content,
    r.created_date
FROM
    reviews r
        JOIN users u ON u.user_id = r.user_id
WHERE
    u.uuid = @user_uuid
ORDER BY
    r.created_date
LIMIT
    @page_size
    OFFSET
    @page;

-- name: GetReviewsForItem :many
SELECT
    r.uuid,
    r.content,
    r.created_date
FROM
    reviews r
        JOIN items i ON i.item_id = r.item_id
WHERE
    i.uuid = @item_uuid
ORDER BY
    r.created_date
LIMIT
    @page_size
    OFFSET
    @page;

-- name: UpdateReview :one
UPDATE reviews
SET
    content = @content
WHERE
    uuid = @review_uuid
RETURNING
    uuid,
    content;

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE
    uuid = @review_uuid;