// Copyright 2014 Krister Svanlund
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	c "github.com/KoFish/pallium/config"
	m "github.com/KoFish/pallium/matrix"
	"time"
)

const events_table = `
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

CREATE INDEX IF NOT EXISTS room_id_depth_room_id ON room_id_depth (room_id)
`

func getCurrentroom_idEvent(db DBI, room_id m.RoomID) (m.EventID, error) {
	var (
		event_id_string string
	)
	row := db.QueryRow("SELECT (event_id) FROM events WHERE room_id=? ORDER BY e.topological_ordering DESC LIMIT 1", room_id.String())
	if err := row.Scan(&event_id_string); err != nil {
		return m.EventID{}, err
	}
	return m.ParseEventID(event_id_string)
}

func getCurrentState(db DBI, room_id m.RoomID, event_type, state_key string) (m.EventID, error) {
	var (
		event_id_string string
	)
	row := db.QueryRow(`SELECT (s.event_id) FROM state_events AS s
                       INNER JOIN events AS e
                       ON e.event_id == state_events.event_id
                       WHERE s.room_id=? AND s.type=? AND s.state_key=? LIMIT 1 ORDER BY e.topological_ordering DESC LIMIT 1`)
	if err := row.Scan(&event_id_string); err != nil {
		return m.EventID{}, err
	}
	return m.ParseEventID(event_id_string)
}

func updateroom_idDepth(db DBI, room_id m.RoomID, depth int64) error {
	if _, err := db.Exec(`INSERT OR REPLACE
        INTO room_id_depth (room_id, min_depth)
        VALUES (?, ?)`, room_id.String(), depth); err != nil {
		return err
	}
	return nil
}

func getroom_idDepth(db DBI, room_id m.RoomID) (int64, error) {
	var room_id_depth int64
	row := db.QueryRow("SELECT min_depth FROM room_id_depth WHERE room_id=?", room_id.String())
	if err := row.Scan(&room_id_depth); err != nil {
		room_id_depth = 0
	}
	return room_id_depth, nil
}

func createEventEdge(db DBI, from_event, to_event m.EventID, room_id m.RoomID) error {
	if _, err := db.Exec(`INSERT OR FAIL
        INTO event_edges (event_id, prev_event_id, room_id)
        VALUES (?, ?, ?)`, to_event.String(), from_event.String(), room_id.String()); err != nil {
		return err
	}
	return nil
}

func GetPowerLevel(db DBI, room_id m.RoomID, state_key string) (int64, error) {
	var power_level int64
	row := db.QueryRow("SELECT power_level FROM state_power_levels WHERE room_id=? AND state_key=?", room_id.String(), state_key)
	if err := row.Scan(&power_level); err != nil {
		return c.Config.DefaultPowerLevel, err
	}
	return power_level, nil
}

func CreateEvent(tx *sql.Tx, event_id m.EventID, user_id m.UserID, room_id m.RoomID, event_type string, content map[string]interface{}) error {
	var (
		content_json string
		ts           int64 = time.Now().Unix()
	)
	if content_bytes, err := json.Marshal(content); err != nil {
		return err
	} else {
		content_json = string(content_bytes)
	}
	curr_event, ce_err := getCurrentroom_idEvent(tx, room_id)
	room_id_depth, _ := getroom_idDepth(tx, room_id)
	event_depth := room_id_depth + 1
	sqlstmt := `INSERT OR FAIL
    INTO events (
        topological_ordering,
        event_id,
        type,
        ts,
        room_id,
        content_json,
        user_id,
        processed,
        outlier)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.Exec(sqlstmt,
		event_depth,
		event_id.String(),
		event_type,
		ts,
		room_id.String(),
		content_json,
		user_id.String(),
		true, false)
	if err != nil {
		return err
	}
	if rowsaffected, err := result.RowsAffected(); err != nil {
		return err
	} else {
		if rowsaffected != 1 {
			return fmt.Errorf("matrix: no rows affected by insert statement")
		}
	}
	if ce_err != nil {
		if err := createEventEdge(tx, curr_event, event_id, room_id); err != nil {
			return err
		}
	}
	err = updateroom_idDepth(tx, room_id, event_depth)
	if err == nil {
		fmt.Printf("storage: created %v event %v by %v\n", event_type, event_id.String(), user_id.String())
	}
	return err
}

func CreateStateEvent(tx *sql.Tx, event_id m.EventID, user_id m.UserID, room_id m.RoomID, event_type, state_key string, content map[string]interface{}) error {
	prev_state, cs_err := getCurrentState(tx, room_id, event_type, state_key)
	if err := CreateEvent(tx, event_id, user_id, room_id, event_type, content); err != nil {
		return err
	}
	sqlstmt := `INSERT OR FAIL
    INTO state_events (
        event_id,
        room_id,
        type,
        state_key,
        prev_state_id)
    VALUES (?, ?, ?, ?, ?)`
	if _, err := tx.Exec(sqlstmt,
		event_id.String(),
		room_id.String(),
		event_type,
		state_key,
		sql.NullString{prev_state.String(), cs_err == nil}); err != nil {
		return err
	}
	return nil
}

func NewEvent(tx *sql.Tx, user_id m.UserID, room_id m.RoomID, event_type string, content map[string]interface{}) (m.EventID, error) {
	event_id, err := m.GenerateEventID()
	if err != nil {
		return event_id, err
	}
	return event_id, CreateEvent(tx, event_id, user_id, room_id, event_type, content)
}

func NewStateEvent(tx *sql.Tx, user_id m.UserID, room_id m.RoomID, event_type, state_key string, content map[string]interface{}) (m.EventID, error) {
	event_id, err := m.GenerateEventID()
	if err != nil {
		return event_id, err
	}
	return event_id, CreateStateEvent(tx, event_id, user_id, room_id, event_type, state_key, content)
}
