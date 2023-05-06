//go:build darwin || dragonfly || freebsd || netbsd || openbsd

package nio

import (
	"math"
	"syscall"
	"unsafe"
)

func rawSendMsg(fd int, bs [][]byte, ivs []syscall.Iovec, zeroCopy bool) (n int, err error) {
	iovLen := iovecs(bs, ivs)
	if iovLen == 0 {
		return 0, nil
	}
	var msghdr = syscall.Msghdr{
		Iovlen: int32(iovLen),
		Iov:    &ivs[0],
	}
	var zeroCopyFlag uintptr
	if zeroCopy {
		zeroCopyFlag = MSG_ZEROCOPY
	}
	r, _, e := syscall.RawSyscall(syscall.SYS_SENDMSG, uintptr(fd), uintptr(unsafe.Pointer(&msghdr)), zeroCopyFlag)
	resetIovecs(bs, ivs[:iovLen])
	if e != 0 {
		return int(r), e
	}
	return int(r), nil
}

//func sendMsgN(fd int, bs [][]byte, ivs []syscall.Iovec, zeroCopy bool) (n int, err error) {
//	zeroCopyFlag := 0
//	if zeroCopy {
//		zeroCopyFlag = 1
//	}
//	return syscall.SendmsgN(fd, )
//}

func sendMsg() {

}

func iovecs(bs [][]byte, ivs []syscall.Iovec) (iovLen int) {
	totalLen := 0
	for i := 0; i < len(bs); i++ {
		chunk := bs[i]
		l := len(chunk)
		if l == 0 {
			continue
		}
		ivs[iovLen].Base = &chunk[0]
		ivs[iovLen].SetLen(l)
		totalLen += l
		iovLen++
	}
	// iovecs <= 2GB(2^31)
	if totalLen <= math.MaxInt32 {
		return iovLen
	}
	// reset here
	totalLen = math.MaxInt32
	for i := 0; i < iovLen; i++ {
		l := int(ivs[i].Len)
		if l < totalLen {
			totalLen -= l
			continue
		}
		ivs[i].SetLen(totalLen)
		iovLen = i + 1
		resetIovecs(nil, ivs[iovLen:])
		return iovLen
	}
	return iovLen
}

func resetIovecs(bs [][]byte, ivs []syscall.Iovec) {
	for i := 0; i < len(bs); i++ {
		bs[i] = nil
	}
	for i := 0; i < len(ivs); i++ {
		ivs[i].Base = nil
	}
}
