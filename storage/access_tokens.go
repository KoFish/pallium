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

import ()

type AccessToken struct {
	ID       DBID
	User     *User
	Token    Token
	LastUsed Timestamp
	Created  Timestamp
}

func NewAccessToken(db DBI, u *User) (*AccessToken, error) {
	now := Now()
	token, err := makeToken(u.UserID.String())
	if err != nil {
		return nil, err
	}
	result, err := db.Exec(`
		INSERT OR FAIL
		INTO access_tokens (
			user_id,
			token,
			last_used
			created
		)
		VALUES (?, ?, ?, ?)`, u.ID, token, now, now)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &AccessToken{
		ID:       DBID(id),
		User:     u,
		Token:    Token(token),
		LastUsed: Timestamp(now),
		Created:  Timestamp(now),
	}, nil
}

func (u *User) GetAccessToken(db DBI) (*AccessToken, error) {
	row := db.QueryRow(`
		SELECT id, token, last_used, created
		FROM access_tokens
		WHERE user_id=?`, u.ID)
	var (
		id        int64
		token     string
		last_used int64
		created   int64
	)
	if err := row.Scan(&id, &token, &last_used, &created); err != nil {
		return NewAccessToken(db, u)
	}
	return &AccessToken{
		ID:       DBID(id),
		User:     u,
		Token:    Token(token),
		LastUsed: Timestamp(last_used),
		Created:  Timestamp(created),
	}, nil
}

func (t *AccessToken) UpdateAccessToken(db DBI) error {
	now := Now()
	t.LastUsed = Timestamp(now)
	result, err := db.Exec(`
		UPDATE OR FAIL access_tokens
		SET last_used=?
		WHERE id=?`, now, t.ID)
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
