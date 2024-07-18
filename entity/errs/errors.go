package errs

import (
	"fmt"
)

const (
	SuccessCode = 0
	SuccessMsg  = "success"
)
const (
	// ErrorTypeBusiness 业务错误
	ErrorTypeBusiness = 0
	// ErrorTypeSystem 系统错误，通常可以做熔断逻辑
	ErrorTypeSystem = 1
)

type Error struct {
	Type int32
	Code int32
	Msg  string
}

func New(code int32, msg string) *Error {
	return &Error{Type: ErrorTypeBusiness, Code: code, Msg: msg}
}

func NewSystem(code int32, msg string) *Error {
	return &Error{Type: ErrorTypeSystem, Code: code, Msg: msg}
}

func (err *Error) Error() string {
	if err.Code == SuccessCode {
		return SuccessMsg
	}
	return fmt.Sprintf("err code %d, msg: %s", err.Code, err.Msg)
}

func (err *Error) Newf(args ...any) *Error {
	formatted := fmt.Sprintf(err.Msg, args...)
	if err.Type == ErrorTypeBusiness {
		return New(err.Code, formatted)
	}
	return NewSystem(err.Code, formatted)
}

func (err *Error) IsSystemError() bool {
	return err.Type == ErrorTypeSystem
}
