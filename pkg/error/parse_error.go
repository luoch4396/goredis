package error

import (
	"errors"
	"goredis/pkg/utils"
)

type ParseError struct {
	Msg string
}

func (e *ParseError) Error() string {
	return utils.NewStringBuilder("parse error: ", e.Msg)
}

func NewParseError(error *ParseError) error {
	return errors.New(error.Error())
}
