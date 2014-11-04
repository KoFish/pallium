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
    "encoding/json"
    "fmt"
    c "github.com/KoFish/pallium/config"
    m "github.com/KoFish/pallium/matrix"
    u "github.com/KoFish/pallium/rest/utils"
    s "github.com/KoFish/pallium/storage"
    "github.com/gorilla/mux"
    "io/ioutil"
    "net/http"
)

type createRoomRequest struct {
    Visibility    string   `json:"visibility,omitempty"`
    RoomAliasName string   `json:"room_alias_name,omitempty"`
    Name          string   `json:"name,omitempty"`
    Topic         string   `json:"topic,omitempty"`
    Invite        []string `json:"invite,omitempty"`
}

type createRoomResponse struct {
    RoomID    string `json:"room_id"`
    RoomAlias string `json:"room_alias,omitempty"`
}

func createRoom(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {
    var (
        req        createRoomRequest
        resp       createRoomResponse
        visibility m.RoomVisibility = m.ROOM_PRIVATE
    )
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
    room, err := s.CreateRoom(tx, user.UserID)
    if err != nil {
        tx.Rollback()
        return nil, u.NewError(m.M_FORBIDDEN, "Could not create room event: "+err.Error())
    }
    if req.Visibility != "" {
        switch req.Visibility {
        case "public":
            visibility = m.ROOM_PUBLIC
        case "private":
            visibility = m.ROOM_PRIVATE
        default:
            tx.Rollback()
            return nil, u.NewError(m.M_BAD_JSON, "Unknown room visibility: "+req.Visibility)
        }
    }
    if err := room.UpdateJoinRule(tx, user.UserID, visibility); err != nil {
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
    if len(req.Invite) > 0 {
        for _, invitee := range req.Invite {
            target, err := m.ParseUserID(invitee)
            if err != nil {
                fmt.Printf("matrix: skipping invalid user for invitation on room creation: %v", err.Error())
            } else {
                room.UpdateMember(tx, user.UserID, target, m.MEMBERSHIP_INVITE)
            }
        }
    }
    var default_power_level int64 = c.DefaultPowerLevel
    if err := room.UpdatePowerLevels(tx, user.UserID, map[m.UserID]int64{user.UserID: c.DefaultCreatorPowerLevel}, &default_power_level); err != nil {
        tx.Rollback()
        return nil, u.NewError(m.M_FORBIDDEN, "Could not set creator powerlevel"+err.Error())
    }
    if err := room.UpdateAddStateLevel(tx, user.UserID, c.DefaultPowerLevel); err != nil {
        tx.Rollback()
        return nil, u.NewError(m.M_FORBIDDEN, "Could not set 'add state' power level"+err.Error())
    }
    if err := room.UpdateSendEventLevel(tx, user.UserID, c.DefaultPowerLevel); err != nil {
        tx.Rollback()
        return nil, u.NewError(m.M_FORBIDDEN, "Could not set 'send event' power level"+err.Error())
    }
    if err := room.UpdateOpsLevel(tx, user.UserID, map[string]int64{"ban": c.DefaultPowerLevel, "kick": c.DefaultPowerLevel, "redact": c.DefaultPowerLevel}); err != nil {
        tx.Rollback()
        return nil, u.NewError(m.M_FORBIDDEN, "Could not set ops levels for room"+err.Error())
    }
    rid := room.ID.String()
    resp = createRoomResponse{rid, ""}
    tx.Commit()
    return resp, nil
}

func setupRooms(root *mux.Router) {
    root.HandleFunc("/createRoom", u.JSONWithAuthReply(createRoom)).Methods("POST")
}
