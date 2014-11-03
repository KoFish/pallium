package matrix

import (
    "bytes"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "fmt"
    c "github.com/KoFish/pallium/config"
    "strings"
    "time"
)

type PowerLevel int64

var (
    event_id_counter int64 = 0
)

type DomainSpecificString struct {
    sigil     string
    localpart string
    domain    string
    is_mine   bool
}

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

func (dss DomainSpecificString) String() string {
    return fmt.Sprintf("%s%s:%s", dss.sigil, dss.localpart, dss.domain)
}

func (dss DomainSpecificString) Localpart() string {
    return dss.localpart
}

func (dss DomainSpecificString) Domain() string {
    return dss.domain
}

func (dss DomainSpecificString) IsMine() bool {
    return dss.is_mine
}

func makeDSS(sigil, localpart, domain string) (dss DomainSpecificString, err error) {
    dss = DomainSpecificString{sigil, localpart, domain, domain == c.Hostname}
    if len(localpart) == 0 {
        err = errors.New("Local part of new domain specific string is empty")
    } else if strings.IndexAny(localpart, sigil+":") > -1 {
        err = errors.New("Local part of new domain specific string contains the sigil or a colon")
    } else if len(domain) == 0 {
        err = errors.New("Domain part of new RoomID is empty")
    } else {
        err = nil
    }
    return
}

func parseDSS(sigil, s string) (dss DomainSpecificString, err error) {
    dss = DomainSpecificString{}
    err = nil
    s = strings.Trim(s, " \t\n\r")
    if !strings.HasPrefix(s, sigil) {
        s = "@" + s + ":" + c.Hostname
    }
    if i := strings.LastIndex(s, ":"); i >= 0 {
        localpart := s[1:i]
        domain := s[i+1:]
        dss, err = makeDSS(sigil, localpart, domain)
    } else {
        err = errors.New("Domain specific string contains no domain")
    }
    return
}

func validateDSS(sigil, s string) bool {
    coli := strings.LastIndex(s, ":")
    colj := strings.Index(s, ":")
    return strings.HasPrefix(s, sigil) && coli == colj && coli >= 0 && coli < len(s)
}
