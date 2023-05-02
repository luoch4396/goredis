package nio

//no support zero copy
func setZeroCopy(fd int) error {
	return syscall.EINVAL
}
