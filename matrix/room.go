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

package matrix

type (
    // Is used when creating rooms to indicate if they are public or private.
    RoomVisibility string
    // Enumerates the different kinds of membership a user can have in a room.
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
