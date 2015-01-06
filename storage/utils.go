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

package storage

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

func generateSalt() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	hash := sha512.Sum512(b)
	return hex.EncodeToString(hash[:]), nil
}

func makeHash(pass, salt string) string {
	hash_bytes := sha512.Sum512([]byte(pass + salt))
	return hex.EncodeToString(hash_bytes[:])

}

func makeToken(base string) (string, error) {
	rstr := make([]byte, 18)
	if _, err := rand.Read(rstr); err != nil {
		return "", err
	}
	b64str := base64.URLEncoding.EncodeToString([]byte(base))
	hexstr := string(hex.EncodeToString(rstr[:]))
	// return (base64.urlsafe_b64encode(user_id).replace('=', '.') + '.' +
	//         stringutils.random_string(18))
	return strings.Replace(b64str, "=", ".", -1) + "." + hexstr, nil
}
