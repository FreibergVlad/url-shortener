-- name: CreateUser :exec
INSERT INTO users (
    id,
    password_hash,
    email,
    full_name,
    role_id,
    created_at

) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;