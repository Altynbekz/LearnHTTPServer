-- name: GetChirps :one
select * from chirps where id = $1 limit 1;

-- name: ExistsByChirpId :one
select exists (select 1 from chirps where id = $1);

-- name: CreateChirp :one
insert into chirps(id, created_at, updated_at, body, user_id)
values(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
returning *;



-- name: GetAllChirps :many
select * from chirps order by created_at;


-- name: DeleteChirp :exec
delete from chirps where id = $1;