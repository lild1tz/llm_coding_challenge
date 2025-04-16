
DELETE FROM hermes_data.tables
WHERE message_id in (
    SELECT id FROM hermes_data.messages
    WHERE chat_id = (SELECT id FROM hermes_data.chat WHERE chat_name = '120363416546569882@g.us')
);

DELETE FROM hermes_data.messages
WHERE chat_id = (SELECT id FROM hermes_data.chat WHERE chat_name = '120363416546569882@g.us');

DELETE FROM hermes_data.chat
WHERE chat_name = '120363416546569882@g.us';

DELETE FROM hermes_data.messages
WHERE chat_id = (SELECT id FROM hermes_data.chat WHERE chat_name = '120363398827953735@g.us');

DELETE FROM hermes_data.tables
WHERE message_id in (
    SELECT id FROM hermes_data.messages
    WHERE chat_id = (SELECT id FROM hermes_data.chat WHERE chat_name = '120363398827953735@g.us')
);

DELETE FROM hermes_data.chat_context
WHERE name = 'тестовый чат 2';
