INSERT INTO hermes_data.chat_context (created_at, name) 
VALUES (NOW(), 'тестовый чат');

INSERT INTO hermes_data.chat (type, chat_name, chat_context_id) 
VALUES ('whatsapp', '120363418020586042@g.us', (SELECT id FROM hermes_data.chat_context WHERE name = 'тестовый чат'));

INSERT INTO hermes_data.worker (name) 
VALUES ('тестовый работник');

INSERT INTO hermes_data.whatsapp (whatsapp_id, worker_id) 
VALUES ('79853651515@s.whatsapp.net', (SELECT id FROM hermes_data.worker WHERE name = 'тестовый работник'));

INSERT INTO hermes_data.listener (chat_id, worker_id, created_at) 
VALUES ((SELECT id FROM hermes_data.chat WHERE chat_name = '120363418020586042@g.us'), (SELECT id FROM hermes_data.worker WHERE name = 'тестовый работник'), NOW());
