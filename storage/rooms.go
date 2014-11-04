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

package storage

import (
    "database/sql"
    "fmt"
    m "github.com/KoFish/pallium/matrix"
)

// These are adapted from https://github.com/matrix-org/synapse/blob/master/synapse/storage/schema/im.sql
const rooms_table = ""

type Room struct {
    ID         m.RoomID
    visibility *m.RoomVisibility
    alias      *[]m.RoomAlias
    name       *string
    topic      *string
}

func CreateRoom(tx *sql.Tx, creator m.UserID) (*Room, error) {
    context, err := m.GenerateRoomID()
    if err != nil {
        return nil, err
    }
    content := map[string]interface{}{"creator": creator.String()}
    if _, err := NewEvent(tx, creator, context, "m.room.create", content); err != nil {
        return nil, err
    }
    return &Room{context, nil, nil, nil, nil}, nil
}

func (r *Room) UpdateJoinRule(tx *sql.Tx, sender m.UserID, new_visibility m.RoomVisibility) (err error) {
    content := map[string]interface{}{"join_rule": string(new_visibility)}
    if _, err = NewStateEvent(tx, sender, r.ID, "m.room.join_rules", "", content); err != nil {
        return
    } else {
        r.visibility = &new_visibility
        return
    }
}

func (r *Room) UpdateName(tx *sql.Tx, sender m.UserID, new_name string) (err error) {
    content := map[string]interface{}{"name": new_name}
    if _, err = NewStateEvent(tx, sender, r.ID, "m.room.name", "", content); err != nil {
        return
    } else {
        r.name = &new_name
        return
    }
}

func (r *Room) UpdateTopic(tx *sql.Tx, sender m.UserID, new_topic string) (err error) {
    content := map[string]interface{}{"topic": new_topic}
    if _, err = NewStateEvent(tx, sender, r.ID, "m.room.topic", "", content); err != nil {
        return
    } else {
        r.topic = &new_topic
        return
    }
}

func (r *Room) UpdateMember(tx *sql.Tx, sender m.UserID, target m.UserID, new_membership m.RoomMembership) (err error) {
    content := map[string]interface{}{"membership": string(new_membership)}
    _, err = NewStateEvent(tx, sender, r.ID, "m.room.member", target.String(), content)
    return
}

func (r *Room) UpdateAliases(tx *sql.Tx, sender m.UserID, new_aliases []m.RoomAlias) (err error) {
    aliases := make([]string, len(new_aliases))
    for i := 0; i < len(new_aliases); i++ {
        aliases[i] = new_aliases[i].String()
    }
    content := map[string]interface{}{"aliases": aliases}
    if _, err = NewStateEvent(tx, sender, r.ID, "m.room.new_aliases", "", content); err != nil {
        return
    } else {
        r.alias = &new_aliases
        return
    }
}

func (r *Room) UpdatePowerLevels(tx *sql.Tx, sender m.UserID, new_levels map[m.UserID]int64, default_level *int64) (err error) {
    var levels map[string]interface{} = make(map[string]interface{})
    for k, v := range new_levels {
        levels[k.String()] = v
    }
    if default_level != nil {
        levels["default"] = default_level
    }
    _, err = NewStateEvent(tx, sender, r.ID, "m.room.power_levels", "", levels)
    return
}

func (r *Room) UpdateAddStateLevel(tx *sql.Tx, sender m.UserID, new_level int64) (err error) {
    content := map[string]interface{}{"level": new_level}
    _, err = NewStateEvent(tx, sender, r.ID, "m.room.add_state_level", "", content)
    return
}

func (r *Room) UpdateSendEventLevel(tx *sql.Tx, sender m.UserID, new_level int64) (err error) {
    content := map[string]interface{}{"level": new_level}
    _, err = NewStateEvent(tx, sender, r.ID, "m.room.send_state_level", "", content)
    return
}

func (r *Room) UpdateOpsLevel(tx *sql.Tx, sender m.UserID, new_levels map[string]int64) (err error) {
    var content map[string]interface{} = make(map[string]interface{})
    for ltype, level := range new_levels {
        switch ltype {
        case "ban":
            content["ban_level"] = level
        case "kick":
            content["kick_level"] = level
        case "redact":
            content["redact_level"] = level
        default:
            err = fmt.Errorf("matrix: unallowed ops-level: " + ltype)
            return
        }
    }
    _, err = NewStateEvent(tx, sender, r.ID, "m.room.send_state_level", "", content)
    return
}
