//go:build darwin || dragonfly || freebsd || netbsd || openbsd

package nio

import (
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
