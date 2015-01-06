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
	m "github.com/KoFish/pallium/matrix"
	o "github.com/KoFish/pallium/objects"
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

func CreateUser(db DBI, user_id m.UserID) (*User, error) {
	now := Now()
	result, err := db.Exec(`
		INSERT OR FAIL
			INTO users (
				user_id,
				password,
				salt,
				creation_ts)
			VALUES (?, '', '', ?)`, user_id.String(), now)
	if err != nil {
		return nil, err
	}
	if row_id, err := result.LastInsertId(); err != nil {
		panic("matrix: could not get last insert id")
	} else {
		return &User{
			ID:      DBID(row_id),
			UserID:  user_id,
			Created: Timestamp(now),
		}, nil
	}
}

func GetUser(db DBI, uid m.UserID) (*User, error) {
	row := db.QueryRow(`
		SELECT id, password, salt, creation_ts
			FROM users
			WHERE user_id=?`, uid.String())
	var (
		id          int64
		password    string
		salt        string
		creation_ts int64
	)
	if err := row.Scan(&id, &password, &salt, &creation_ts); err != nil {
		return nil, err
	}
	return &User{
		ID:       DBID(id),
		UserID:   uid,
		Password: PasswordHash{hash: password, salt: salt},
		Created:  Timestamp(creation_ts),
		Profile:  nil,
	}, nil
}

func GetUserByToken(db DBI, token Token) (*User, error) {
	row := db.QueryRow(`
		SELECT users.id, users.user_id, users.password, users.salt, users.creation_ts
			FROM users, access_tokens
			WHERE access_tokens.token=?
				AND users.id=access_tokens.user_id`, string(token))
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
	return &User{
		ID:       DBID(id),
		UserID:   uid,
		Password: PasswordHash{hash: password, salt: salt},
		Created:  Timestamp(creation_ts),
		Profile:  nil,
	}, nil
}

// Fetch the users joined rooms and return it in a list for initial sync.
// limit is used to limit the number of messages that is returned for each room.
func (u *User) GetRoomMemberships(db DBI, limit uint64) ([]o.InitialSyncRoomData, error) {
	rows, err := db.Query(`
		SELECT r.room_id as id
			FROM room_memberships AS r
			WHERE r.user_id = ? `, u.UserID.String())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []o.InitialSyncRoomData
	// TODO(): Properly fill out the sync room data

	for rows.Next() {
		var roomId string
		err := rows.Scan(&roomId)
		if err != nil {
			return nil, err
		}
		room := o.InitialSyncRoomData{
			Membership: "joined",
			RoomID:     roomId,
			State:      []o.Event{},
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (u *User) GetInitialPresence(db DBI) ([]o.InitialSyncEvent, error) {
	// TODO(): Get initial presence
	return []o.InitialSyncEvent{}, nil
}

func (u *User) UpdatePassword(db DBI, hash PasswordHash) error {
	result, err := db.Exec(`
		UPDATE OR FAIL users
			SET password=?, salt=?
			WHERE id=?`, hash.Hash(), hash.Salt(), u.ID)
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
	salt, err := generateSalt()
	if err != nil {
		return err
	}
	hash := password.MakeHash(salt)

	err = u.UpdatePassword(db, hash)

	return err
}

func (p Password) MakeHash(salt string) PasswordHash {
	hash := makeHash(string(p), salt)
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
