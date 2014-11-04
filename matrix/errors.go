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

type ErrorCode string

const (
    M_FORBIDDEN               ErrorCode = "M_FORBIDDEN"
    M_UNKNOWN_TOKEN           ErrorCode = "M_UNKNOWN_TOKEN"
    M_BAD_JSON                ErrorCode = "M_BAD_JSON"
    M_NOT_JSON                ErrorCode = "M_NOT_JSON"
    M_NOT_FOUND               ErrorCode = "M_NOT_FOUND"
    M_LIMIT_EXCEEDED          ErrorCode = "M_LIMIT_EXCEEDED"
    M_USER_IN_USE             ErrorCode = "M_USER_IN_USE"
    M_ROOM_IN_USE             ErrorCode = "M_ROOM_IN_USE"
    M_BAD_PAGINATION          ErrorCode = "M_BAD_PAGINATION"
    M_LOGIN_EMAIL_URL_NOT_YET ErrorCode = "M_LOGIN_EMAIL_URL_NOT_YET"
)
