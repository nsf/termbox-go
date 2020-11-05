package main

import (
	"github.com/nsf/termbox-go"
	"strconv"
)

var currentR int = 150
var currentG int = 100
var currentB int = 50

func draw_red_text() {
	termbox.SetCell(5, 5, 'H', termbox.ColorRed, termbox.ColorDefault)
	termbox.Flush()
}

func draw_colored_text() {
	termbox.SetOutputMode(termbox.OutputRGB)
	fg := termbox.RGBToAttribute(uint8(currentR), uint8(currentG), uint8(currentB))
	bg := termbox.RGBToAttribute(5, 5, 5)
	//panic("old color is " + fmt.Sprint(uint32(attr)))
	// 9779300
	for i, v := range "Here is some rgb text in #" + strconv.Itoa(currentR) + ";" + strconv.Itoa(currentG) + ";" + strconv.Itoa(currentB) + ";   " {
		termbox.SetCell(5+i, 5, v, fg, bg)
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

	draw_colored_text()
	//draw_red_text()
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				currentR = currentR - 10
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				currentR = currentR + 10
			case termbox.KeyArrowUp:
				currentG = currentG - 10
			case termbox.KeyArrowDown:
				currentG = currentG + 10
			default:
				if ev.Ch == 'q' {
					break mainloop
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		draw_colored_text()
		//draw_red_text()
	}
}
