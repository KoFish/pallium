CREATE TABLE IF NOT EXISTS
access_tokens(
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL,
  token TEXT NOT NULL,
  last_used INTEGER,
  created INTEGER,
  FOREIGN KEY(user_id) REFERENCES users(id),
  UNIQUE(token)
);
