package main

import (
	_ "unicode/utf8"

	_ "github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"

	"fmt"
)

// This gives a table of the 256-color-set,
// both the foreground and background variants.
// It is ordered to produce many color ramps.
// Quit with ESC

func draw_ramp() {
	var i int
	for i = 0; i < 256; i++ {
		row := ((i + 2) / 6) + 3
		col := ((i + 2) % 6) * 4
		var text string = fmt.Sprintf("%03d", i)
		for j := 0; j < 3; j++ {
			termbox.SetCell(col+j, row, []rune(text)[j],
				termbox.Attribute(i+1), termbox.ColorDefault)
			termbox.SetCell(col+j+30, row, []rune(text)[j],
				termbox.ColorDefault, termbox.Attribute(i+1))
		}
	}
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.Output256)

	draw_ramp()

	termbox.PollEvent()
	termbox.Close()
}
