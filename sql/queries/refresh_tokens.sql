-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (
    token,
    user_id,
    expires_at,
    revoked_at
)
VALUES ($1, $2, $3, null);

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE token = $1
  AND revoked_at IS NULL;


-- name: FindRefreshToken :one
select * from refresh_tokens where token = $1 limit 1;

-- name: ExistsByRefreshToken :one
select exists(select 1 from refresh_tokens where token = $1);