-- +goose Up
create table users(
    id UUID primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    email text not null,
    hashed_password text not null,
    is_chirpy_red BOOLEAN NOT NULL DEFAULT FALSE
);


-- +goose Down
drop table users;
