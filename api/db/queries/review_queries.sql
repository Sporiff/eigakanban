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