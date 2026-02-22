package main

import (
	"os"
	"syscall"
	"time"
)

// byteAvailable checks whether at least one byte is ready to read from stdin
// within the given timeout.
func byteAvailable(timeout time.Duration) bool {
	fd := int(os.Stdin.Fd())
	var fds syscall.FdSet
	fds.Bits[fd/64] |= 1 << (uint(fd) % 64)
	tv := syscall.NsecToTimeval(int64(timeout))
	err := syscall.Select(fd+1, &fds, nil, nil, &tv)
	if err != nil {
		return false
	}
	return fds.Bits[fd/64]&(1<<(uint(fd)%64)) != 0
}
