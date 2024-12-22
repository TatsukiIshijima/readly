-- name: GetAllUsers :many
SELECT *
FROM users
ORDER BY id LIMIT $1
OFFSET $2;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (name,
                   email,
                   hashed_password)
VALUES ($1,
        $2,
        $3) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name            = $2,
    email           = $3,
    hashed_password = $4,
    updated_at      = now()
WHERE id = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;