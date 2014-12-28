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
	"fmt"
	c "github.com/KoFish/pallium/config"
	m "github.com/KoFish/pallium/matrix"
)

const rooms_table = `
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
)
`

//type Room struct {
//    ID m.RoomID `sql:"room_id"`
//}

type Room struct {
	Aliases       []string `json:"aliases"`
	Name          string   `json:"name"`
	JoinedMembers int      `json:"num_joined_members"`
	RoomId        string   `json:"room_id"`
	Topic         string   `json:"topic"`
	ID            m.RoomID `json:"-"`
}

func LookupRoomAlias(db DBI, room_alias m.RoomAlias) (m.RoomID, []string, error) {
	// TODO(kofish): Implement resolving room aliases to room ids
	return m.NoRoomID, []string{}, fmt.Errorf("matrix: LookupRoomAlias is not yet implemented")
}

func GetRoom(db DBI, room_id m.RoomID) (*Room, error) {
	var exists bool
	row := db.QueryRow("SELECT true FROM rooms WHERE room_id=?", room_id.String())
	if err := row.Scan(&exists); err != nil {
		return nil, fmt.Errorf("matrix: no such room exists")
	} else {
		return &Room{ID: room_id}, nil
	}
}

func GetPublicRooms(tx *sql.Tx) []Room {
	rows, err := tx.Query(
		`SELECT r.room_id as id, count(rm.room_id) as membercount
        FROM rooms r, room_memberships rm
        WHERE is_public = 1 AND r.room_id = rm.room_id`)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	var rooms []Room

	for rows.Next() {
		var (
			roomId      string
			memberCount int
		)
		err := rows.Scan(&roomId, &memberCount)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		room := Room{RoomId: roomId, JoinedMembers: memberCount}
		rooms = append(rooms, room)
	}

	return rooms
}

func CreateRoom(tx *sql.Tx, creator m.UserID, is_public bool) (*Room, error) {
	var event_id m.EventID
	room_id, err := m.GenerateRoomID()
	if err != nil {
		return nil, err
	}
	content := map[string]interface{}{"creator": creator.String()}
	if event_id, err = NewEvent(tx, creator, room_id, "m.room.create", content); err != nil {
		return nil, err
	} else {
		if _, err = tx.Exec(`INSERT
            INTO rooms (
                event_id,
                room_id,
                is_public,
                creator)
            VALUES (?, ?, ?, ?)`, event_id.String(), room_id.String(), is_public, creator.String()); err != nil {
			return nil, err
		}
	}
	return &Room{ID: room_id}, nil
}

func (r *Room) GetUserPowerLevel(db DBI, user_id m.UserID) (int64, error) {
	var powerlevel int64
	row := db.QueryRow(`SELECT level
        FROM room_power_levels
        WHERE user_id=? AND room_id=?`, user_id.String(), r.ID.String())
	if err := row.Scan(&powerlevel); err != nil {
		return powerlevel, err
	}
	return powerlevel, nil
}

func (r *Room) GetUserMembership(db DBI, user_id m.UserID) (m.RoomMembership, error) {
	var membership string
	row := db.QueryRow(`SELECT membership
        FROM room_memberships
        WHERE user_id=? AND room_id=?`, user_id.String(), r.ID.String())
	if err := row.Scan(&membership); err != nil {
		return "", err
	}
	switch membership {
	case "invite":
		return m.MEMBERSHIP_INVITE, nil
	case "join":
		return m.MEMBERSHIP_JOIN, nil
	case "leave":
		return m.MEMBERSHIP_LEAVE, nil
	case "ban":
		return m.MEMBERSHIP_BAN, nil
	default:
		return "", fmt.Errorf("matrix: unknown membership in database: %v", membership)
	}
}

func (r *Room) GetJoinRule(db DBI) (m.RoomJoinRule, error) {
	var join_rule string
	row := db.QueryRow(`SELECT join_rule
        FROM room_join_rules
        WHERE room_id=?`, r.ID.String())
	if err := row.Scan(&join_rule); err != nil {
		return m.JOIN_PRIVATE, err
	}
	switch join_rule {
	case "public":
		return m.JOIN_PUBLIC, nil
	case "knock":
		return m.JOIN_KNOCK, nil
	case "invite":
		return m.JOIN_INVITE, nil
	case "private":
		return m.JOIN_PRIVATE, nil
	default:
		return m.JOIN_PRIVATE, fmt.Errorf("matrix: unknown join rule in database: %v", join_rule)
	}
}

