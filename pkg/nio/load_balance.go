package nio

// LoadBalancer 负载均衡器
type LoadBalancer int

const (
	roundRobin LoadBalancer = iota
	hash
	consistentHash
)

type LoadBalance interface {
	register(eventLoop *poller)

	balance() *poller

	forEach(func(int, *poller) bool)

	getSize() int
}

type roundRobinLoadBalancer struct {
	nextIndex  int
	eventLoops []*poller
	size       int
}

func (lb *roundRobinLoadBalancer) register(el *poller) {
	//el.idx = lb.size
	lb.eventLoops = append(lb.eventLoops, el)
	lb.size++
}

func (lb *roundRobinLoadBalancer) balance() (el *poller) {
	el = lb.eventLoops[lb.nextIndex]
	if lb.nextIndex++; lb.nextIndex >= lb.size {
		lb.nextIndex = 0
	}
	return
}

func (lb *roundRobinLoadBalancer) forEach(f func(int, *poller) bool) {
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
