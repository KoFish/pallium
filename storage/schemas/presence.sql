CREATE TABLE IF NOT EXISTS
presence (
    user_id INTEGER NOT NULL,
    state INTEGER,
    status_msg TEXT,
    mtime INTEGER,
    FOREIGN KEY(user_id) REFERENCES users(id)
);
