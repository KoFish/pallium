CREATE TABLE IF NOT EXISTS
events(
    stream_ordering INTEGER PRIMARY KEY AUTOINCREMENT,
    topological_ordering INTEGER NOT NULL,
    event_id TEXT NOT NULL,
    type TEXT NOT NULL,
    ts INTEGER,
    room_id TEXT NOT NULL,
    content_json TEXT NOT NULL,
    user_id TEXT NOT NULL,
    unrecognized_keys TEXT,
    processed BOOL NOT NULL,
    outlier BOOL NOT NULL,
    CONSTRAINT ev_uniq UNIQUE (event_id)
);

CREATE INDEX IF NOT EXISTS events_event_id ON events (event_id);
CREATE INDEX IF NOT EXISTS events_stream_ordering ON events (stream_ordering);
CREATE INDEX IF NOT EXISTS events_topological_ordering ON events (topological_ordering);
CREATE INDEX IF NOT EXISTS events_room_id ON events (room_id);

CREATE TABLE IF NOT EXISTS
state_events(
    event_id TEXT NOT NULL,
    room_id TEXT NOT NULL,
    type TEXT NOT NULL,
    state_key TEXT NOT NULL,
    prev_state_id TEXT,
    CONSTRAINT event_id UNIQUE (event_id)
    CONSTRAINT prev_event_id UNIQUE (prev_state_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS state_events_event_id ON state_events (event_id);
CREATE INDEX IF NOT EXISTS state_events_room_id ON state_events (room_id);
CREATE INDEX IF NOT EXISTS state_events_type ON state_events (type);
CREATE INDEX IF NOT EXISTS state_events_state_key ON state_events (state_key);

CREATE TABLE IF NOT EXISTS
event_edges(
    event_id TEXT,
    prev_event_id TEXT,
    room_id TEXT,
    CONSTRAINT uniqueness UNIQUE (event_id, prev_event_id, room_id)
);

CREATE INDEX IF NOT EXISTS event_edges_id ON event_edges (event_id);

CREATE TABLE IF NOT EXISTS
state_power_levels(
    state_key TEXT NOT NULL,
    room_id TEXT NOT NULL,
    power_level INTEGER,
    CONSTRAINT uniqueness UNIQUE (state_key, room_id) ON CONFLICT REPLACE
);

CREATE INDEX IF NOT EXISTS state_power_levels_state_room_id ON state_power_levels (state_key, room_id);

CREATE TABLE IF NOT EXISTS
room_id_depth(
    room_id TEXT,
    min_depth INTEGER NOT NULL,
    CONSTRAINT uniqueness UNIQUE (room_id) ON CONFLICT REPLACE
);

CREATE INDEX IF NOT EXISTS room_id_depth_room_id ON room_id_depth (room_id);
