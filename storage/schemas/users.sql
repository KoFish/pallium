CREATE TABLE IF NOT EXISTS
users(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id TEXT,
  password TEXT,
  salt TEXT,
  creation_ts INTEGER,
  UNIQUE(user_id)
);
