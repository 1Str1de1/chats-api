-- +goose Up
CREATE TABLE IF NOT EXISTS chats
(
    id         BIGSERIAL PRIMARY KEY,
    title      TEXT UNIQUE,
    created_at TIMESTAMP
    CHECK ( length(title) > 0 AND length(title) < 200 )
);

-- +goose Down
DROP TABLE chats;