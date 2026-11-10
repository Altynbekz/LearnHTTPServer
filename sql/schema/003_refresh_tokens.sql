-- +goose Up
create table if not exists refresh_tokens(
    token text primary key,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    user_id uuid not null references users(id) on delete cascade,
    expires_at timestamp not null,
    revoked_at timestamp
);


-- +goose Down
DROP TABLE refresh_tokens;
