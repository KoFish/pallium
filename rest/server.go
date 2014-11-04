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

package rest

import (
    "fmt"
    c "github.com/KoFish/pallium/config"
    "github.com/gorilla/mux"
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

    // http.Handle("/", client_api_v1)
    http.Handle("/", root)
}

func Start() {
    if err := http.ListenAndServe(fmt.Sprintf(":%i", c.Port), nil); err != nil {
        fmt.Printf("matrix: could not start up server\n >> %v\n", err)
    }
}
