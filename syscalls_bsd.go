//go:build darwin || freebsd || openbsd || netbsd

package termbox

import (
	"golang.org/x/sys/unix"
)

const (
	tcgets = unix.TIOCGETA
	tcsets = unix.TIOCSETA
)
