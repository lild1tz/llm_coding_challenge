INSERT INTO hermes_data.chat_context (created_at, name)
VALUES (NOW(), 'тестовый чат telegram');

INSERT INTO hermes_data.chat (type, chat_name, chat_context_id)
VALUES ('telegram', 'tg@-4759347163', (SELECT id FROM hermes_data.chat_context WHERE name = 'тестовый чат telegram'));
