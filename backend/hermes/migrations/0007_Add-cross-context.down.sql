DELETE FROM hermes_data.image WHERE message_id IN (
    SELECT id FROM hermes_data.message WHERE chat_id IN (
        SELECT id FROM hermes_data.chat WHERE chat_name = '120363402118486398@g.us' OR chat_name = 'tg@-4663851323'
    )
);

DELETE FROM hermes_data.verbiage WHERE chat_id IN (
    SELECT id FROM hermes_data.chat WHERE chat_name = '120363402118486398@g.us' OR chat_name = 'tg@-4663851323'
);

DELETE FROM hermes_data.message WHERE chat_id IN (
    SELECT id FROM hermes_data.chat WHERE chat_name = '120363402118486398@g.us' OR chat_name = 'tg@-4663851323'
);

DELETE FROM hermes_data.listener WHERE chat_id IN (
    SELECT id FROM hermes_data.chat WHERE chat_name = '120363402118486398@g.us'
);

DELETE FROM hermes_data.chat WHERE chat_name = '120363402118486398@g.us';
DELETE FROM hermes_data.chat WHERE chat_name = 'tg@-4663851323';

DELETE FROM hermes_data.chat_context WHERE name = 'вотс + телеграм';