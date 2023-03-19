package error

type NilPointerError struct {
	message string
}

func (e *NilPointerError) Error() string {
	return e.message
}

func CheckIsNotNull(object interface{}) (bool, error) {
	if object == nil {
		return false, &NilPointerError{
			message: "must be not null!",
		}
	}
	return true, nil
}
