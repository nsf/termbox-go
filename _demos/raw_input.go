package main

import (
	"github.com/nsf/termbox-go"
	"fmt"
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

var current string

func redraw_all() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	tbprint(0, 0, coldef, coldef, current)
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

mainloop:
	for {
		var data [64]byte
		switch ev := termbox.PollRawEvent(data[:]); ev.Type {
		case termbox.EventRaw:
			d := data[:ev.N]
			current = fmt.Sprintf("%q", d)
			if current == `"\x1b"` {
				break mainloop
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redraw_all()
	}
}
