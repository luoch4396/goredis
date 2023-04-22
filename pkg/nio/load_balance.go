package nio

//LoadBalancer 负载均衡器
type LoadBalancer int

const (
	roundRobin LoadBalancer = iota
	hash
	consistentHash
)

type LoadBalance interface {
	register(eventLoop *eventLoop)

	balance() *eventLoop

	forEach(func(int, *eventLoop) bool)

	getSize() int
}

type roundRobinLoadBalancer struct {
	nextIndex  int
	eventLoops []*eventLoop
	size       int
}

func (lb *roundRobinLoadBalancer) register(el *eventLoop) {
	//el.idx = lb.size
	lb.eventLoops = append(lb.eventLoops, el)
	lb.size++
}

func (lb *roundRobinLoadBalancer) balance() (el *eventLoop) {
	el = lb.eventLoops[lb.nextIndex]
	if lb.nextIndex++; lb.nextIndex >= lb.size {
		lb.nextIndex = 0
	}
	return
}

func (lb *roundRobinLoadBalancer) forEach(f func(int, *eventLoop) bool) {
	for i, v := range lb.eventLoops {
		if !f(i, v) {
			break
		}
	}
}

func (lb *roundRobinLoadBalancer) getSize() int {
	return lb.size
}

//re-balance
