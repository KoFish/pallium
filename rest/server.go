package rest

import (
    "fmt"
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
    if err := http.ListenAndServe(":8008", nil); err != nil {
        fmt.Printf("matrix: could not start up server\n >> %v\n", err)
    }
}
