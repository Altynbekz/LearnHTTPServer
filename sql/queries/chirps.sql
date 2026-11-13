-- name: GetChirps :one
select * from chirps where id = $1 limit 1;

-- name: GetChirpByUserId :many
select * from chirps where user_id = $1;

-- -- name: GetChirpsByUserIdAsc :many
-- select * from chirps where user_id = $1 order by user_id ASC;

-- -- name: GetChirpsByUserIdDesc :many
-- select * from chirps where user_id = $1 order by user_id DESC

-- name: GetChirpBySortUserId :many
select * from chirps where user_id = $1 order by 
    case when @sort_type::text = 'asc' then created_at end asc,
    case when @sort_type::text = 'desc' then created_at end DESC,
    created_at asc;

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
select * from chirps order by
    case when @sort_type::text = 'desc' then created_at end DESC,
    created_at asc;


-- name: DeleteChirp :exec
delete from chirps where id = $1;