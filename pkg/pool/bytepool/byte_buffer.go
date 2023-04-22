package bytepool

import "sync"

type Allocator interface {
	// Malloc 分配
	Malloc(size int) []byte
	// Recycle 回收
	Recycle(size int) []byte
	// Write 写入
	Write(more ...byte) []byte
	// WriteString 写入字符串
	WriteString(more string) []byte
	// Free 释放
	Free()

	GetCap() int
}

// ByteBuffer 字节池
type ByteBuffer struct {
	bufferSize int
	freeSize   int
	cap        int64
	buf        []byte
}

var (
	pool        *sync.Pool
	defaultSize = 512
)

func init() {
	pool = &sync.Pool{}
	pool.New = func() interface{} {
		buf := make([]byte, defaultSize)
		return &buf
	}
}

func NewByteBuffer(bufferSize, freeSize int) Allocator {
	if bufferSize <= 0 {
		bufferSize = 32
	}
	if freeSize <= 0 {
		freeSize = 1024
	}
	if freeSize < bufferSize {
		freeSize = bufferSize
	}

	b := &ByteBuffer{
		bufferSize: bufferSize,
		freeSize:   freeSize,
		buf:        make([]byte, defaultSize),
	}
	return b
}

func (b *ByteBuffer) GetCap() int {
	return 0
}

func (b *ByteBuffer) Malloc(size int) []byte {
	if size > b.freeSize {
		return make([]byte, size)
	}
	freeBuf := pool.Get().(*[]byte)
	n := cap(*freeBuf)
	if n < size {
		*freeBuf = append((*freeBuf)[:n], make([]byte, size-n)...)
	}
	return (*freeBuf)[:size]
}

func (b *ByteBuffer) Recycle(size int) []byte {
	n := cap(b.buf)
	if size <= n {
		return b.buf[:size]
	}
	if n < b.freeSize {
		freeBuf := pool.Get().(*[]byte)
		if n < size {
			*freeBuf = append((*freeBuf)[:n], make([]byte, size-n)...)
		}
		*freeBuf = (*freeBuf)[:size]
		copy(*freeBuf, b.buf)
		b.Free()
		return *freeBuf
	}
	return append(b.buf[:n], make([]byte, size-n)...)[:size]
}

func (b *ByteBuffer) Write(more ...byte) []byte {
	return append(b.buf, more...)
}

func (b *ByteBuffer) WriteString(more string) []byte {
	return append(b.buf, more...)
}

func (b *ByteBuffer) Free() {
	if cap(b.buf) > b.freeSize {
		return
	}
	pool.Put(&b.buf)
}
