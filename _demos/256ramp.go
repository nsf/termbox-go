package main

import (
	_ "unicode/utf8"

	_ "github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"

	"fmt"

	"bufio"
	"os"
)

func draw_ramp() {
	log_file, _ := os.OpenFile("/tmp/rgb.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	log_writer := bufio.NewWriter(log_file)
	termbox.SetOutputMode(termbox.Output256)
	var i int
	for i = 0; i < 256; i++ {
		row := ((i + 2) / 6) + 3
		col := ((i + 2) % 6) * 4
		//if i < 20 {
		log_writer.WriteString(fmt.Sprintln("printing", i, "to", row, "x", col))
		//}
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
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	draw_ramp()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		draw_ramp()
	}
}
