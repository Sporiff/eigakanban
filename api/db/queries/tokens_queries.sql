-- name: AddRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expires_at)
VALUES (
        @user_id, @token, @expires_at
       )
RETURNING token, expires_at;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = @token;

-- name: GetRefreshTokenByToken :one
SELECT
    user_id,
    token,
    expires_at
FROM refresh_tokens
WHERE token = @token;