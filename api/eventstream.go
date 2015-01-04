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

package api

import (
	o "github.com/KoFish/pallium/objects"
	s "github.com/KoFish/pallium/storage"
	"io"
	"strconv"
)

func GetInitialSync(user *s.User, request io.Reader, vars Vars, query Query) (interface{}, error) {
	var sync o.InitialSync
	limit_s, _ := query.GetOne("limit", "16")
	limit, err := strconv.ParseUint(limit_s, 10, 64)
	if err != nil {
		return nil, EBadJSON("Invalid data in limit parameter")
	}
	db := s.GetDatabase()
	// TODO(): Fill sync struct with proper data
	if sync.Rooms, err = user.GetRoomMemberships(db, limit); err != nil {
		return nil, ENotFound("Could not fetch users room memberships")
	}
	if sync.Presence, err = user.GetInitialPresence(db); err != nil {
		return nil, ENotFound("Could not fetch users initial presence information")
	}
	return sync, nil
}

func GetEvents(user *s.User, request io.Reader, vars Vars, query Query) (interface{}, error) {
	from_s, _ := query.GetOne("from", "")
	timeout_s, _ := query.GetOne("timeout", "")
	_ = timeout_s
	events := o.PaginationChunk{
		Start: from_s,
		End:   from_s, // TODO(): Replace temporary value with proper end
	}

	return events, nil
}
