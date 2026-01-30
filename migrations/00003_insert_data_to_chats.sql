-- +goose Up
INSERT INTO chats (title, created_at)
VALUES ('Family', NOW()),
       ('Friends', NOW()),
       ('Schoolmates', NOW());

-- +goose Down
TRUNCATE TABLE chats CASCADE;