func (r *Room) UpdateJoinRule(tx *sql.Tx, sender m.UserID, new_join_rule m.RoomJoinRule) (err error) {
	var event_id m.EventID
	content := map[string]interface{}{"join_rule": string(new_join_rule)}
	if event_id, err = NewStateEvent(tx, sender, r.ID, "m.room.join_rules", "", content); err != nil {
		return
	} else {
		_, err = tx.Exec(`INSERT OR REPLACE
            INTO room_join_rules (
                event_id,
                room_id,
                join_rule)
            VALUES (?, ?, ?)`, event_id.String(), r.ID.String(), string(new_join_rule))
		return
	}
}

func (r *Room) UpdateName(tx *sql.Tx, sender m.UserID, new_name string) (err error) {
	var event_id m.EventID
	content := map[string]interface{}{"name": new_name}
	if event_id, err = NewStateEvent(tx, sender, r.ID, "m.room.name", "", content); err != nil {
		return
	} else {
		_, err = tx.Exec(`INSERT OR REPLACE
            INTO room_names (
                event_id,
                room_id,
                name)
            VALUES (?, ?, ?)`, event_id.String(), r.ID.String(), new_name)
		return
	}
}

func (r *Room) UpdateTopic(tx *sql.Tx, sender m.UserID, new_topic string) (err error) {
	var event_id m.EventID
	content := map[string]interface{}{"topic": new_topic}
	if event_id, err = NewStateEvent(tx, sender, r.ID, "m.room.topic", "", content); err != nil {
		return
	} else {
		_, err = tx.Exec(`INSERT OR REPLACE
            INTO room_topics (
                event_id,
                room_id,
                topic)
            VALUES (?, ?, ?)`, event_id.String(), r.ID.String(), new_topic)
		return
	}
}

func (r *Room) UpdateMember(tx *sql.Tx, sender m.UserID, target m.UserID, new_membership m.RoomMembership) (err error) {
	var event_id m.EventID
	content := map[string]interface{}{"membership": string(new_membership)}
	if _, err = NewStateEvent(tx, sender, r.ID, "m.room.member", target.String(), content); err != nil {
		return
	} else {
		_, err = tx.Exec(`INSERT OR REPLACE
            INTO room_memberships (
                event_id,
                room_id,
                user_id,
                membership)
            VALUES (?, ?, ?, ?)`, event_id.String(), r.ID.String(), target.String(), string(new_membership))
		return
	}
}

func (r *Room) UpdateAliases(tx *sql.Tx, sender m.UserID, new_aliases []m.RoomAlias) (err error) {
	var event_id m.EventID
	aliases := make([]string, len(new_aliases))
	for i := 0; i < len(new_aliases); i++ {
		aliases[i] = new_aliases[i].String()
	}
	content := map[string]interface{}{"aliases": aliases}
	if event_id, err = NewStateEvent(tx, sender, r.ID, "m.room.aliases", "", content); err != nil {
		return
	} else {
		for i := range new_aliases {
			alias := new_aliases[i]
			_, err = tx.Exec(`INSERT OR REPLACE
            INTO room_aliases (
                event_id,
                room_id,
                room_alias,
                server)
            VALUES (?, ?, ?, ?)`, event_id.String(), r.ID.String(), alias.String(), alias.Domain())
		}
		return
	}
}

