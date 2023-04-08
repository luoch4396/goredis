package btyepool

import "sync"

type Allocator interface {
	// Malloc 分配
	Malloc(size int) []byte
	// Recycle 回收
	Recycle(buf []byte, size int) []byte
	// Write 写入
	Write(buf []byte, more ...byte) []byte
	// WriteString 写入字符串
	WriteString(buf []byte, more string) []byte
	// Free 释放
	Free(buf []byte)
}

// ByteBuffer 字节池
type ByteBuffer struct {
	bufferSize int
	freeSize   int
	pool       *sync.Pool
	cap        int64
}

var Default = NewByteBuffer(1024, 1024*32)

func NewByteBuffer(bufferSize, freeSize int) Allocator {
	if bufferSize <= 0 {
		bufferSize = 32
	}
	if freeSize <= 0 {
		freeSize = 1024 * 32
	}
	if freeSize < bufferSize {
		freeSize = bufferSize
	}

	b := &ByteBuffer{
		bufferSize: bufferSize,
		freeSize:   freeSize,
		pool:       &sync.Pool{},
	}
	b.pool.New = func() interface{} {
		buf := make([]byte, bufferSize)
		return &buf
	}

	return b
}

func (b *ByteBuffer) Malloc(size int) []byte {
	if size > b.freeSize {
		return make([]byte, size)
	}
	freeBuf := b.pool.Get().(*[]byte)
	n := cap(*freeBuf)
	if n < size {
		*freeBuf = append((*freeBuf)[:n], make([]byte, size-n)...)
	}
	return (*freeBuf)[:size]
}

func (b *ByteBuffer) Recycle(buf []byte, size int) []byte {
	n := cap(buf)
	if size <= n {
		return buf[:size]
	}
	if n < b.freeSize {
		freeBuf := b.pool.Get().(*[]byte)
		if n < size {
			*freeBuf = append((*freeBuf)[:n], make([]byte, size-n)...)
		}
		*freeBuf = (*freeBuf)[:size]
		copy(*freeBuf, buf)
		b.Free(buf)
		return *freeBuf
	}
	return append(buf[:n], make([]byte, size-n)...)[:size]
}

func (b *ByteBuffer) Write(buf []byte, more ...byte) []byte {
	return append(buf, more...)
}

func (b *ByteBuffer) WriteString(buf []byte, more string) []byte {
	return append(buf, more...)
}

func (b *ByteBuffer) Free(buf []byte) {
	if cap(buf) > b.freeSize {
		return
	}
	b.pool.Put(&buf)
}

type ByteBufAllocator struct {
}

func (a *ByteBufAllocator) Malloc(size int) []byte {
	return make([]byte, size)
}

func (a *ByteBufAllocator) Recycle(buf []byte, size int) []byte {
	if size <= cap(buf) {
		return buf[:size]
	}
	newBuf := make([]byte, size)
	copy(newBuf, buf)
	return newBuf
}

//---------------------- use default byte buffer-----------------------------

func DefaultMalloc(size int) []byte {
	return Default.Malloc(size)
}

func DefaultRecycle(buf []byte, size int) []byte {
	return Default.Recycle(buf, size)
}

func DefaultWrite(buf []byte, more ...byte) []byte {
	return Default.Write(buf, more...)
}

func DefaultWriteString(buf []byte, more string) []byte {
	return Default.WriteString(buf, more)
}

func DefaultFree(buf []byte) {
	Default.Free(buf)
}

func DefaultInit(bufSize, freeSize int) {
	Default = NewByteBuffer(bufSize, freeSize)
}
