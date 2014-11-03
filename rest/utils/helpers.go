package utils

import (
    "github.com/gorilla/mux"
    "net/http"
    "strconv"
    // "strings"
)

var (
    prevrequest = map[string]uint64{}
)

func CheckTxnId(r *http.Request) bool {
    vars := mux.Vars(r)
    token := r.URL.Query().Get("access_token")
    if stxnId := vars["txnId"]; r.Method == "PUT" && stxnId != "" {
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

func AllowMatrixOrg(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "http://matrix.org")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.WriteHeader(200)
}
