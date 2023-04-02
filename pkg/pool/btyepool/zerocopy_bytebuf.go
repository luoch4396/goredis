package btyepool

import "sync"

// ByteBuf 实现一个no copy字节池
type ByteBuf struct {
	mallocSize int
	freeSize   int
	pool       *sync.Pool
	length     int64
}

type Allocator interface {
	// Malloc 分配
	Malloc(size int) []byte
	// Recycle 回收
	Recycle(buf []byte, size int) []byte
	// Write 写入
	Write(buf []byte, more ...byte) []byte
	// WriteString 写入字符串
	WriteString(buf []byte, more string) []byte
	// Read 读取
	Read()
	// ReadString 读取字符串
	ReadString()
	// Free 释放
	Free(buf []byte)
}

type ByteBufAllocator struct {
}

func (a *ByteBufAllocator) Malloc(size int) []byte {
	return make([]byte, size)
}
