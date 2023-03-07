package exception

type UnsupportedException struct {
	message string
}

func (e *UnsupportedException) Error() string {
	return e.message
}

func NewUnsupportedException() {
	panic("unsupported operation")
}

func NewUnsupportedException0(message string) {
	panic(message)
}
