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

package utils

import (
	"encoding/json"
	"github.com/KoFish/pallium/api"
	s "github.com/KoFish/pallium/storage"
	"io"
	"log"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type AuthJSONReply func(*s.User, *http.Request) (interface{}, error)
type JSONReply func(*http.Request) (interface{}, error)

func (fn JSONReply) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if response, err := fn(r); err != nil {
		switch apierr := err.(type) {
		case api.Error:
			log.Printf("[request from %v] %v", r.RemoteAddr, apierr.Error())
			w.WriteHeader(400)
			apierr.WriteTo(w)
		default:
			io.WriteString(w, err.Error())
		}
	} else {
		if err = json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(400)
			api.EBadJSON(err.Error()).WriteTo(w)
		}
	}
}

func RequireAuth(fn AuthJSONReply) JSONReply {
	return func(r *http.Request) (interface{}, error) {
		token := r.URL.Query().Get("access_token")
		if token == "" {
			return nil, api.EForbidden("Not authorized")
		}
		db := s.GetDatabase()
		user, err := s.GetUserByToken(db, s.Token(token))
		if err != nil {
			return nil, api.EUnknownToken("Unknown token")
		}
		return fn(user, r)
	}
}

func OptionsReply() HandlerFunc {
	responsefunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Allow", "POST,GET,PUT,OPTIONS")
	}
	return responsefunc
}
