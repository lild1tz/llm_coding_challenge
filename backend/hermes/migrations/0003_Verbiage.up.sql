CREATE TABLE hermes_data.verbiage (
    id SERIAL,
    worker_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    content TEXT NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (worker_id) REFERENCES hermes_data.worker
);
