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

// The API package is a the transport independent logic to perform the Matrix
// API functionalities from the specification. Any exported function in the
// module is intended to be used by any of the implemented transport protocols
// and only takes and produces json-serializable objects.
package api

import (
	s "github.com/KoFish/pallium/storage"
	"io"
)

/// This file is primarily for documentation of the API package

type (
	Vars  map[string]string
	Query map[string][]string
)

func (q Query) GetOne(key, defval string) (string, bool) {
	vs, ok := q[key]
	if ok {
		return vs[0], true
	} else {
		return defval, false
	}
}

type (
	SimpleEndpoint func(io.Reader, Vars, Query) (interface{}, error)
	AuthEndpoint   func(*s.User, io.Reader, Vars, Query) (interface{}, error)
)
