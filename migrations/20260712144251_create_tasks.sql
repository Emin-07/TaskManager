-- +goose Up
CREATE TABLE IF NOT EXISTS tasks (
                                     id        SERIAL PRIMARY KEY,
                                     title     TEXT NOT NULL,
                                     text      TEXT NOT NULL DEFAULT '',
                                     priority  INTEGER NOT NULL DEFAULT 0,
                                     created   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     expires   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '7 days',
                                     user_id   INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_tasks_user_id ON tasks(user_id);

CREATE INDEX idx_tasks_title ON tasks(title);

-- +goose Down
DROP INDEX IF EXISTS idx_tasks_user_id;
DROP INDEX IF EXISTS idx_tasks_title;

DROP TABLE IF EXISTS tasks;
