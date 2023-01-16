-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users
    (username,first_name,last_name,email,password) 
    VALUES($1,$2,$3,$4,$5) RETURNING *;

-- name: MakeSuperuser :one
UPDATE users SET is_superuser=true WHERE id=$1 RETURNING *;

-- name: UpdateUser :one
UPDATE users SET 
    username=$1, 
    first_name=$2, 
    last_name=$3, 
    email=$4
WHERE id=$5 RETURNING *;

-- name: UpdatePassword :exec
UPDATE users SET password=$1 WHERE id=$2;

-- name: GetUserById :one
SELECT * FROM users WHERE id=$1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email=$1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE email=$1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id=$1;

-- name: ActivateUserAccount :exec
UPDATE users SET active=true WHERE id=$1;

-- name: DeactivateUserAccount :exec
UPDATE users SET active=false WHERE id=$1;

-- name: UpdateLastLogin :one
UPDATE users SET last_login=CURRENT_TIMESTAMP WHERE id=$1 RETURNING *;

-- name: GetUsers :many
SELECT id, username, first_name, last_name FROM users 
WHERE id IN (@ids::bigint[]);
