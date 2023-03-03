package biz

import "fmt"

type Error struct {
	arg    int
	errMsg string
}

func NewError(errMsg string) Error {
	return Error{
		errMsg: errMsg,
	}
}

//实现Error接口
func (e *Error) Error() string {
	return fmt.Sprintf("%s", e.errMsg)
}
