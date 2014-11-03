package matrix

type (
    RoomVisibility string
    RoomMembership string
)

const (
    ROOM_PUBLIC  RoomVisibility = "public"
    ROOM_PRIVATE RoomVisibility = "private"
)

const (
    MEMBERSHIP_INVITE RoomMembership = "invite"
    MEMBERSHIP_JOIN   RoomMembership = "join"
    MEMBERSHIP_LEAVE  RoomMembership = "leave"
    MEMBERSHIP_BAN    RoomMembership = "ban"
)
