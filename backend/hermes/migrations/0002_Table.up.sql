CREATE TABLE hermes_data.tables (
    id SERIAL,
    created_at TIMESTAMP NOT NULL,
    data JSONB NOT NULL,

    PRIMARY KEY (id)
);
