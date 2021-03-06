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
	_ "github.com/mxk/go-sqlite/sqlite3"
	"log"
	"time"
)

type (
	DBID      int64
	Token     string
	Password  string
	Timestamp int64
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

type DB interface {
	DBI
	Begin() (*sql.Tx, error)
}

type TX interface {
	DBI
	Commit() error
	Rollback() error
}

func GetDatabase() DB {
	if gDB == nil {
		if db, err := sql.Open("sqlite3", "test.db"); err != nil {
			panic("matrix: could not open database: " + err.Error())
		} else {
			gDB = db
		}
	}
	return gDB
}

func Now() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func Setup() error {
	db := GetDatabase()
	if tx, err := db.Begin(); err != nil {
		panic("matrix: could not open database transaction")
	} else {
		schemas, err := AssetDir("schemas")
		if err != nil {
			log.Fatal("matrix: no schemas found")
		}
		for _, schema := range schemas {
			log.Printf("matrix: setting up DB table %v\n", schema)
			table, err := Asset("schemas/" + schema)
			if err != nil {
				tx.Rollback()
				panic("matrix: could not load table")
			}
			if _, err := tx.Exec(string(table)); err != nil {
				tx.Rollback()
				log.Println(err.Error())
				panic("matrix: could not setup table")
			}
		}
		tx.Commit()
	}
	return nil
}
