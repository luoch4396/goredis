package bytepool

// RingBuffer 环形缓冲区
type RingBuffer struct {
	bs      [][]byte
	buf     []byte
	size    int
	read    int
	write   int
	isEmpty bool
}
