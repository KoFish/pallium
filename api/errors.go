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

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Error struct {
	ErrorCode string `json:"errcode"`
	ErrorMsg  string `json:"error"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%v: %v", e.ErrorCode, e.ErrorMsg)
}

// WriteTo should really be able to implement the io.WriterTo interface but
// can't since the json.Encode does not return the correct values.
func (e Error) WriteTo(w io.Writer) {
	if err := json.NewEncoder(w).Encode(e); err != nil {
		panic("could not write error to writer")
	}
}

// Create a new JSON-serializable matrix-formatted error message. The intention
// for having this in the api package is that it should be extendable for any
// future transport protocol that is added later on.
func NewError(code, msg string) Error {
	return Error{
		ErrorCode: strings.ToUpper(code),
		ErrorMsg:  msg,
	}
}

func EForbidden(msg string) Error {
	return NewError("M_FORBIDDEN", msg)
}

func EUnknownToken(msg string) Error {
	return NewError("M_UNKNOWN_TOKEN", msg)
}

func EBadJSON(msg string) Error {
	return NewError("M_BAD_JSON", msg)
}

func ENotJSON(msg string) Error {
	return NewError("M_NOT_JSON", msg)
}

func ENotFound(msg string) Error {
	return NewError("M_NOT_FOUND", msg)
}

func EUserInUse(msg string) Error {
	return NewError("M_USER_IN_USE", msg)
}
