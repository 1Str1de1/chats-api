-- +goose Up
CREATE TABLE IF NOT EXISTS chats (
    id BIGSERIAL PRIMARY KEY,
    title TEXT UNIQUE,
    created_at TIMESTAMP
);

-- +goose Down
DROP TABLE chats;