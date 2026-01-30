-- +goose Up
CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    text TEXT,
    chat_id INT,
    created_at TIMESTAMP,
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    CHECK ( length(text) > 0 AND length(text) < 5000 )
);

-- +goose Down
DROP TABLE messages;