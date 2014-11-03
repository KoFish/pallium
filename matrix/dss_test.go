package matrix

import "testing"
import c "github.com/KoFish/pallium/config"

func TestParseValidUID(t *testing.T) {
    uid, err := ParseUserID("@kofish:matrix.org")
    t.Log("Parse '@kofish:matrix.org':", uid.String())
    if uid.String() != "@kofish:matrix.org" || err != nil {
        t.Fail()
    }

    uid, err = ParseUserID("kofish")
    t.Log("Parse 'kofish':", uid.String())
    if uid.String() != "@kofish:"+c.Hostname || err != nil {
        t.Fail()
    }
}

func TestParseInvalidUID(t *testing.T) {
    uid, err := ParseUserID("@kofish:matrix.org")
    t.Log("Parse '@kofish:matrix.org':", uid.String())
    if err != nil && uid.IsMine() {
        t.Fail()
    }

    uid, err = ParseUserID("kofish:matrix.org")
    t.Log("Parse 'kofish:matrix.org':", uid.String())
    if err == nil {
        t.Fail()
    }

    uid, err = ParseUserID("kofish:")
    t.Log("Parse 'kofish:':", uid.String())
    if err == nil {
        t.Fail()
    }

    uid, err = ParseUserID("k@fish")
    t.Log("Parse 'k@fish':", uid.String())
    if err == nil {
        t.Fail()
    }
}
