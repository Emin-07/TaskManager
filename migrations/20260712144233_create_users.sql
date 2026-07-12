-- +goose Up
CREATE TABLE IF NOT EXISTS users (
                                     id            SERIAL PRIMARY KEY,
                                     username      TEXT NOT NULL,
                                     role          TEXT NOT NULL DEFAULT 'user',
                                     email         TEXT NOT NULL UNIQUE,
                                     password_hash BYTEA NOT NULL,
                                     created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_users_username ON users(username);

-- +goose Down
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS users;
