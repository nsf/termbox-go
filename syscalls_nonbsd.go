//go:build !darwin && !freebsd && !netbsd && !openbsd && !windows

package termbox

import (
	"golang.org/x/sys/unix"
)

const (
	tcgets = unix.TCGETS
	tcsets = unix.TCSETS
)
