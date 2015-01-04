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
	"time"
)

const presence_table = `
CREATE TABLE IF NOT EXISTS
presence (
    user_id INTEGER NOT NULL,
    state INTEGER,
    status_msg TEXT,
    mtime INTEGER,
    FOREIGN KEY(user_id) REFERENCES users(id)
)`

type PresenceState int

type Presence struct {
	User           *User
	Presence       PresenceState
	PresenceString string
	Status         string
	MTime          int64
}

const (
	P_FREEFORCHAT PresenceState = 20
	P_ONLINE      PresenceState = 10
	P_UNAVAILABLE PresenceState = 5
	P_HIDDEN      PresenceState = 1
	P_OFFLINE     PresenceState = 0
)

var (
	PresenceStates map[string]PresenceState = map[string]PresenceState{
		"free for chat": P_FREEFORCHAT,
		"online":        P_ONLINE,
		"unavailable":   P_UNAVAILABLE,
		"hidden":        P_HIDDEN,
		"offline":       P_OFFLINE,
	}
)

func (p PresenceState) String() string {
	for k, v := range PresenceStates {
		if v == p {
			return k
		}
	}
	return ""
}

func (u *User) UpdatePresence(db DBI, newpresence PresenceState, message string) error {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	result, err := db.Exec("INSERT OR FAIL INTO presence VALUES(?,?,?,?)", u.ID, int(newpresence), message, now)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if count < 1 {
		panic("Inserting new presence did not create new rows")
	}
	return err
}

func (u *User) GetPresence(db DBI) (*Presence, error) {
	var (
		state int64
		msg   string
		mtime int64
	)
	row := db.QueryRow(`SELECT state, status_msg mtime
		FROM presence
		WHERE user_id=?`)
	if err := row.Scan(&state, &msg, &mtime); err != nil {
		return nil, err
	}
	curstate := PresenceState(state)
	return &Presence{
		User:           u,
		Presence:       curstate,
		PresenceString: curstate.String(),
		Status:         msg,
		MTime:          mtime,
	}, nil
}
