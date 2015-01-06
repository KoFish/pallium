CREATE TABLE IF NOT EXISTS
rooms (
    event_id TEXT,
    room_id TEXT PRIMARY KEY NOT NULL,
    is_public INTEGER,
    creator TEXT
);

CREATE TABLE IF NOT EXISTS
room_names (
    event_id TEXT NOT NULL,
    room_id TEXT NOT NULL,
    name TEXT NOT NULL,
    CONSTRAINT uniqueness UNIQUE (room_id)
);

CREATE TABLE IF NOT EXISTS
room_topics (
    event_id TEXT NOT NULL,
    room_id TEXT NOT NULL,
    topic TEXT NOT NULL,
    CONSTRAINT uniqueness UNIQUE (room_id)
);

CREATE TABLE IF NOT EXISTS
room_join_rules (
    event_id TEXT NOT NULL,
    room_id TEXT NOT NULL,
    join_rule TEXT NOT NULL,
    CONSTRAINT uniqueness UNIQUE (room_id)
);

CREATE TABLE IF NOT EXISTS
room_power_levels (
    event_id TEXT NOT NULL,
    room_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    level INTEGER NOT NULL,
    CONSTRAINT uniqueness UNIQUE (room_id, user_id)
);

CREATE TABLE IF NOT EXISTS
room_default_levels (
    event_id TEXT NOT NULL,
    room_id TEXT NOT NULL,
    level INTEGER NOT NULL,
    CONSTRAINT uniqueness UNIQUE (room_id)
);

CREATE TABLE IF NOT EXISTS
room_add_state_levels (
    event_id TEXT NOT NULL,
    room_id TEXT NOT NULL,
    level INTEGER NOT NULL,
    CONSTRAINT uniqueness UNIQUE (room_id)
);

CREATE TABLE IF NOT EXISTS
room_send_event_levels (
    event_id TEXT NOT NULL,
    room_id TEXT NOT NULL,
    level INTEGER NOT NULL,
    CONSTRAINT uniqueness UNIQUE (room_id)
);

CREATE TABLE IF NOT EXISTS
room_ops_levels (
    event_id TEXT NOT NULL,
    room_id TEXT NOT NULL,
    ban_level INTEGER NOT NULL,
    kick_level INTEGER NOT NULL,
    redact_level INTEGER NOT NULL,
    CONSTRAINT uniqueness UNIQUE (room_id)
);

CREATE TABLE IF NOT EXISTS
room_aliases (
    event_id TEXT NOT NULL,
    room_alias TEXT NOT NULL,
    room_id TEXT NOT NULL,
    server TEXT NOT NULL,
    CONSTRAINT uniqueness UNIQUE (room_alias)
);

CREATE TABLE IF NOT EXISTS
room_memberships (
    event_id TEXT NOT NULL,
    room_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    membership TEXT,
    CONSTRAINT uniqueness UNIQUE (room_id, user_id)
);
