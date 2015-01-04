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

func getEvents(user *s.User, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	return api.GetEvents(user, r.Body, mux.Vars(r), api.Query(r.URL.Query()))
}

func getInitialSync(user *s.User, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	return api.GetInitialSync(user, r.Body, mux.Vars(r), api.Query(r.URL.Query()))
}

func setupEvents(root *mux.Router) {
	root.HandleFunc("/events", u.OptionsReply()).Methods("OPTIONS")
	root.HandleFunc("/initialSync", u.OptionsReply()).Methods("OPTIONS")
	root.Handle("/events", u.JSONReply(u.RequireAuth(getEvents))).Methods("GET")
	root.Handle("/initialSync", u.JSONReply(u.RequireAuth(getInitialSync))).Methods("GET")
}
