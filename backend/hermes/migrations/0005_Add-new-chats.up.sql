INSERT INTO hermes_data.chat_context (created_at, name)
VALUES (NOW(), 'тестовый чат 2');

INSERT INTO hermes_data.chat (type, chat_name, chat_context_id)
VALUES ('whatsapp', '120363416546569882@g.us', (SELECT id FROM hermes_data.chat_context WHERE name = 'тестовый чат 2'));

INSERT INTO hermes_data.listener (chat_id, worker_id, created_at)
VALUES ((SELECT id FROM hermes_data.chat WHERE chat_name = '120363416546569882@g.us'), (SELECT id FROM hermes_data.worker WHERE name = 'тестовый работник'), NOW());

INSERT INTO hermes_data.chat (type, chat_name, chat_context_id)
VALUES ('whatsapp', '120363398827953735@g.us', (SELECT id FROM hermes_data.chat_context WHERE name = 'тестовый чат 2'));

INSERT INTO hermes_data.listener (chat_id, worker_id, created_at)
VALUES ((SELECT id FROM hermes_data.chat WHERE chat_name = '120363398827953735@g.us'), (SELECT id FROM hermes_data.worker WHERE name = 'тестовый работник'), NOW());

