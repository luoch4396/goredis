package nio

func setZeroCopy(fd int) error {
	return syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, SoZerocopy, 1)
}
