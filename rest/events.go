package rest

import (
    "database/sql"
    "fmt"
    m "github.com/KoFish/pallium/matrix"
    u "github.com/KoFish/pallium/rest/utils"
    s "github.com/KoFish/pallium/storage"
    "github.com/gorilla/mux"
    "net/http"
    "strconv"
)

var (
    _ = fmt.Println
)

type (
    Content map[string]interface{}
)

type Event struct {
    EventID   string  `json:"event_id"`
    EventType string  `json:"type"`
    Content   Content `json:"content"`
    RoomID    string  `json:"room_id"`
    UserID    string  `json:"user_id"`
}

type StateEvent struct {
    Event
    StateKey      string  `json:"state_key"`
    ReqPowerLevel int64   `json:"required_power_level"`
    PrevContent   Content `json:"prev_content"`
}

type InitialSyncEvent struct {
    Type    string  `json:"type"`
    Content Content `json:"content"`
}

type InitialSyncRoomData struct {
    Membership string           `json:"membership"`
    RoomID     string           `json:"room_id"`
    Messages   *PaginationChunk `json:"messages,omitempty"`
    State      []Event          `json:"state,omitempty"`
}

type PaginationChunk struct {
    Start string  `json:"start"`
    End   string  `json:"end"`
    Chunk []Event `json:"chunk"`
}

type InitialSync struct {
    End      string                `json:"end"`
    Presence []InitialSyncEvent    `json:"presence,omitempty"`
    Rooms    []InitialSyncRoomData `json:"rooms,omitempty"`
}

func getLimit(r *http.Request, def uint64) (limit uint64, err error) {
    slimit := r.URL.Query().Get("limit")
    if slimit != "" {
        limit, err = strconv.ParseUint(slimit, 10, 64)
    } else {
        limit = def
    }
    return
}

func getInitialRoomStates(db s.DBI, user *s.User, limit uint64) ([]InitialSyncRoomData, error) {
    return nil, nil
}

func getInitialEvents(db s.DBI, user *s.User) ([]InitialSyncEvent, error) {
    content := Content{}
    content["user_id"] = "@kofish:kofish.org"
    ev := InitialSyncEvent{"m.presence", content}
    return []InitialSyncEvent{ev}, nil
}

func getInitialSync(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {
    limit, err := getLimit(r, 16)
    if err != nil {
        return nil, u.NewError(m.M_FORBIDDEN, "Limit is not a number")
    }
    var sync InitialSync
    sync.End = "stuff"
    sync.Rooms, _ = getInitialRoomStates(db, user, limit)
    sync.Presence, _ = getInitialEvents(db, user)
    return sync, nil
}

func setupEvents(root *mux.Router) {
    root.HandleFunc("/initialSync", u.AllowMatrixOrg).Methods("OPTIONS")
    root.HandleFunc("/initialSync", u.JSONWithAuthReply(getInitialSync)).Methods("GET")
}
