CREATE TABLE hermes_data.audios (
    id SERIAL PRIMARY KEY,
    message_id INT NOT NULL,
    audio_url VARCHAR(255) NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (message_id) REFERENCES hermes_data.messages(id)
);
