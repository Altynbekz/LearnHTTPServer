-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: FindUser :one
select * from users where id = $1 limit 1;

-- name: DeleteAll :exec
delete from users;
