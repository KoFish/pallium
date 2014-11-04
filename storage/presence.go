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

const (
    P_FREEFORCHAT PresenceState = 20
    P_ONLINE      PresenceState = 10
    P_UNAVAILABLE PresenceState = 5
    P_HIDDEN      PresenceState = 1
    P_OFFLINE     PresenceState = 0
)

var (
    PresenceStates map[int]string = map[int]string{
        20: "free for chat",
        10: "online",
        5:  "unavailable",
        1:  "hidden",
        0:  "offline",
    }
)

func (u *User) UpdatePresence(db DBI, newpresence PresenceState, message string) error {
    now := time.Now().UnixNano() / int64(time.Millisecond)
    result, err := db.Exec("INSERT OR FAIL INTO presence SET user_id=?, state=?, status_msg=?, mtime=?", u.ID, int(newpresence), message, now)
    if err != nil {
        return err
    }
    count, err := result.RowsAffected()
    if count < 1 {
        panic("Inserting new presence did not create new rows")
    }
    return err
}
