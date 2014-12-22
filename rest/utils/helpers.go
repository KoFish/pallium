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

// The rest/utils package supplies some tools that are useful when processing
// REST requests and generating their responses.
package utils

import (
    "github.com/gorilla/mux"
    "net/http"
    "strconv"
)

var (
    prevrequest = map[string]uint64{}
)

// Check if a certain request is a PUT and has a txnId in their request. If it
// does and the access token of the requests hasn't already made a request with
// the same txnId `true` is returned, otherwise `false`.
func CheckTxnId(r *http.Request) bool {
    vars := mux.Vars(r)
    token := r.URL.Query().Get("access_token")
    if stxnId, ok := vars["txnId"]; ok && r.Method == "PUT" && stxnId != "" {
        if txnId, err := strconv.ParseUint(stxnId, 10, 64); err == nil {
            // ip := strings.Split(r.RemoteAddr, ":")[0]
            // if prevrequest[ip] == txnId {
            if prevrequest[token] == txnId {
                return false
            } else {
                // prevrequest[ip] = txnId
                prevrequest[token] = txnId
            }
        }
    }
    return true
}

// By calling this in a request handler the http://matrix.org live test tool
// is able to make requests.
func AllowMatrixOrg(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.WriteHeader(200)
}
