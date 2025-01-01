package common

import (
	"fmt"
)

// 业务错误码
var (
	ErrCodeCanceled = 10000
	ErrCodeIgnored  = 10001
	ErrCodeIllegal  = 10002
	ErrCodeNotFound = 10004
)

// 系统错误码
var (
	ErrCodeIo      = 20001
	ErrCodeMarshal = 20002
)

var (
	ErrCanceled = NewErr(ErrCodeCanceled, "canceled")
	ErrIgnored  = NewErr(ErrCodeIgnored, "ignored")
	ErrNotFound = NewErr(ErrCodeNotFound, "not found")
	ErrIllegal  = NewErr(ErrCodeIllegal, "illegal")
)

//--------------------------------------------------------------------

type ErrWithCode struct {
	Code  int
	Msg   string
	Cause error
}

func NewErr(code int, msg string) *ErrWithCode {
	return &ErrWithCode{
		Code:  code,
		Msg:   msg,
		Cause: nil,
	}
}

func NewErrWithCause(code int, msg string, cause error) *ErrWithCode {
	return &ErrWithCode{
		Code:  code,
		Msg:   msg,
		Cause: cause,
	}
}

func (e *ErrWithCode) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Msg, e.Cause.Error())
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Msg)
}
