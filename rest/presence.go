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

func setupPresence(root *mux.Router) {
	root.HandleFunc("/presence/{user}/status", u.OptionsReply()).Methods("OPTIONS")
	root.Handle("/presence/{user}/status", u.JSONReply(u.RequireAuth(updatePresence))).Methods("PUT")
	root.Handle("/presence/{user}/status", u.JSONReply(u.RequireAuth(getPresence))).Methods("GET")
}

func updatePresence(user *s.User, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	return api.UpdatePresence(user, r.Body, mux.Vars(r))
}

func getPresence(user *s.User, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	return api.GetPresence(user, r.Body, mux.Vars(r))
}
