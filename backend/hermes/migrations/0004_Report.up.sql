CREATE TABLE hermes_data.chat (
    id SERIAL,
    listener_id INTEGER NOT NULL,
    chat_id VARCHAR(1023) NOT NULL UNIQUE,
    
    PRIMARY KEY (id),
    FOREIGN KEY (listener_id) REFERENCES hermes_data.worker
);

CREATE TABLE hermes_data.report (
    id SERIAL,
    chat_id INTEGER NOT NULL,
    started_at TIMESTAMP NOT NULL,
    last_updated_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP,
    
    PRIMARY KEY (id),
    FOREIGN KEY (chat_id) REFERENCES hermes_data.chat
);
