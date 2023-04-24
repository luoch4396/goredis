package nio

//LoadBalancer 负载均衡器
type LoadBalancer int

const (
	roundRobin LoadBalancer = iota
	hash
	consistentHash
)

type LoadBalance interface {
	register(eventLoop *poll)

	balance() *poll

	forEach(func(int, *poll) bool)

	getSize() int
}

type roundRobinLoadBalancer struct {
	nextIndex  int
	eventLoops []*poll
	size       int
}

func (lb *roundRobinLoadBalancer) register(el *poll) {
	//el.idx = lb.size
	lb.eventLoops = append(lb.eventLoops, el)
	lb.size++
}

func (lb *roundRobinLoadBalancer) balance() (el *poll) {
	el = lb.eventLoops[lb.nextIndex]
	if lb.nextIndex++; lb.nextIndex >= lb.size {
		lb.nextIndex = 0
	}
	return
}

func (lb *roundRobinLoadBalancer) forEach(f func(int, *poll) bool) {
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
