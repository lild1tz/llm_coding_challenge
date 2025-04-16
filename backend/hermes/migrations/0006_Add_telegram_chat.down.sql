DELETE FROM hermes_data.tables
WHERE message_id in (
    SELECT id FROM hermes_data.messages
    WHERE chat_id = (SELECT id FROM hermes_data.chat WHERE chat_name = 'tg@-4759347163')
);

DELETE FROM hermes_data.messages
WHERE chat_id = (SELECT id FROM hermes_data.chat WHERE chat_name = 'tg@-4759347163');

DELETE FROM hermes_data.chat
WHERE chat_name = 'tg@-4759347163';

DELETE FROM hermes_data.chat_context
WHERE name = 'тестовый чат telegram';
