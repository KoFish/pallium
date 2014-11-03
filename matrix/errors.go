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
