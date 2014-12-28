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

package rest

import (
	"database/sql"
	"fmt"
	m "github.com/KoFish/pallium/matrix"
	o "github.com/KoFish/pallium/objects"
	u "github.com/KoFish/pallium/rest/utils"
	s "github.com/KoFish/pallium/storage"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	_ = fmt.Println
)

type StateEvent struct {
	Event         o.Event
	StateKey      string    `json:"state_key"`
	ReqPowerLevel int64     `json:"required_power_level"`
	PrevContent   o.Content `json:"prev_content"`
}

type InitialSyncEvent struct {
	Type    string    `json:"type"`
	Content o.Content `json:"content"`
}

type InitialSync struct {
	End      string                  `json:"end"`
	Presence []InitialSyncEvent      `json:"presence,omitempty"`
	Rooms    []o.InitialSyncRoomData `json:"rooms,omitempty"`
}

func getLimit(r *http.Request, def uint64) (limit uint64, err error) {
	slimit := r.URL.Query().Get("limit")
	if slimit != "" {
		limit, err = strconv.ParseUint(slimit, 10, 64)
	} else {
		limit = def
	}
	return
}

func getInitialRoomStates(db s.DBI, user *s.User, limit uint64) ([]o.InitialSyncRoomData, error) {

	value, err := user.GetRoomMemberships(db)

	if err != nil {
		fmt.Println(err)
	}
	return value, nil
}

func getInitialEvents(db s.DBI, user *s.User) ([]InitialSyncEvent, error) {
	content := o.Content{}
	content["user_id"] = user.UserID.String()
	ev := InitialSyncEvent{"m.presence", content}
	return []InitialSyncEvent{ev}, nil
}

func getInitialSync(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	limit, err := getLimit(r, 16)
	if err != nil {
		return nil, u.NewError(m.M_FORBIDDEN, "Limit is not a number")
	}
	var sync InitialSync
	sync.End = "stuff"
	sync.Rooms, _ = getInitialRoomStates(db, user, limit)
	sync.Presence, _ = getInitialEvents(db, user)
	return sync, nil
}

func getEvents(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {

	fromToken := r.URL.Query().Get("from")
	events := o.PaginationChunk{Start: fromToken, End: fmt.Sprintf("%d", rand.Intn(100000)), Chunk: []o.Event{o.Event{}}}
	time.Sleep(5000 * time.Millisecond)

	return events, nil
}

func setupEvents(root *mux.Router) {
	root.HandleFunc("/initialSync", u.AllowMatrixOrg).Methods("OPTIONS")
	root.HandleFunc("/initialSync", u.JSONWithAuthReply(getInitialSync)).Methods("GET")
	root.HandleFunc("/events", u.AllowMatrixOrg).Methods("OPTIONS")
	root.HandleFunc("/events", u.JSONWithAuthReply(getEvents)).Methods("GET")
}
