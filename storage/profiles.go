package storage

import ()

const profile_table = `
CREATE TABLE IF NOT EXISTS
profiles(
    user_id INTEGER NOT NULL,
    display_name TEXT,
    avatar_url TEXT,
    UNIQUE(user_id),
    FOREIGN KEY(user_id) REFERENCES users(id)
)`

type Profile struct {
    User        *User
    DisplayName string
    AvatarURL   string
}

func NewProfile(db DBI, u *User) (*Profile, error) {
    result, err := db.Exec("INSERT OR FAIL INTO profiles SET display_name=?, avatar_url=?, user_id=?", "", "", u.ID)
    if err != nil {
        return nil, err
    }
    count, err := result.RowsAffected()
    if count < 1 {
        panic("No arrows affected on creating new profile")
    }
    if err != nil {
        return nil, err
    }
    profile := &Profile{u, "", ""}
    u.Profile = profile
    return profile, nil
}

func (u *User) GetProfile(db DBI) (*Profile, error) {
    if u.Profile != nil {
        return u.Profile, nil
    }
    row := db.QueryRow("SELECT display_name, avatar_url FROM profiles WHERE user_id=?", u.ID)
    var (
        displayname string
        avatar_url  string
    )
    if err := row.Scan(&displayname, &avatar_url); err != nil {
        return NewProfile(db, u)
    }
    profile := &Profile{u, displayname, avatar_url}
    u.Profile = profile
    return profile, nil
}

func (p *Profile) UpdateDisplayName(db DBI, newname string) error {
    result, err := db.Exec("UPDATE OR FAIL profiles SET display_name=? WHERE user_id=?", newname, p.User.ID)
    if err != nil {
        return err
    }
    count, err := result.RowsAffected()
    if count < 1 {
        panic("No rows affected when updating profile")
    }
    if err == nil {
        p.DisplayName = newname
    }
    return err
}

func (p *Profile) UpdateAvatarURL(db DBI, newurl string) error {
    result, err := db.Exec("UPDATE OR FAIL profiles SET avatar_url=? WHERE user_id=?", newurl, p.User.ID)
    if err != nil {
        return err
    }
    count, err := result.RowsAffected()
    if count < 1 {
        panic("No rows affected when updating profile")
    }
    if err == nil {
        p.AvatarURL = newurl
    }
    return err
}
