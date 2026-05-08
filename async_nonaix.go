//go:build !aix && !windows

package termbox

import (
	"runtime"

	"golang.org/x/sys/unix"
)

func setFdAsync(fd int) error {
	_, err := unix.FcntlInt(uintptr(fd), unix.F_SETFL, unix.O_ASYNC|unix.O_NONBLOCK)
	if err != nil {
		return err
	}
	_, err = unix.FcntlInt(uintptr(fd), unix.F_SETOWN, unix.Getpid())
	if runtime.GOOS != "darwin" && err != nil {
		return err
	}
	return nil
}