func (r *Room) UpdatePowerLevels(tx *sql.Tx, sender m.UserID, new_levels map[m.UserID]int64, default_level *int64) (err error) {
	var (
		levels   map[string]interface{} = make(map[string]interface{})
		event_id m.EventID
	)
	for k, v := range new_levels {
		levels[k.String()] = v
	}
	if default_level != nil {
		levels["default"] = default_level
	}
	if event_id, err = NewStateEvent(tx, sender, r.ID, "m.room.power_levels", "", levels); err != nil {
		return
	} else {
		if _, err = tx.Exec(`DELETE
            FROM room_power_levels
            WHERE room_id=?`, r.ID.String()); err != nil {
			return
		}
		for user_id := range new_levels {
			if _, err = tx.Exec(`INSERT OR REPLACE
            INTO room_power_levels (
                event_id,
                room_id,
                user_id,
                level)
            VALUES (?, ?, ?, ?)`, event_id.String(), r.ID.String(), user_id.String(), new_levels[user_id]); err != nil {
				return
			}
		}
		if default_level != nil {
			_, err = tx.Exec(`INSERT OR REPLACE
                INTO room_default_levels (
                    event_id,
                    room_id,
                    level)
                VALUES (?, ?, ?)`, event_id.String(), r.ID.String(), default_level)
		}
	}
	return
}

func (r *Room) UpdateAddStateLevel(tx *sql.Tx, sender m.UserID, new_level int64) (err error) {
	var event_id m.EventID
	content := map[string]interface{}{"level": new_level}
	if event_id, err = NewStateEvent(tx, sender, r.ID, "m.room.add_state_level", "", content); err != nil {
		_, err = tx.Exec(`INSERT OR REPLACE
            INTO room_add_state_levels (
                event_id,
                room_id,
                level)
            VALUES (?, ?, ?)`, event_id.String(), r.ID.String(), new_level)
	}
	return
}

func (r *Room) UpdateSendEventLevel(tx *sql.Tx, sender m.UserID, new_level int64) (err error) {
	var event_id m.EventID
	content := map[string]interface{}{"level": new_level}
	if event_id, err = NewStateEvent(tx, sender, r.ID, "m.room.send_event_level", "", content); err != nil {
		_, err = tx.Exec(`INSERT OR REPLACE
            INTO room_send_event_levels (
                event_id,
                room_id,
                level)
            VALUES (?, ?, ?)`, event_id.String(), r.ID.String(), new_level)
	}
	return
}

func (r *Room) GetOpsLevels(db DBI) (ban int64, kick int64, redact int64, err error) {
	row := db.QueryRow(`SELECT ban_level, kick_level, redact_level
        FROM room_ops_levels
        WHERE room_id=?`, r.ID.String())
	err = row.Scan(&ban, &kick, &redact)
	return
}

func (r *Room) UpdateOpsLevel(tx *sql.Tx, sender m.UserID, new_levels map[string]int64) (err error) {
	var (
		event_id     m.EventID
		content      map[string]interface{} = make(map[string]interface{})
		kick_level   bool                   = false
		ban_level    bool                   = false
		redact_level bool                   = false
	)
	for ltype, level := range new_levels {
		switch ltype {
		case "ban":
			content["ban_level"] = level
			ban_level = true
		case "kick":
			content["kick_level"] = level
			kick_level = true
		case "redact":
			content["redact_level"] = level
			redact_level = true
		default:
			err = fmt.Errorf("matrix: unallowed ops-level: " + ltype)
			return
		}
	}
	if !(kick_level && ban_level && redact_level) {
		err = fmt.Errorf("matrix: not all ops levels specified in request")
		return
	}
	if event_id, err = NewStateEvent(tx, sender, r.ID, "m.room.ops_levels", "", content); err != nil {
		_, err = tx.Exec(`INSERT OR REPLACE
            INTO room_ops_levels (
                event_id,
                room_id,
                kick_level,
                ban_level,
                redact_level)
            VALUES (?, ?, ?, ?, ?)`, event_id.String(), r.ID.String(), content["kick_level"], content["ban_level"], content["redact_level"])
	}
	return
}

func (r *Room) GetInviteLevel(db DBI) (invite int64, err error) {
	row := db.QueryRow(`SELECT level
        FROM room_default_levels
        WHERE room_id=?`, r.ID.String())
	err = row.Scan(&invite)
	return
}

func canJoin(is_self, r_public, u_invited bool) bool {
	return is_self && (r_public || u_invited)
}

func canInvite(is_self, u_joined, s_joined, u_banned bool, s_level, r_inv_level, r_ban_level int64) bool {
	if is_self {
		return false // user can not invite themselves
	} else if s_joined && s_level >= r_inv_level {
		if u_banned && s_level < r_ban_level {
			return false // user is banned and sender is not allowed to remove that
		} else {
			return !u_joined // true if the user isn't already in the room
		}
	} else {
		return false
	}
}

