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
	"github.com/KoFish/pallium/api"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var (
	prevrequest = map[string]uint64{}
)

func TxnID(fn JSONReply) JSONReply {
	return func(r *http.Request) (interface{}, error) {
		txnid, txnok := mux.Vars(r)["txnId"]
		token := r.URL.Query().Get("access_token")
		if r.Method == "PUT" && txnok && txnid != "" {
			if itxnid, err := strconv.ParseUint(txnid, 10, 64); err == nil {
				if prevrequest[token] == itxnid {
					return nil, api.EForbidden("Repeated transaction id")
				}
				return fn(r)
			} else {
				return nil, api.ENotFound("Invalid transaction id, needs to be an integer")
			}
		} else {
			return nil, api.EForbidden("Not allowed without transaction id")
		}
	}
}
