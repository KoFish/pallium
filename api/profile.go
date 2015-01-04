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
	"crypto/md5"
	"encoding/json"
	"fmt"
	m "github.com/KoFish/pallium/matrix"
	s "github.com/KoFish/pallium/storage"
	"io"
)

func GetDisplayName(req_user *s.User, request io.Reader, vars map[string]string) (interface{}, error) {
	var response struct {
		DisplayName string `json:"display_name"`
	}
	if requser, ok := vars["user"]; ok {
		userid, err := m.ParseUserID(requser)
		if err != nil {
			return nil, ENotFound("Invalid User ID in request")
		}
		db := s.GetDatabase()
		user, err := s.GetUser(db, userid)
		if err != nil {
			return nil, ENotFound("Unknown user")
		}
		profile, err := user.GetProfile(db)
		if err != nil {
			return nil, ENotFound(err.Error())
		}
		response.DisplayName = profile.DisplayName
		return response, nil
	} else {
		return nil, ENotFound("Unspecified User ID")
	}
}

func UpdateDisplayName(user *s.User, request io.Reader, vars map[string]string) (interface{}, error) {
	var response struct {
		DisplayName string `json:"display_name"`
	}
	if requser, ok := vars["user"]; !ok {
		return nil, ENotFound("Unspecified user")
	} else {
		ruserid, err := m.ParseUserID(requser)
		if err != nil {
			return nil, ENotFound("Invalid User ID")
		}
		if !ruserid.Compare(user.UserID.DomainSpecificString) {
			return nil, EForbidden("Can not change other users profile information")
		}
	}
	if err := json.NewDecoder(request).Decode(&response); err != nil {
		return nil, ENotJSON(err.Error())
	}
	db := s.GetDatabase()
	if profile, err := user.GetProfile(db); err != nil {
		return nil, ENotFound("Could not get user profile")
	} else {
		if err := profile.UpdateDisplayName(db, response.DisplayName); err != nil {
			return nil, err
		}
	}
	return response, nil
}

func GetAvatarURL(req_user *s.User, request io.Reader, vars map[string]string) (interface{}, error) {
	var response struct {
		AvatarURL string `json:"avatar_url"`
	}
	if requser, ok := vars["user"]; ok {
		userid, err := m.ParseUserID(requser)
		if err != nil {
			return nil, ENotFound("Invalid User ID in request")
		}
		db := s.GetDatabase()
		user, err := s.GetUser(db, userid)
		if err != nil {
			return nil, ENotFound("Unknown user")
		}
		profile, err := user.GetProfile(db)
		if err != nil {
			return nil, ENotFound(err.Error())
		}
		if profile.AvatarURL == "" {
			response.AvatarURL = fmt.Sprintf("https://www.gravatar.com/avatar/%x?d=mm", md5.Sum([]byte(req_user.UserID.String())))
		} else {
			response.AvatarURL = profile.AvatarURL
		}
		return response, nil
	} else {
		return nil, ENotFound("Unspecified User ID")
	}
}

func UpdateAvatarURL(user *s.User, request io.Reader, vars map[string]string) (interface{}, error) {
	var response struct {
		AvatarURL string `json:"avatar_url"`
	}
	if requser, ok := vars["user"]; !ok {
		return nil, ENotFound("Unspecified user")
	} else {
		ruserid, err := m.ParseUserID(requser)
		if err != nil {
			return nil, ENotFound("Invalid UserID")
		}
		if !ruserid.Compare(user.UserID.DomainSpecificString) {
			return nil, EForbidden("Can not change other users profile information")
		}
	}
	if err := json.NewDecoder(request).Decode(&response); err != nil {
		return nil, ENotJSON(err.Error())
	}
	db := s.GetDatabase()
	if profile, err := user.GetProfile(db); err != nil {
		return nil, ENotFound("Could not get user profile")
	} else {
		if err := profile.UpdateAvatarURL(db, response.AvatarURL); err != nil {
			return nil, err
		}
	}
	return response, nil
}
