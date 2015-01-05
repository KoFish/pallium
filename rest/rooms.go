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
	"github.com/KoFish/pallium/api"
	u "github.com/KoFish/pallium/rest/utils"
	// s "github.com/KoFish/pallium/storage"
	"github.com/gorilla/mux"
	// "net/http"
)

func setupRooms(root *mux.Router) {
	// root.Handle("/directory/room/{roomalias}", u.JSONReply(u.RequireAuth(roomAliasCreate))).Methods("PUT")
	// root.Handle("/directory/room/{roomalias}", u.JSONReply(u.RequireAuth(roomAliasLookup))).Methods("GET")
	// root.Handle("/directory/room/{roomalias}", u.JSONReply(u.RequireAuth(roomAliasDelete))).Methods("DELETE")

	root.HandleFunc("/createRoom", u.OptionsReply).Methods("OPTIONS")
	root.Handle("/createRoom", u.AuthAPIEndpoint(api.CreateRoom)).Methods("POST")
	root.Handle("/createRoom/{txnId:[0-9]+}", u.TxnID(u.AuthAPIEndpoint(api.CreateRoom))).Methods("PUT")

	root.Handle("/join/{room}", u.AuthAPIEndpoint(api.JoinRoom)).Methods("POST")
	root.Handle("/join/{room}/{txnId:[0-9]+}", u.TxnID(u.AuthAPIEndpoint(api.JoinRoom))).Methods("PUT")

	root.Handle("/publicRooms", u.AuthAPIEndpoint(api.ListPublicRooms)).Methods("GET")

	root.Handle("/rooms/{room}/join", u.AuthAPIEndpoint(api.JoinRoom)).Methods("POST")
	root.Handle("/rooms/{room}/join/{txnId:[0-9]+}", u.TxnID(u.AuthAPIEndpoint(api.JoinRoom))).Methods("PUT")

	// root.HandleFunc("/rooms/{room}/leave", u.JSONWithAuthReply(leaveRoom)).Methods("POST")
	// root.HandleFunc("/rooms/{room}/leave/{txnId:[0-9]+}", u.JSONWithAuthReply(leaveRoom)).Methods("PUT")
	// root.HandleFunc("/rooms/{room}/invite", u.JSONWithAuthReply(inviteRoom)).Methods("POST")
	// root.HandleFunc("/rooms/{room}/invite/{txnId:[0-9]+}", u.JSONWithAuthReply(inviteRoom)).Methods("PUT")
	// root.HandleFunc("/rooms/{room}/ban", u.JSONWithAuthReply(banRoom)).Methods("POST")
	// root.HandleFunc("/rooms/{room}/ban/{txnId:[0-9]+}", u.JSONWithAuthReply(banRoom)).Methods("PUT")
	// root.HandleFunc("/rooms/{room}/state/{state_type}/{state_key}", u.JSONWithAuthReply(setRoomState)).Methods("POST")
	//	root.HandleFunc("/rooms/{room}/state/{state_type}/{state_key}", u.JSONWithAuthReply(setRoomState)).Methods("PUT")
	// root.HandleFunc("/rooms/{room}/state/{state_type}/{state_key}/{txnId:[0-9]+}", u.JSONWithAuthReply(inviteRoom)).Methods("PUT")
}
