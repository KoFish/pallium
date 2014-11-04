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

import (
    "bytes"
    "crypto/rand"
    "encoding/base64"
    c "github.com/KoFish/pallium/config"
    "strings"
    "time"
)

type (
    UserID    struct{ DomainSpecificString }
    RoomID    struct{ DomainSpecificString }
    RoomAlias struct{ DomainSpecificString }
    EventID   struct{ DomainSpecificString }
)

func NewUserID(localpart, domain string) (UserID, error) {
    dss, err := makeDSS("@", localpart, domain)
    return UserID{dss}, err
}

func ParseUserID(s string) (UserID, error) {
    dss, err := parseDSS("@", s)
    return UserID{dss}, err
}

func NewRoomID(localpart, domain string) (RoomID, error) {
    dss, err := makeDSS("!", localpart, domain)
    return RoomID{dss}, err
}

func GenerateRoomID() (RoomID, error) {
    rstr := make([]byte, 10)
    if _, err := rand.Read(rstr); err != nil {
        return RoomID{}, err
    }
    localpart := strings.Replace(base64.URLEncoding.EncodeToString(rstr), "=", "", -1)
    return NewRoomID(localpart, c.Hostname)
}

func ParseRoomID(s string) (RoomID, error) {
    dss, err := parseDSS("!", s)
    return RoomID{dss}, err
}

func NewRoomAlias(localpart, domain string) (RoomAlias, error) {
    dss, err := makeDSS("#", localpart, domain)
    return RoomAlias{dss}, err
}

func ParseRoomAlias(s string) (RoomAlias, error) {
    dss, err := parseDSS("#", s)
    return RoomAlias{dss}, err
}

func NewEventID(localpart, domain string) (EventID, error) {
    dss, err := makeDSS("$", localpart, domain)
    return EventID{dss}, err
}

func ParseEventID(s string) (EventID, error) {
    dss, err := parseDSS("$", s)
    return EventID{dss}, err
}

func toBytes(nr int64) []byte {
    var b [8]byte
    for i := 0; i < 8; i++ {
        b[i] = byte((nr >> uint(8*(7-i))) & 0xff)
    }
    return b[:]
}

func GenerateEventID() (ev EventID, err error) {
    idb := bytes.TrimLeft(toBytes(event_id_counter), "\x00")
    event_id_counter += 1
    nowb := bytes.TrimLeft(toBytes(time.Now().Unix()), "\x00")
    rstr := make([]byte, 5)
    if _, err := rand.Read(rstr); err != nil {
        return NewEventID("", "")
    }
    evid := bytes.Join([][]byte{rstr, nowb, idb}, []byte{})
    return NewEventID(strings.Replace(base64.URLEncoding.EncodeToString(evid), "=", "", -1), c.Hostname)
}
