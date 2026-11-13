-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1, 
    $2, 
    FALSE
)
RETURNING *;

-- name: ExistsByUserId :one
select exists (
    select 1 from users where id = $1
);

-- name: FindUser :one
select * from users where id = $1 limit 1;

-- name: FindUserByEmail :one
select * from users where email = $1 limit 1;

-- name: DeleteAll :exec
delete from users;

-- name: UpdateUser :exec
update users set email = $1, hashed_password = $2 where id = $3;

-- name: UpdateChripStatus :exec
update users set is_chirpy_red = $1 where id = $2;