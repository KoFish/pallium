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
