package error

type ParseError struct {
	message string
}

func (e *ParseError) Error() string {
	return e.message
}
