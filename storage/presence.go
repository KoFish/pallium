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
