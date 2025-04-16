CREATE SCHEMA hermes_data

CREATE TABLE hermes_data.worker (
    id SERIAL,
    name VARCHAR(1023),

    PRIMARY KEY (id)
);

CREATE TABLE hermes_data.chat_context (
    id SERIAL,
    created_at TIMESTAMP NOT NULL,
    name VARCHAR(1023) NOT NULL,

    PRIMARY KEY (id)
);

CREATE TABLE hermes_data.chat (
    id SERIAL,
    type VARCHAR(255) NOT NULL,
    chat_name VARCHAR(1023) NOT NULL UNIQUE,
    chat_context_id INTEGER NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (chat_context_id) REFERENCES hermes_data.chat_context
);

CREATE TABLE hermes_data.whatsapp (
    id SERIAL,
    whatsapp_id VARCHAR(1023) NOT NULL UNIQUE,
    worker_id INTEGER NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (worker_id) REFERENCES hermes_data.worker
);

CREATE TABLE hermes_data.telegram (
    id SERIAL,
    telegram_id VARCHAR(1023) NOT NULL UNIQUE,
    worker_id INTEGER NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (worker_id) REFERENCES hermes_data.worker
);

CREATE TABLE hermes_data.verbiage (
    id SERIAL,
    worker_id INTEGER NOT NULL,
    chat_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    content TEXT NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (worker_id) REFERENCES hermes_data.worker,
    FOREIGN KEY (chat_id) REFERENCES hermes_data.chat
);

CREATE TABLE hermes_data.messages (
    id SERIAL,
    worker_id INTEGER NOT NULL,
    chat_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    content TEXT NOT NULL,
    role VARCHAR(1023) NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (worker_id) REFERENCES hermes_data.worker,
    FOREIGN KEY (chat_id) REFERENCES hermes_data.chat
);

CREATE TABLE hermes_data.images (
    id SERIAL,
    message_id INTEGER NOT NULL,
    image_url VARCHAR(1023) NOT NULL,
    
    PRIMARY KEY (id),
    FOREIGN KEY (message_id) REFERENCES hermes_data.messages
);

CREATE TABLE hermes_data.tables (
    id SERIAL,
    created_at TIMESTAMP NOT NULL,
    data JSONB NOT NULL,
    message_id INTEGER NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (message_id) REFERENCES hermes_data.messages
);

CREATE TABLE hermes_data.listener (
    id SERIAL,
    worker_id INTEGER NOT NULL,
    chat_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (worker_id) REFERENCES hermes_data.worker,
    FOREIGN KEY (chat_id) REFERENCES hermes_data.chat
);

CREATE TABLE hermes_data.report (
    id SERIAL,
    chat_context_id INTEGER NOT NULL,
    started_at TIMESTAMP NOT NULL,
    last_updated_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP,
    
    PRIMARY KEY (id),
    FOREIGN KEY (chat_context_id) REFERENCES hermes_data.chat_context
);
