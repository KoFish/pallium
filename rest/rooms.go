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
	s "github.com/KoFish/pallium/storage"
	"github.com/gorilla/mux"
	"net/http"
)

type room struct {
	Aliases       []string `json:"aliases"`
	Name          string   `json:"name"`
	JoinedMembers int      `json:"num_joined_members"`
	RoomId        string   `json:"room_id"`
	Topic         string   `json:"topic"`
}

func joinRoom(user *s.User, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	return api.JoinRoom(user, r.Body, mux.Vars(r), api.Query(r.URL.Query()))
}

func createRoom(user *s.User, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	return api.CreateRoom(user, r.Body, mux.Vars(r), api.Query(r.URL.Query()))
}

func listPublicRooms(user *s.User, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	return api.ListPublicRooms(user, r.Body, mux.Vars(r), api.Query(r.URL.Query()))
}

// func roomAliasLookup(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {
//     var (
//         req  struct{}
//         resp struct {
//             RoomID  string   `json:"room_id"`
//             Servers []string `json:"servers"`
//         }
//     )
//     vars := mux.Vars(r)
//     q_room, _ := vars["roomalias"]
//     body, err := ioutil.ReadAll(r.Body)
//     if err != nil {
//         return nil, err
//     }
//     if err = json.Unmarshal(body, &req); err != nil {
//         return nil, u.NewError(m.M_NOT_JSON, "Could not parse json: "+err.Error())
//     }
//     if room_alias, err := m.ParseRoomAlias(q_room); err != nil {
//         return nil, u.NewError(m.M_FORBIDDEN, "Not a valid room alias")
//     } else {
//         room_id, servers, err := s.LookupRoomAlias(db, room_alias)
//         if err != nil {
//             return nil, u.NewError(m.M_FORBIDDEN, "Unknown room alias")
//         } else {
//             resp.RoomID = room_id.String()
//             resp.Servers = servers
//             return resp, nil
//         }
//     }
// }

// func roomAliasCreate(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {
//     var (
//         req struct {
//             RoomID string `json:"room_id"`
//         }
//         resp struct {
//             RoomID  string   `json:"room_id"`
//             Servers []string `json:"servers"`
//         }
//         room_id m.RoomID
//     )
//     vars := mux.Vars(r)
//     q_room, _ := vars["roomalias"]
//     body, err := ioutil.ReadAll(r.Body)
//     if err != nil {
//         return nil, err
//     }
//     if err = json.Unmarshal(body, &req); err != nil {
//         return nil, u.NewError(m.M_NOT_JSON, "Could not parse json: "+err.Error())
//     }
//     if room_alias, err := m.ParseRoomAlias(q_room); err != nil {
//         return nil, u.NewError(m.M_FORBIDDEN, "Not a valid room alias")
//     } else {
//         return nil, u.NewError(m.M_FORBIDDEN, "Not allowed to delete alias")
//     }
// }

// func roomAliasDelete(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {
//     var (
//         req  struct{}
//         resp struct {
//             RoomID string `json:"room_id"`
//         }
//         room_id m.RoomID
//     )
//     vars := mux.Vars(r)
//     q_room, _ := vars["roomalias"]
//     body, err := ioutil.ReadAll(r.Body)
//     if err != nil {
//         return nil, err
//     }
//     if err = json.Unmarshal(body, &req); err != nil {
//         return nil, u.NewError(m.M_NOT_JSON, "Could not parse json: "+err.Error())
//     }
//     if room_alias, err := m.ParseRoomAlias(q_room); err != nil {
//         return nil, u.NewError(m.M_FORBIDDEN, "Not a valid room alias")
//     } else {
//         return nil, u.NewError(m.M_FORBIDDEN, "Not allowed to delete alias")
//     }
// }

// func leaveRoom(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {
//     var (
//         req  leaveRoomRequest
//         resp leaveRoomResponse
//     )
//     if !u.CheckTxnId(r) {
//         return nil, u.NewError(m.M_FORBIDDEN, "Request has already been sent")
//     }
//     vars := mux.Vars(r)
//     q_room := r.URL.Query().Get("room")
//     body, err := ioutil.ReadAll(r.Body)
//     if err != nil {
//         return nil, err
//     }
//     if err = json.Unmarshal(body, &req); err != nil {
//         return nil, u.NewError(m.M_NOT_JSON, "Could not parse json: "+err.Error())
//     }
//     resp = leaveRoomResponse{rid}
//     return resp, nil
// }

