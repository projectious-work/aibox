//go:build darwin

package output

import (
	"syscall"
	"unsafe"
)

// IsTerminal reports whether the descriptor is a terminal, by the TIOCGETA
// ioctl, without cgo.
func IsTerminal(fd uintptr) bool {
	var termios syscall.Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TIOCGETA, uintptr(unsafe.Pointer(&termios)))
	return errno == 0
}
