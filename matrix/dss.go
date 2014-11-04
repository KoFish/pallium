package matrix

import (
    "errors"
    "fmt"
    c "github.com/KoFish/pallium/config"
    "strings"
)

var (
    // event_id_counter is used to ensure reasonably unique event IDs
    event_id_counter int64 = 0
)

type DomainSpecificString struct {
    sigil     string
    localpart string
    domain    string
    is_mine   bool
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