// func inviteRoom(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {
//     var (
//         req  inviteRoomRequest
//         resp inviteRoomResponse
//     )
//     if !u.CheckTxnId(r) {
//         return nil, u.NewError(m.M_FORBIDDEN, "Request has already been sent")
//     }
//     vars := mux.Vars(r)
//     q_room := r.URL.Query().Get("room")
//     body, err := ioutil.ReadAll(r.Body)
//     if err != nil {
//         return nil, err
//     }
//     if err = json.Unmarshal(body, &req); err != nil {
//         return nil, u.NewError(m.M_NOT_JSON, "Could not parse json: "+err.Error())
//     }
//     resp = inviteRoomResponse{rid}
//     return resp, nil
// }

// func banRoom(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {
//     var (
//         req  banRoomRequest
//         resp banRoomResponse
//     )
//     if !u.CheckTxnId(r) {
//         return nil, u.NewError(m.M_FORBIDDEN, "Request has already been sent")
//     }
//     vars := mux.Vars(r)
//     q_room := r.URL.Query().Get("room")
//     body, err := ioutil.ReadAll(r.Body)
//     if err != nil {
//         return nil, err
//     }
//     if err = json.Unmarshal(body, &req); err != nil {
//         return nil, u.NewError(m.M_NOT_JSON, "Could not parse json: "+err.Error())
//     }
//     resp = banRoomResponse{rid}
//     return resp, nil
// }
//
//func setRoomState(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {
//    var (
//        req  roomStateRequest
//        resp roomStateResponse
//    )
//    if !u.CheckTxnId(r) {
//        return nil, u.NewError(m.M_FORBIDDEN, "Request has already been sent")
//    }
//    vars := mux.Vars(r)
//    q_room := r.URL.Query().Get("room")
//    q_state_type := r.URL.Query().Get("state_type")
//    q_state_key := r.URL.Query().Get("state_key")
//    body, err := ioutil.ReadAll(r.Body)
//    if err != nil {
//        return nil, err
//    }
//    if err = json.Unmarshal(body, &req); err != nil {
//        return nil, u.NewError(m.M_NOT_JSON, "Could not parse json: "+err.Error())
//    }
//    room_id, err := m.ParseRoomID(q_room)
//    if err != nil {
//        return nil, u.NewError(m.M_FORBIDDEN, "Room is not a valid room ID")
//    }
//    room, err := s.GetRoom(room_id)
//    if err != nil {
//        return nil, u.NewError(m.M_FORBIDDEN, "No such room exists")
//    }
//    resp = roomStateResponse{rid}
//    return resp, nil
//}

func setupRooms(root *mux.Router) {
	// root.HandleFunc("/directory/room/{roomalias}", u.JSONWithAuthReply(roomAliasCreate)).Methods("PUT")
	// root.HandleFunc("/directory/room/{roomalias}", u.JSONWithAuthReply(roomAliasLookup)).Methods("GET")
	// root.HandleFunc("/directory/room/{roomalias}", u.JSONWithAuthReply(roomAliasDelete)).Methods("DELETE")
	root.HandleFunc("/createRoom", u.OptionsReply()).Methods("OPTIONS")
	root.Handle("/createRoom", u.JSONReply(u.RequireAuth(createRoom))).Methods("POST")
	root.Handle("/createRoom/{txnId:[0-9]+}", u.JSONReply(u.TxnID(u.RequireAuth(createRoom)))).Methods("PUT")
	root.Handle("/join/{room}", u.JSONReply(u.RequireAuth(joinRoom))).Methods("POST")
	root.Handle("/join/{room}/{txnId:[0-9]+}", u.JSONReply(u.TxnID(u.RequireAuth(joinRoom)))).Methods("PUT")
	root.Handle("/rooms/{room}/join", u.JSONReply(u.RequireAuth(joinRoom))).Methods("POST")
	root.Handle("/rooms/{room}/join/{txnId:[0-9]+}", u.JSONReply(u.TxnID(u.RequireAuth(joinRoom)))).Methods("PUT")
	root.Handle("/publicRooms", u.JSONReply(u.RequireAuth(listPublicRooms))).Methods("GET")
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
