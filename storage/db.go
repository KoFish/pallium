package storage

import (
    "database/sql"
    _ "github.com/mxk/go-sqlite/sqlite3"
)

var (
    gDB *sql.DB
)

type DBI interface {
    Exec(string, ...interface{}) (sql.Result, error)
    Query(string, ...interface{}) (*sql.Rows, error)
    QueryRow(string, ...interface{}) *sql.Row
}

func GetDatabase() (db *sql.DB, err error) {
    if gDB == nil {
        db, err = sql.Open("sqlite3", "test.db")
        gDB = db
    } else {
        db = gDB
        err = nil
    }
    return
}

func Setup() error {
    var db_tables map[string]string = map[string]string{
        "user_table":         user_table,
        "rooms_table":        rooms_table,
        "events_table":       events_table,
        "profile_table":      profile_table,
        "presence_table":     presence_table,
        "access_token_table": access_token_table,
    }

    db, err := GetDatabase()
    if err != nil {
        return err
    }
    tx, err := db.Begin()
    if err != nil {
        panic("Could not open database transaction")
    }
    for name, table := range db_tables {
        if _, err := tx.Exec(table); err != nil {
            tx.Rollback()
            panic("Could not setup " + name + ": " + err.Error())
        }
    }
    tx.Commit()
    return nil
}
