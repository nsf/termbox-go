module github.com/nsf/termbox-go

go 1.23.0

require (
	github.com/mattn/go-runewidth v0.0.19
	golang.org/x/sys v0.31.0
)

require github.com/clipperhouse/uax29/v2 v2.2.0 // indirect

retract v1.1.0 // panics on BSD
