package nio

type Poll interface {
	newPoll(g *Engine, isListener bool, index int) (*poll, error)
}
