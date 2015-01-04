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

// Handles login requests.
package rest

import (
	"fmt"
	"github.com/KoFish/pallium/api"
	u "github.com/KoFish/pallium/rest/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func getLoginFlows(r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	return api.GetLoginFlows(r.Body)
}

func submitLogin(r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	return api.LoginRequest(r.Body)
}

func submitRegister(r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	return api.RegistrationRequest(r.Body)
}

func fallbackLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not implemented")
}

func setupLogin(root *mux.Router) {
	root.HandleFunc("/login", u.OptionsReply()).Methods("OPTIONS")
	root.Handle("/login", u.JSONReply(getLoginFlows)).Methods("GET")
	root.Handle("/login", u.JSONReply(submitLogin)).Methods("POST")
	root.HandleFunc("/login/fallback", fallbackLogin).Methods("GET")
}

func setupRegister(root *mux.Router) {
	// TODO: Change the following function if the registration and login paths
	//       differ.
	root.HandleFunc("/register", u.OptionsReply()).Methods("OPTIONS")
	root.Handle("/register", u.JSONReply(getLoginFlows)).Methods("GET")
	root.Handle("/register", u.JSONReply(submitRegister)).Methods("POST")
}
