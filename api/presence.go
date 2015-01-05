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

package api

import (
	"encoding/json"
	m "github.com/KoFish/pallium/matrix"
	s "github.com/KoFish/pallium/storage"
	"io"
	"log"
)

func GetPresence(cur_user *s.User, request io.Reader, vars Vars, query Query) (interface{}, error) {
	var presence struct {
		Presence  string `json:"presence"`
		StatusMsg string `json:"status_msg,omitempty"`
	}
	if requser, ok := vars["user"]; !ok {
		return nil, ENotFound("Unspecified user")
	} else {
		db := s.GetDatabase()
		userid, err := m.ParseUserID(requser)
		if err != nil {
			log.Println(err.Error())
			return nil, ENotFound("Incorrect User ID")
		}

		if user, err := s.GetUser(db, userid); err != nil {
			log.Println(err.Error())
			return nil, ENotFound("Unknown user")
		} else {
			if cur_presence, err := user.GetPresence(db); err != nil {
				log.Println(err.Error())
				return nil, ENotFound(err.Error())
			} else {
				presence.Presence = cur_presence.PresenceString
				presence.StatusMsg = cur_presence.Status
			}
		}
		return presence, nil
	}
}

func UpdatePresence(user *s.User, request io.Reader, vars Vars, query Query) (interface{}, error) {
	var presence struct {
		Presence  string `json:"presence"`
		StatusMsg string `json:"status_msg,omitempty"`
	}
	if requser, ok := vars["user"]; !ok {
		return nil, ENotFound("Unspecified user")
	} else {
		ruserid, err := m.ParseUserID(requser)
		if err != nil {
			log.Println(err.Error())
			return nil, ENotFound("Invalid UserID")
		}
		if !ruserid.Compare(user.UserID.DomainSpecificString) {
			return nil, EForbidden("Can not change other users presence information")
		}
	}
	if err := json.NewDecoder(request).Decode(&presence); err != nil {
		return nil, ENotJSON(err.Error())
	}

	if new_state, ok := s.PresenceStates[presence.Presence]; ok {
		db := s.GetDatabase()
		if err := user.UpdatePresence(db, new_state, presence.StatusMsg); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	} else {
		return nil, EBadJSON("Unknown new presence level")
	}

	return presence, nil
}
