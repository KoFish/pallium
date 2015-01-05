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

// TODO(): This code needs to be moved to some other place in pallium. This is
//         not rest specific code.
package utils

import (
	"errors"
	"strconv"
	"strings"
)

type StreamToken interface {
	String() string
}

type ForwardStreamToken struct {
	TopologicalOrder uint64
}
type PaginationStreamToken struct {
	ForwardStreamToken
	StreamingOrder uint64
}

func ParseStreamToken(s string) (StreamToken, error) {
	switch s[0] {
	case 's':
		n, err := strconv.ParseUint(s[1:], 10, 64)
		return ForwardStreamToken{n}, err
	case 't':
		ns := strings.SplitN(s[1:], "-", 1)
		if len(ns) == 2 {
			n1, err := strconv.ParseUint(ns[0], 10, 64)
			if err != nil {
				return PaginationStreamToken{}, err
			}
			n2, err := strconv.ParseUint(ns[0], 10, 64)
			return PaginationStreamToken{ForwardStreamToken{n1}, n2}, err
		}
	}
	return ForwardStreamToken{}, errors.New("matrix stream: Could not parse stream token: " + s)
}

func (t ForwardStreamToken) String() string {
	return "s" + strconv.FormatUint(t.TopologicalOrder, 10)
}

func (t PaginationStreamToken) String() string {
	return "t" + strconv.FormatUint(t.TopologicalOrder, 10) + "-" + strconv.FormatUint(t.StreamingOrder, 10)
}

func NewForwardStreamToken(o uint64) ForwardStreamToken {
	return ForwardStreamToken{o}
}

func NewPaginationStreamToken(to, so uint64) PaginationStreamToken {
	return PaginationStreamToken{ForwardStreamToken{to}, so}
}
