package error

type UnsupportedError struct {
	message string
}

func (e *UnsupportedError) Error() string {
	return e.message
}
