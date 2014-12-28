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
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	m "github.com/KoFish/pallium/matrix"
	o "github.com/KoFish/pallium/objects"
	"time"
)

var _ = fmt.Println

const user_table = `
CREATE TABLE IF NOT EXISTS
users(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id TEXT,
  password TEXT,
  salt TEXT,
  creation_ts INTEGER,
  UNIQUE(user_id)
)`

type (
	DBID      int64
	Token     string
	Password  string
	Timestamp int64
)

type PasswordHash struct {
	hash string
	salt string
}

type User struct {
	ID       DBID
	UserID   m.UserID
	Password PasswordHash
	Created  Timestamp
	Profile  *Profile
}

func GetUserByToken(db DBI, token Token) (*User, error) {
	row := db.QueryRow("SELECT users.id, users.user_id, users.password, users.salt, users.creation_ts FROM users, access_tokens WHERE access_tokens.token=? AND users.id=access_tokens.user_id", string(token))
	var (
		id          int64
		user_id     string
		password    string
		salt        string
		creation_ts int64
	)
	if err := row.Scan(&id, &user_id, &password, &salt, &creation_ts); err != nil {
		return nil, err
	}

	uid, err := m.ParseUserID(user_id)
	if err != nil {
		return nil, err
	}
	return &User{DBID(id), uid, PasswordHash{password, salt}, Timestamp(creation_ts), nil}, nil
}

func GetUser(db DBI, uid m.UserID) (*User, error) {
	row := db.QueryRow("SELECT id, password, salt, creation_ts FROM users WHERE user_id=?", uid.String())
	var (
		id          int64
		password    string
		salt        string
		creation_ts int64
	)
	if err := row.Scan(&id, &password, &salt, &creation_ts); err != nil {
		return nil, err
	}
	return &User{DBID(id), uid, PasswordHash{password, salt}, Timestamp(creation_ts), nil}, nil
}

func (u *User) GetRoomMemberships(db DBI) ([]o.InitialSyncRoomData, error) {
	rows, err := db.Query(
		`SELECT r.room_id as id
	FROM room_memberships r
	WHERE r.user_id = ? `, u.UserID.String())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []o.InitialSyncRoomData

	for rows.Next() {
		var roomId string
		err := rows.Scan(&roomId)
		if err != nil {
			return nil, err
		}
		room := o.InitialSyncRoomData{Membership: "joined", RoomID: roomId, State: []o.Event{}}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func CreateUser(db DBI, user_id m.UserID) error {
	result, err := db.Exec("INSERT OR FAIL INTO users (user_id, password, salt, creation_ts) VALUES (?, '', '', ?)", user_id.String(), time.Now().Unix())
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows < 1 {
		panic("Creating users did not affect any row")
	}
	return nil
}

func (u *User) UpdatePassword(db DBI, hash PasswordHash) error {
	result, err := db.Exec("UPDATE OR FAIL users SET password=?, salt=? WHERE id=?", hash.Hash(), hash.Salt(), u.ID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if count < 1 {
		panic("No rows affected when updating user password.")
	}
	if err == nil {
		u.Password = hash
	}
	return err
}

func (u *User) SetPassword(db DBI, passwordstring string) error {
	password := Password(passwordstring)
	salt, err := GenerateSalt()
	if err != nil {
		return err
	}
	hash := password.MakeHash(salt)

	err = u.UpdatePassword(db, hash)

	return err
}

func GenerateSalt() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	hash := sha512.Sum512(b)
	return hex.EncodeToString(hash[:]), nil
}

func (p Password) MakeHash(salt string) PasswordHash {
	salted := string(p) + salt
	hash_bytes := sha512.Sum512([]byte(salted))
	hash := hex.EncodeToString(hash_bytes[:])
	return PasswordHash{hash, salt}
}

func (p PasswordHash) Hash() string {
	return p.hash
}

func (p PasswordHash) Salt() string {
	return p.salt
}

func (p PasswordHash) Equal(password string) bool {
	return bytes.Equal([]byte(Password(password).MakeHash(p.salt).Hash()), []byte(p.hash))
}
