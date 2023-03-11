package exception

import "goredis/pkg/log"

type UnsupportedException struct {
	message string
}

func (e *UnsupportedException) Error() string {
	return e.message
}

func NewUnsupportedException() {
	log.Error("unsupported operation, please check")
}

func NewUnsupportedException0(message string) {
	log.Error(message)
}
