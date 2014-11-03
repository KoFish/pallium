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
