package errors

import (
	"errors"
	"goredis/pkg/utils"
)

var (
	CRLF            = "\r\n"
	unknownErrBytes = []byte("-Err unknown\r\n")
)

// NilPointerError 空指针
type NilPointerError struct {
	Msg string
}

func (e *NilPointerError) Error() string {
	return e.Msg
}

func NewNilPointerError(msg string) error {
	return errors.New(error.Error(
		&NilPointerError{
			Msg: msg,
		},
	))
}

func CheckIsNotNull(object interface{}) (bool, error) {
	if object == nil {
		return false, &NilPointerError{
			Msg: "nil errors: [object] must be not null!",
		}
	}
	return true, nil
}

// ParseError 解析错误
type ParseError struct {
	Msg string
}

func (e *ParseError) Error() string {
	return utils.NewStringBuilder("parse errors: ", e.Msg)
}

func NewParseError(msg string) error {
	return errors.New(error.Error(
		&ParseError{
			Msg: msg,
		},
	))
}

// StandardError 通用错误
type StandardError struct {
	Status string
}

func (r *StandardError) Info() []byte {
	return []byte("-" + r.Status + CRLF)
}

func NewStandardError(status string) *StandardError {
	return &StandardError{
		//redis 错误需要添加ERR前缀
		Status: "ERR: " + status,
	}
}

func (r *StandardError) Error() string {
	return r.Status
}

// UnknownError 未知错误
type UnknownError struct{}

func (r *UnknownError) Info() []byte {
	return unknownErrBytes
}

func (r *UnknownError) Error() string {
	return "Err unknown"
}
