-- name: RevokeAllSessionsByUserID :exec
UPDATE sessions SET revoked = TRUE WHERE user_id = $1;

-- name: RevokeSessionByID :exec
UPDATE sessions SET revoked = TRUE WHERE id = $1 AND user_id = $2;

-- name: CreateSession :one
INSERT INTO sessions (
  user_id,
  device_id,
  refresh_token_id,
  expires_at,
  meta
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;