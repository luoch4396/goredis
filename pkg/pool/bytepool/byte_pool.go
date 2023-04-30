package bytepool

import (
	"sync"
)

type Allocator interface {
	Malloc(size int) []byte
	Recycle(buf []byte, size int) []byte
	Append(buf []byte, more ...byte) []byte
	AppendString(buf []byte, more string) []byte
	Free(buf []byte)
}

var DefaultBytePool = NewBytePool(1024, 1024*1024*128)

type BytePool struct {
	bufSize  int
	freeSize int
	pool     *sync.Pool
}

func NewBytePool(bufSize, freeSize int) Allocator {
	if bufSize <= 0 {
		bufSize = 64
	}
	if freeSize <= 0 {
		freeSize = 64 * 1024
	}
	if freeSize < bufSize {
		freeSize = bufSize
	}

	mp := &BytePool{
		bufSize:  bufSize,
		freeSize: freeSize,
		pool:     &sync.Pool{},
	}
	mp.pool.New = func() interface{} {
		buf := make([]byte, bufSize)
		return &buf
	}

	return mp
}

func (mp *BytePool) Malloc(size int) []byte {
	if size > mp.freeSize {
		return make([]byte, size)
	}
	pbuf := mp.pool.Get().(*[]byte)
	n := cap(*pbuf)
	if n < size {
		*pbuf = append((*pbuf)[:n], make([]byte, size-n)...)
	}
	return (*pbuf)[:size]
}

func (mp *BytePool) Recycle(buf []byte, size int) []byte {
	if size <= cap(buf) {
		return buf[:size]
	}

	if cap(buf) < mp.freeSize {
		pbuf := mp.pool.Get().(*[]byte)
		n := cap(buf)
		if n < size {
			*pbuf = append((*pbuf)[:n], make([]byte, size-n)...)
		}
		*pbuf = (*pbuf)[:size]
		copy(*pbuf, buf)
		mp.Free(buf)
		return *pbuf
	}
	return append(buf[:cap(buf)], make([]byte, size-cap(buf))...)[:size]
}

// Append .
func (mp *BytePool) Append(buf []byte, more ...byte) []byte {
	return append(buf, more...)
}

// AppendString .
func (mp *BytePool) AppendString(buf []byte, more string) []byte {
	return append(buf, more...)
}

// Free .
func (mp *BytePool) Free(buf []byte) {
	if cap(buf) > mp.freeSize {
		return
	}
	mp.pool.Put(&buf)
}

// NativeAllocator definition.
type NativeAllocator struct{}

// Malloc .
func (a *NativeAllocator) Malloc(size int) []byte {
	return make([]byte, size)
}

// Recycle .
func (a *NativeAllocator) Recycle(buf []byte, size int) []byte {
	if size <= cap(buf) {
		return buf[:size]
	}
	newBuf := make([]byte, size)
	copy(newBuf, buf)
	return newBuf
}

// Free .
func (a *NativeAllocator) Free(buf []byte) {
}

// Malloc exports default package method.
func Malloc(size int) []byte {
	return DefaultBytePool.Malloc(size)
}

func Recycle(buf []byte, size int) []byte {
	return DefaultBytePool.Recycle(buf, size)
}

// Append exports default package method.
func Append(buf []byte, more ...byte) []byte {
	return DefaultBytePool.Append(buf, more...)
}

// AppendString exports default package method.
func AppendString(buf []byte, more string) []byte {
	return DefaultBytePool.AppendString(buf, more)
}

// Free exports default package method.
func Free(buf []byte) {
	DefaultBytePool.Free(buf)
}
