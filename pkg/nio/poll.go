package nio

type poll interface {
	accept() error
}
