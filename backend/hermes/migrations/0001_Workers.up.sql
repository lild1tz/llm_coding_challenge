CREATE SCHEMA hermes_data

CREATE TABLE hermes_data.worker (
    id SERIAL,
    whatsapp_id VARCHAR(1023) NULL,
    telegram_id VARCHAR(1023) NULL,
    name VARCHAR(1023),

    PRIMARY KEY (id),
    CONSTRAINT unique_whatsapp_id UNIQUE (whatsapp_id),
    CONSTRAINT unique_telegram_id UNIQUE (telegram_id)
);

CREATE TABLE hermes_data.messages (
    id SERIAL,
    worker_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    content TEXT NOT NULL,
    role VARCHAR(1023) NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (worker_id) REFERENCES hermes_data.worker
);

CREATE TABLE hermes_data.images (
    id SERIAL,
    message_id INTEGER NOT NULL,
    image_url VARCHAR(1023) NOT NULL,
    
    PRIMARY KEY (id),
    FOREIGN KEY (message_id) REFERENCES hermes_data.messages
);
