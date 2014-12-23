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
	u "github.com/KoFish/pallium/rest/utils"
	s "github.com/KoFish/pallium/storage"
	"github.com/gorilla/mux"
	"net/http"
	"database/sql"
	"encoding/json"
	"fmt"
)

type presence struct {
	Presence string `json:"presence"`
}


func setupPresence(root *mux.Router) {
	root.HandleFunc("/presence/{user}/status", u.OptionsReply()).Methods("OPTIONS")
	root.HandleFunc("/presence/{user}/status", u.JSONWithAuthReply(updatePresence)).Methods("PUT")
}


func updatePresence(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface {}, error) {
	presenceState := presence{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&presenceState)

	err := user.UpdatePresence(db, s.PresenceStates[presenceState.Presence], "foobar")

	if(err != nil) {
		fmt.Println(err)
	}

	fmt.Println(presenceState.Presence)
	return nil, nil
}
