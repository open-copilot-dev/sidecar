package common

import (
	"fmt"
)

var ErrCanceled = NewErr(10000, "canceled")

var ErrIgnored = NewErr(10001, "ignored")

type ErrWithCode struct {
	Code int
	Msg  string
}

func NewErr(code int, msg string) *ErrWithCode {
	return &ErrWithCode{
		Code: code,
		Msg:  msg,
	}
}

func (e *ErrWithCode) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Msg)
}
