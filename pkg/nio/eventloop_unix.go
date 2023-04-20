package nio

import "sync"

type eventLoop struct {
	//*timer.Timer
	sync.WaitGroup
	pollers []*poll
	*loadBalancer
}

func NewEventLoop() {

}
