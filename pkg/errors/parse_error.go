package errors

import (
	"errors"
	"goredis/pkg/utils"
)

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
