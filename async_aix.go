//go:build aix

package termbox

import (
	"os"

	"golang.org/x/sys/unix"
)

const (
	fioasync  = -2147195267
	fiosetown = -2147195268
)

func setFdAsync(fd int) error {
	err := unix.IoctlSetPointerInt(fd, fiosetown, os.Getpid())
	if err != nil {
		return err
	}

	err = unix.IoctlSetPointerInt(fd, fioasync, 1)
	if err != nil {
		return err
	}
	return nil
}
