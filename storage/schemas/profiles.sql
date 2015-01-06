CREATE TABLE IF NOT EXISTS
profiles(
    user_id INTEGER NOT NULL,
    display_name TEXT,
    avatar_url TEXT,
    UNIQUE(user_id),
    FOREIGN KEY(user_id) REFERENCES users(id)
);
