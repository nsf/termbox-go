package main

import "github.com/nsf/termbox-go"
import "fmt"

// This program can demonstrate the 16 basic colors available
// for foreground and background.

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += 1
	}
}

func main() {
	termbox.Init()

	var i, j int
	var fg, bg termbox.Attribute
	var colorRange []termbox.Attribute = []termbox.Attribute{
		termbox.ColorDefault,
		termbox.ColorBlack,
		termbox.ColorRed,
		termbox.ColorGreen,
		termbox.ColorYellow,
		termbox.ColorBlue,
		termbox.ColorMagenta,
		termbox.ColorCyan,
		termbox.ColorWhite,
		termbox.ColorDarkGray,
		termbox.ColorLightRed,
		termbox.ColorLightGreen,
		termbox.ColorLightYellow,
		termbox.ColorLightBlue,
		termbox.ColorLightMagenta,
		termbox.ColorLightCyan,
		termbox.ColorLightGray,
	}

	var row, col int
	var text string
	for i, fg = range colorRange {
		for j, bg = range colorRange {
			row = i + 1
			col = j * 8
			text = fmt.Sprintf(" %02d/%02d ", fg, bg)
			tbprint(col, row+0, fg, bg, text)
			/*text = fmt.Sprintf(" on ")
			tbprint(col, row+1, fg, bg, text)
			text = fmt.Sprintf(" %2d ", bg)
			tbprint(col, row+2, fg, bg, text)*/
			//fmt.Println(text, col, row)
		}
	}
	for j, bg = range colorRange {
		tbprint(j*8, 0, termbox.ColorDefault, bg, "       ")
		tbprint(j*8, i+2, termbox.ColorDefault, bg, "       ")
	}

	tbprint(15, i+4, termbox.ColorDefault, termbox.ColorDefault,
		"Press any key to close...")
	termbox.Flush()
	termbox.PollEvent()
	termbox.Close()
}
