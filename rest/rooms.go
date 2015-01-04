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
	"encoding/json"
	c "github.com/KoFish/pallium/config"
	m "github.com/KoFish/pallium/matrix"
	u "github.com/KoFish/pallium/rest/utils"
	s "github.com/KoFish/pallium/storage"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type room struct {
	Aliases       []string `json:"aliases"`
	Name          string   `json:"name"`
	JoinedMembers int      `json:"num_joined_members"`
	RoomId        string   `json:"room_id"`
	Topic         string   `json:"topic"`
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

func createRoom(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var (
		req struct {
			Visibility    string   `json:"visibility,omitempty"`
			RoomAliasName string   `json:"room_alias_name,omitempty"`
			Name          string   `json:"name,omitempty"`
			Topic         string   `json:"topic,omitempty"`
			Invite        []string `json:"invite,omitempty"`
		}
		resp struct {
			RoomID    string `json:"room_id"`
			RoomAlias string `json:"room_alias,omitempty"`
		}
		join_rule m.RoomJoinRule = m.JOIN_INVITE
	)
	if !u.CheckTxnId(r) {
		return nil, u.NewError(m.M_FORBIDDEN, "Request has already been sent")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &req); err != nil {
		return nil, u.NewError(m.M_NOT_JSON, "Could not parse json: "+err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, u.NewError(m.M_FORBIDDEN, "Could not create database transaction: "+err.Error())
	}
	if req.Visibility != "" {
		switch req.Visibility {
		case "public":
			join_rule = m.JOIN_PUBLIC
		case "private":
			join_rule = m.JOIN_INVITE
		default:
			tx.Rollback()
			return nil, u.NewError(m.M_BAD_JSON, "Unknown room visibility: "+req.Visibility)
		}
	}
	is_public := join_rule == m.JOIN_PUBLIC || join_rule == m.JOIN_KNOCK
	room, err := s.CreateRoom(tx, user.UserID, is_public)
	if err != nil {
		tx.Rollback()
		return nil, u.NewError(m.M_FORBIDDEN, "Could not create room event: "+err.Error())
	}
	if err := room.UpdateJoinRule(tx, user.UserID, join_rule); err != nil {
		tx.Rollback()
		return nil, u.NewError(m.M_FORBIDDEN, "Could not update join rule: "+err.Error())
	}
	if req.RoomAliasName != "" {
		roomalias, err := m.ParseRoomAlias(req.RoomAliasName)
		if err != nil {
			tx.Rollback()
			return nil, u.NewError(m.M_BAD_JSON, "Error parsing room alias: "+err.Error())
		}
		if err = room.UpdateAliases(tx, user.UserID, []m.RoomAlias{roomalias}); err != nil {
			tx.Rollback()
			return nil, u.NewError(m.M_FORBIDDEN, "Could not update room aliases: "+err.Error())
		}
	}
	if req.Topic != "" {
		if err := room.UpdateTopic(tx, user.UserID, req.Name); err != nil {
			tx.Rollback()
			return nil, u.NewError(m.M_FORBIDDEN, "Could not update room topic")
		}
	}
	if req.Name != "" {
		if err := room.UpdateName(tx, user.UserID, req.Name); err != nil {
			tx.Rollback()
			return nil, u.NewError(m.M_FORBIDDEN, "Could not update room name")
		}
	}
	if err := room.UpdateMember(tx, user.UserID, user.UserID, m.MEMBERSHIP_JOIN); err != nil {
		tx.Rollback()
		return nil, u.NewError(m.M_FORBIDDEN, "Could not add creator as member of new room: "+err.Error())
	}
	if len(req.Invite) > 0 {
		for _, invitee := range req.Invite {
			target, err := m.ParseUserID(invitee)
			if err != nil {
				log.Printf("matrix: skipping invalid user for invitation on room creation: %v", err.Error())
			} else {
				room.UpdateMember(tx, user.UserID, target, m.MEMBERSHIP_INVITE)
			}
		}
	}
	var default_power_level int64 = c.Config.DefaultPowerLevel
	if err := room.UpdatePowerLevels(tx, user.UserID, map[m.UserID]int64{user.UserID: c.Config.DefaultCreatorPowerLevel}, &default_power_level); err != nil {
		tx.Rollback()
		return nil, u.NewError(m.M_FORBIDDEN, "Could not set creator powerlevel"+err.Error())
	}
	if err := room.UpdateAddStateLevel(tx, user.UserID, c.Config.DefaultPowerLevel); err != nil {
		tx.Rollback()
		return nil, u.NewError(m.M_FORBIDDEN, "Could not set 'add state' power level"+err.Error())
	}
	if err := room.UpdateSendEventLevel(tx, user.UserID, c.Config.DefaultPowerLevel); err != nil {
		tx.Rollback()
		return nil, u.NewError(m.M_FORBIDDEN, "Could not set 'send event' power level"+err.Error())
	}
	if err := room.UpdateOpsLevel(tx, user.UserID, map[string]int64{"ban": c.Config.DefaultPowerLevel, "kick": c.Config.DefaultPowerLevel, "redact": c.Config.DefaultPowerLevel}); err != nil {
		tx.Rollback()
		return nil, u.NewError(m.M_FORBIDDEN, "Could not set ops levels for room"+err.Error())
	}
	resp.RoomID = room.ID.String()
	tx.Commit()
	return resp, nil
}

func joinRoom(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var (
		req  struct{}
		resp struct {
			RoomID string `json:"room_id"`
		}
		room_id m.RoomID
	)
	if !u.CheckTxnId(r) {
		return nil, u.NewError(m.M_FORBIDDEN, "Request has already been sent")
	}
	vars := mux.Vars(r)
	q_room, _ := vars["room"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &req); err != nil {
		return nil, u.NewError(m.M_NOT_JSON, "Could not parse json: "+err.Error())
	}
	if room_alias, err := m.ParseRoomAlias(q_room); err != nil {
		if room_id, err = m.ParseRoomID(q_room); err != nil {
			return nil, u.NewError(m.M_FORBIDDEN, "Could not parse room identifier argument to a room ID")
		}
	} else {
		if room_id, _, err = s.LookupRoomAlias(db, room_alias); err != nil {
			return nil, u.NewError(m.M_FORBIDDEN, "Could not resolve room alias to room id: "+err.Error())
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, u.NewError(m.M_FORBIDDEN, "Could not start database transaction: "+err.Error())
	}
	room := s.Room{ID: room_id}
	if err := room.CheckedUpdateMember(tx, user.UserID, user.UserID, m.MEMBERSHIP_JOIN); err != nil {
		tx.Rollback()
		return nil, u.NewError(m.M_FORBIDDEN, "Could not join room: "+err.Error())
	}
	tx.Commit()
	resp.RoomID = room_id.String()
	return resp, nil
}

func listPublicRooms(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var (
		res struct {
			Chunk []s.Room `json:"chunk"`
			End   string   `json:"end"`
			Start string   `json:"start"`
		}
	)
	res.Start = "START"
	res.End = "END"

	tx, err := db.Begin()
	if err != nil {
		return nil, u.NewError(m.M_FORBIDDEN, "Could not start database transaction: "+err.Error())
	}
	rooms := s.GetPublicRooms(tx)

	res.Chunk = rooms

	tx.Commit()
	return res, nil
}

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
	root.HandleFunc("/createRoom", u.JSONWithAuthReply(createRoom)).Methods("POST")
	root.HandleFunc("/createRoom", u.OptionsReply()).Methods("OPTIONS")
	root.HandleFunc("/createRoom/{txnId:[0-9]+}", u.JSONWithAuthReply(createRoom)).Methods("PUT")
	root.HandleFunc("/join/{room}", u.JSONWithAuthReply(joinRoom)).Methods("POST")
	root.HandleFunc("/join/{room}/{txnId:[0-9]+}", u.JSONWithAuthReply(joinRoom)).Methods("PUT")
	root.HandleFunc("/rooms/{room}/join", u.JSONWithAuthReply(joinRoom)).Methods("POST")
	root.HandleFunc("/rooms/{room}/join/{txnId:[0-9]+}", u.JSONWithAuthReply(joinRoom)).Methods("PUT")
	root.HandleFunc("/publicRooms", u.JSONWithAuthReply(listPublicRooms)).Methods("GET")
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