func canLeave(is_self, u_joined, u_invited, s_joined bool, s_level, r_kick_level int64) bool {
	if u_joined || u_invited {
		return is_self || (s_joined && s_level >= r_kick_level)
	} else {
		return false // The user isn't even in the room
	}
}

func canBan(is_self, s_joined bool, s_level, r_ban_level int64) bool {
	return s_joined && s_level >= r_ban_level
}

func (r *Room) CheckedUpdateMember(tx *sql.Tx, sender, target m.UserID, new_membership m.RoomMembership) (err error) {
	var (
		room_join_rule    m.RoomJoinRule
		room_invite_level int64
		sender_membership m.RoomMembership
		sender_level      int64
		user_membership   m.RoomMembership
		room_ban_level    int64
		room_kick_level   int64
	)
	is_self := sender.Compare(target.DomainSpecificString)
	if user_membership, err = r.GetUserMembership(tx, target); err != nil {
		user_membership = m.MEMBERSHIP_NONE
	}
	if is_self && user_membership == new_membership {
		return nil
	} else {
		if room_join_rule, err = r.GetJoinRule(tx); err != nil {
			return fmt.Errorf("matrix: no join rules defined: " + err.Error())
		}
		if room_ban_level, room_kick_level, _, err = r.GetOpsLevels(tx); err != nil {
			room_ban_level = c.Config.DefaultPowerLevel
			room_kick_level = c.Config.DefaultPowerLevel
		}
		if room_invite_level, err = r.GetInviteLevel(tx); err != nil {
			room_invite_level = c.Config.DefaultPowerLevel
		}
		if sender_membership, err = r.GetUserMembership(tx, sender); err != nil {
			sender_membership = m.MEMBERSHIP_NONE
		}
		if sender_level, err = r.GetUserPowerLevel(tx, sender); err != nil {
			sender_level = -1
		}
		var (
			u_banned  bool = user_membership == m.MEMBERSHIP_BAN
			u_invited bool = user_membership == m.MEMBERSHIP_INVITE
			u_joined  bool = user_membership == m.MEMBERSHIP_JOIN
			// u_leaved  bool = user_membership == m.MEMBERSHIP_LEAVE
			s_banned bool = sender_membership == m.MEMBERSHIP_BAN
			// s_invited bool = sender_membership == m.MEMBERSHIP_INVITE
			s_joined bool = sender_membership == m.MEMBERSHIP_JOIN
			// s_leaved  bool = sender_membership == m.MEMBERSHIP_LEAVE
			r_public bool = room_join_rule == m.JOIN_PUBLIC
			// r_knock   bool = room_join_rule == m.JOIN_KNOCK
			// r_invite  bool = room_join_rule == m.JOIN_INVITE
			// r_private bool = room_join_rule == m.JOIN_PRIVATE
		)
		if s_banned {
			return fmt.Errorf("matrix: user %v is banned from this room", sender.String())
		}
		if !is_self && sender_membership != m.MEMBERSHIP_JOIN {
			return fmt.Errorf("matrix: not joined to room")
		}
		switch new_membership {
		case m.MEMBERSHIP_JOIN:
			if !canJoin(is_self, r_public, u_invited) {
				return fmt.Errorf("matrix: can not join room")
			}
		case m.MEMBERSHIP_INVITE:
			if !canInvite(is_self, u_joined, s_joined, u_banned, sender_level, room_invite_level, room_ban_level) {
				return fmt.Errorf("matrix: can not invite to room")
			}
		case m.MEMBERSHIP_LEAVE:
			if !canLeave(is_self, u_joined, u_invited, s_joined, sender_level, room_kick_level) {
				return fmt.Errorf("matrix: can not leave room")
			}
		case m.MEMBERSHIP_BAN:
			if !canBan(is_self, s_joined, sender_level, room_ban_level) {
				return fmt.Errorf("matrix: can not ban user from room")
			}
		}
		err = r.UpdateMember(tx, sender, target, new_membership)
	}
	return
}
