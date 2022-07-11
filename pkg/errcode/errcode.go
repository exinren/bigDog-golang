package errcode

import (
	"fmt"
	"net/http"
)

type Error struct {
	code	int			`json:"code"`
	msg		string		`json:"msg"`
	details []string	`json:"details"`
}

// 存储新建的状态码
var codes = map[int]string{}


func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
	}
	codes[code] = msg
	return &Error{code: code, msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误码： %d，错误信息：%s", e.code, e.msg)
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Msgf(args []interface{}) string {
	return	fmt.Sprintf(e.msg, args)
}

func (e *Error) Details() []string {
	return e.details
}

func (e *Error) WithDetails(details ...string) *Error {
	NewError := *e
	NewError.details = []string{}
	for _, d := range details {
		NewError.details = append(NewError.details, d)
	}
	return &NewError
}

// 根据错误码，返回状态码
func (e *Error) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK	// 200
	case ServerError.Code():
		return http.StatusInternalServerError	// 500
	case InvalidParams.Code():
		return http.StatusBadGateway	// 502
	case UnauthorizedAuthNotExist.Code():
		fallthrough
	case UnauthorizedTokenError.Code():
		fallthrough
	case UnauthorizedTokenTimeout.Code():
		fallthrough
	case UnauthorizedTokenGenerate.Code():
		return http.StatusUnauthorized	// 401
	case TooManyRequests.Code():
		return http.StatusTooManyRequests	// 429
	}
	return http.StatusInternalServerError // 500
}