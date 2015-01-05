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

type loginFlow struct {
	Type   string   `json:"type"`
	Stages []string `json:"stages,omitempty"`
}

type loginFlows struct {
	Flows []loginFlow `json:"flows"`
}

func GetLoginFlows(request io.Reader, v Vars, q Query) (interface{}, error) {
	return loginFlows{
		Flows: []loginFlow{
			loginFlow{
				Type:   "m.login.password",
				Stages: []string{},
			},
		},
	}, nil
}

// Request login authentication for a user, an api.Error will be returned
// if the authentication did not succeed. Otherwise a proper response
// according to chosen login type should be returned.
func LoginRequest(request io.Reader, v Vars, q Query) (interface{}, error) {
	var requestdata map[string]interface{} = make(map[string]interface{})
	if err := json.NewDecoder(request).Decode(&requestdata); err != nil {
		return nil, ENotJSON(err.Error())
	}
	reqtype, ok := requestdata["type"]
	if typ, tok := reqtype.(string); ok && tok {
		switch typ {
		case "m.login.password":
			requser := requestdata["user"]
			reqpass := requestdata["password"]
			user, uok := requser.(string)
			pass, pok := reqpass.(string)
			if uok && pok {
				return loginUserByPassword(user, pass)
			} else if !uok {
				return nil, EBadJSON("Invalid or missing user ID")
			} else {
				return nil, EBadJSON("Invalid or missing password")
			}
		default:
			return nil, EBadJSON("Requested login type not supported")
		}
	} else {
		return nil, EBadJSON("Invalid request")
	}
}

// Request registration of a new user, an api.Error will be returned
// if the registration did not succeed. Otherwise a proper response
// according to chosen registration type should be returned.
func RegistrationRequest(request io.Reader, v Vars, q Query) (interface{}, error) {
	var requestdata map[string]interface{} = make(map[string]interface{})
	if err := json.NewDecoder(request).Decode(&requestdata); err != nil {
		return nil, ENotJSON(err.Error())
	}
	reqtype, ok := requestdata["type"]
	if typ, tok := reqtype.(string); ok && tok {
		switch typ {
		case "m.login.password":
			requser := requestdata["user"]
			reqpass := requestdata["password"]
			user, uok := requser.(string)
			pass, pok := reqpass.(string)
			if uok && pok {
				return registerUserByPassword(user, pass)
			} else if !uok {
				return nil, EBadJSON("Invalid or missing user ID")
			} else {
				return nil, EBadJSON("Invalid or missing password")
			}
		default:
			return nil, EBadJSON("Requested registration type not supported")
		}
	} else {
		return nil, EBadJSON("Invalid request")
	}
}

func registerUserByPassword(user, password string) (interface{}, error) {
	var response struct {
		UserID string `json:"user_id"`
		Token  string `json:"access_token"`
	}
	db := s.GetDatabase()
	if tx, err := db.Begin(); err != nil {
		log.Println("matrix-database:", err.Error())
		panic("matrix: could not begin database transaction")
	} else {
		user_id, err := m.ParseUserID(user)
		response.UserID = user_id.String()
		if err != nil || !user_id.IsMine() {
			tx.Rollback()
			if err != nil {
				log.Println(err.Error())
			}
			return nil, EBadJSON("Invalid User ID")
		}
		user, err := s.CreateUser(tx, user_id)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return nil, EUserInUse("Username already taken")
		}
		if err = user.SetPassword(tx, password); err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return nil, EForbidden("Could not set password")
		}
		token, err := user.GetAccessToken(tx)
		response.Token = string(token.Token)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return nil, EForbidden("Could not fetch access token for user")
		}
		return response, tx.Commit()
	}
}

func loginUserByPassword(user, password string) (interface{}, error) {
	var response struct {
		UserID string `json:"user_id"`
		Token  string `json:"access_token"`
	}
	db := s.GetDatabase()
	user_id, err := m.ParseUserID(user)
	response.UserID = user_id.String()
	if err != nil || !user_id.IsMine() {
		if err != nil {
			log.Println(err.Error())
		}
		return nil, EBadJSON("Invalid User ID")
	}
	if user, err := s.GetUser(db, user_id); err != nil {
		log.Println(err.Error())
		return nil, EForbidden("Invalid user and password combination")
	} else {
		if user.Password.Equal(password) {
			token, err := user.GetAccessToken(db)
			response.Token = string(token.Token)
			if err != nil {
				log.Println(err.Error())
				return nil, ENotFound("Could not fetch access token")
			}
			token.UpdateAccessToken(db)
		} else {
			return nil, EForbidden("Invalid user and password combination")
		}
	}
	return response, nil
}
