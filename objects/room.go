package objects

type InitialSyncRoomData struct {
	Membership string           `json:"membership"`
	RoomID     string           `json:"room_id"`
	Messages   *PaginationChunk `json:"messages,omitempty"`
	State      []Event          `json:"state"`
}

type PaginationChunk struct {
	Start string  `json:"start"`
	End   string  `json:"end"`
	Chunk []Event `json:"chunk"`
}

type Event struct {
	EventID   string  `json:"event_id"`
	EventType string  `json:"type"`
	Content   Content `json:"content"`
	RoomID    string  `json:"room_id"`
	UserID    string  `json:"user_id"`
}


type (
	Content map[string]interface{}
)
