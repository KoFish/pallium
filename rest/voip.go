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
	u "github.com/KoFish/pallium/rest/utils"
	s "github.com/KoFish/pallium/storage"
	"github.com/gorilla/mux"
	"net/http"
)

func setupVoip(root *mux.Router) {
	root.HandleFunc("/voip/turnServer", u.OptionsReply).Methods("OPTIONS")
	root.Handle("/voip/turnServer", u.JSONReply(u.RequireAuth(getTurnServer))).Methods("GET")
}

type turnServer struct{}

func getTurnServer(user *s.User, r *http.Request) (interface{}, error) {
	return turnServer{}, nil
}
