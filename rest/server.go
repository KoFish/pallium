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

// The rest package deals with serving up the REST API for the server. This is
// where all the processing goes and requests are made to the storage package
// for fetching the actual information.
package rest

import (
	"fmt"
	c "github.com/KoFish/pallium/config"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Setup() {
	root := mux.NewRouter()
	root.StrictSlash(false)
	client_api_v1 := root.PathPrefix("/_matrix/client/api/v1").Subrouter()
	federation_v1 := root.PathPrefix("/_matrix/federation/v1").Subrouter()

	setupLogin(client_api_v1)
	setupRegister(client_api_v1)
	setupEvents(client_api_v1)
	setupRooms(client_api_v1)
	setupFederation(federation_v1)
	setupProfile(client_api_v1)
	setupPresence(client_api_v1)
	setupVoip(client_api_v1)

	// http.Handle("/", client_api_v1)
	http.Handle("/", root)
}

func Start() {
	log.Printf("matrix: starting service at 0.0.0.0:%v", c.Config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", c.Config.Port), nil); err != nil {
		log.Printf("matrix: could not start up server \"%v\"", err)
	}
}
