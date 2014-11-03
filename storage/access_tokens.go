package storage

import (
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "strings"
    "time"
)

const access_token_table = `
CREATE TABLE IF NOT EXISTS
access_tokens(
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL,
  token TEXT NOT NULL,
  last_used INTEGER,
  FOREIGN KEY(user_id) REFERENCES users(id),
  UNIQUE(token)
)`

type AccessToken struct {
    ID       DBID
    User     *User
    Token    Token
    LastUsed Timestamp
}

func makeToken(base string) (string, error) {
    rstr := make([]byte, 18)
    if _, err := rand.Read(rstr); err != nil {
        return "", err
    }
    b64str := base64.URLEncoding.EncodeToString([]byte(base))
    hexstr := string(hex.EncodeToString(rstr[:]))
    // return (base64.urlsafe_b64encode(user_id).replace('=', '.') + '.' +
    //         stringutils.random_string(18))
    return strings.Replace(b64str, "=", ".", -1) + "." + hexstr, nil
}

func NewAccessToken(db DBI, u *User) (*AccessToken, error) {
    var (
        id    int64
        token string
        now   int64
    )
    now = time.Now().UnixNano() / int64(time.Millisecond)
    token, err := makeToken(u.UserID.String())
    if err != nil {
        return nil, err
    }
    result, err := db.Exec("INSERT OR FAIL INTO access_tokens (user_id, token, last_used) VALUES (?, ?, ?)", u.ID, token, now)
    if err != nil {
        return nil, err
    }
    id, err = result.LastInsertId()
    if err != nil {
        return nil, err
    }
    return &AccessToken{DBID(id), u, Token(token), Timestamp(now)}, nil
}

func (t *AccessToken) UpdateAccessToken(db DBI) error {
    now := time.Now().Unix()
    t.LastUsed = Timestamp(now)
    result, err := db.Exec("UPDATE OR FAIL access_tokens SET last_used=? WHERE id=?", now, t.ID)
    if err != nil {
        return err
    }
    count, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if count < 1 {
        panic("No rows affected when updating token.")
    }
    return nil
}

func (u *User) GetAccessToken(db DBI) (*AccessToken, error) {
    row := db.QueryRow("SELECT id, token, last_used FROM access_tokens WHERE user_id=?", u.ID)
    var (
        id        int64
        token     string
        last_used int64
    )
    if err := row.Scan(&id, &token, &last_used); err != nil {
        return NewAccessToken(db, u)
    }
    return &AccessToken{DBID(id), u, Token(token), Timestamp(last_used)}, nil
}
