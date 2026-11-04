-- +goose Up
CREATE TABLE IF NOT EXISTS chirps (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    body TEXT NOT NULL,
    -- Removing a user also removes their chirps.
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- An earlier local chirps table allowed NULLs; make the constraint explicit there too.
ALTER TABLE chirps
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET NOT NULL,
    ALTER COLUMN body SET NOT NULL;

-- +goose Down
DROP TABLE chirps;
