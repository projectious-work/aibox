//go:build linux

package output

import (
	"syscall"
	"unsafe"
)

// IsTerminal reports whether the descriptor is a terminal. It uses the
// TCGETS ioctl, without cgo, so a character device such as /dev/null does
// not count.
func IsTerminal(fd uintptr) bool {
	var termios syscall.Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(&termios)))
	return errno == 0
}
