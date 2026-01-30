-- +goose Up
INSERT INTO messages (text, chat_id, created_at)
VALUES ('Hi!', 1, NOW()),
       ('Hi!', 2, NOW()),
       ('Hello', 3, NOW()),
       ('Hi!', 1, NOW()),
       ('Hi!', 2, NOW()),
       ('Hello', 3, NOW()),
       ('How are you?', 1, NOW()),
       ('What about that party tonight?', 2, NOW()),
       ('What is the hometask for tomorrow?', 3, NOW()),
       ('Fine, thanks', 1, NOW()),
       ('It starts at 19:00 at Beth''s house', 2, NOW()),
       ('There is no hometask just chill', 3, NOW());

-- +goose Down
TRUNCATE TABLE messages;