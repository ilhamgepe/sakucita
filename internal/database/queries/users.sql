-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByIDWithRoles :one
SELECT 
    u.*,
    COALESCE(
      json_agg(
        json_build_object(
          'id', r.id,
          'name', r.name
        )
      ) FILTER (WHERE r.id IS NOT NULL),'[]'
    )::JSONB AS roles
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.id = $1
GROUP BY u.id;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL;

-- name: CreateUser :one
INSERT INTO users (
  email,
  phone,
  name,
  nickname
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

