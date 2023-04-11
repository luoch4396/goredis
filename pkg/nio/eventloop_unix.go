package nio

import "sync"

type eventLoop struct {
	//*timer.Timer
	sync.WaitGroup
	pollers []*poll
}

func NewEventLoop() {

}
