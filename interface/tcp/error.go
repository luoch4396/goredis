package tcp

type ErrorInfo interface {
	Error() string
	Info() []byte
}
