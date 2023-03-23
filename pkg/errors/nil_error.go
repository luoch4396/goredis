package errors

type NilPointerError struct {
	message string
}

func (e *NilPointerError) Error() string {
	return e.message
}

func CheckIsNotNull(object interface{}) (bool, error) {
	if object == nil {
		return false, &NilPointerError{
			message: "nil errors: must be not null!",
		}
	}
	return true, nil
}
