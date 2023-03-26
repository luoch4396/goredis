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
	message string
}

func (e *NilPointerError) Error() string {
	return e.message
}

func CheckIsNotNull(object interface{}) (bool, error) {
	if object == nil {
		return false, &NilPointerError{
			message: "nil errors: [object] must be not null!",
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
		Status: status,
	}
}

func (r *StandardError) Error() string {
	return r.Status
}

// UnsupportedError 不支持的操作
type UnsupportedError struct {
	message string
}

func (e *UnsupportedError) Error() string {
	return "Err unknown"
}

func (e *UnsupportedError) Info() []byte {
	return unknownErrBytes
}
