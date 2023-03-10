package exception

import (
	"goredis/pkg/log"
	"goredis/pkg/utils"
)

type NullPointerException struct {
	message string
}

func (e *NullPointerException) Error() string {
	return e.message
}

func CheckIsNotNull(name string) (bool, error) {
	if name == "" {
		return false, &NullPointerException{
			message: name,
		}
	}
	return true, nil
}

func NewNullPointerException(object interface{}, message string) {
	if object == nil {
		log.Error(utils.NewStringBuilder(message, "must be not null!"))
	}
}
