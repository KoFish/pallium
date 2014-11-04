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

package utils

import (
    "encoding/json"
    m "github.com/KoFish/pallium/matrix"
    "io"
    "net/http"
)

type WriterError interface {
    Error(io.Writer)
}

type JSONError struct {
    ErrCode  m.ErrorCode `json:"errcode"`
    ErrorMsg string      `json:"error"`
}

func (err JSONError) Error() string {
    errcode := string(err.ErrCode)
    return errcode + ": " + err.ErrorMsg
}

func NewError(errcode m.ErrorCode, message string) JSONError {
    return JSONError{errcode, message}
}

func (err JSONError) WriteError(w http.ResponseWriter) {
    enc := json.NewEncoder(w)
    if encerr := enc.Encode(err); encerr != nil {
        io.WriteString(w, encerr.Error())
    }
}
