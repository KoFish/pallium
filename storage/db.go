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

// The storage package manages setting up the database ond knows how to fetch
// and update all relevant data.
package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/mxk/go-sqlite/sqlite3"
)

var (
	gDB *sql.DB
)

// The DBI is something that can do something on a database, normally either a
// `*sql.DB` or a `*sql.Tx`.
type DBI interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

func GetDatabase() *sql.DB {
	if gDB == nil {
		if db, err := sql.Open("sqlite3", "test.db"); err != nil {
			panic("matrix: could not open database: " + err.Error())
		} else {
			gDB = db
		}
	}
	return gDB
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

	db := GetDatabase()
	tx, err := db.Begin()
	if err != nil {
		panic("Could not open database transaction")
	}
	for name, table := range db_tables {
		fmt.Printf("matrix: setting up DB table %v\n", name)
		if _, err := tx.Exec(table); err != nil {
			tx.Rollback()
			panic("Could not setup " + name + ": " + err.Error())
		}
	}
	tx.Commit()
	return nil
}
