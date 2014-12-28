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

package utils

import (
	"database/sql"
	"encoding/json"
	m "github.com/KoFish/pallium/matrix"
	s "github.com/KoFish/pallium/storage"
	"io"
	"net/http"
)

type JSONResponse interface{}

type HandlerFunc func(http.ResponseWriter, *http.Request)
type JSONResponseFunc func(http.ResponseWriter, *http.Request) (interface{}, error)
type JSONDBResponseFunc func(*sql.DB, http.ResponseWriter, *http.Request) (interface{}, error)

func JSONReplyHandler(handler JSONResponseFunc) HandlerFunc {
	responsefunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		enc := json.NewEncoder(w)
		if data, err := handler(w, r); err != nil {
			switch err.(type) {
			case JSONError:
				w.WriteHeader(400)
				err.(JSONError).WriteError(w)
			default:
				io.WriteString(w, err.Error())
			}
		} else {
			if err := enc.Encode(data); err != nil {
				error := NewError(m.M_BAD_JSON, err.Error())
				w.WriteHeader(400)
				error.WriteError(w)
			}
		}
	}
	return responsefunc
}

func DBAccessHandler(handler JSONDBResponseFunc) JSONResponseFunc {
	responsefunc := func(w http.ResponseWriter, r *http.Request) (ret interface{}, err error) {
		db := s.GetDatabase()
		defer func() {
			if r := recover(); r != nil {
				ret = nil
				err = NewError(m.M_FORBIDDEN, "Could not access database")
			}
		}()
		ret, err = handler(db, w, r)
		return
	}
	return responsefunc
}

type AuthorizedFunc func(*sql.DB, *s.User, http.ResponseWriter, *http.Request) (interface{}, error)

func RequireAuth(handler AuthorizedFunc) JSONResponseFunc {
	responsefunc := func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		token := r.URL.Query().Get("access_token")
		if token == "" {
			return nil, NewError(m.M_FORBIDDEN, "You are not authorized to access this resource")
		}
		db := s.GetDatabase()
		user, err := s.GetUserByToken(db, s.Token(token))
		if err != nil {
			return nil, NewError(m.M_UNKNOWN_TOKEN, "Could not verify token")
		}
		return handler(db, user, w, r)
	}
	return responsefunc
}

func JSONWithDBReply(handler JSONDBResponseFunc) HandlerFunc {
	return JSONReplyHandler(DBAccessHandler(handler))
}

func JSONWithAuthReply(handler AuthorizedFunc) HandlerFunc {
	return JSONReplyHandler(RequireAuth(handler))
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
