package error

import "goredis/pkg/utils"

type UnsupportedError struct {
	message string
}

func (e *UnsupportedError) Error() string {
	return utils.NewStringBuilder("unsupported operation error: ", e.message)
}
