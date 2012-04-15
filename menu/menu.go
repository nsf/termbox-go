package termbox

import (
	"github.com/nsf/termbox-go"
	"strings"
)

const (
	TOPLEFT     = "┌"
	TOPRIGHT    = "┐"
	BOTTOMLEFT  = "└"
	BOTTOMRIGHT = "┘"
	VERTICAL    = "│"
	HORIZONTAL  = "─"
)

func DrawText(x, y int, text string) {
	j := 0
	for _, r := range text {
		termbox.SetCell(x+j, y, r, termbox.ColorDefault, termbox.ColorDefault)
		j += 1
	}
}

func DrawTextMulti(x, y int, text string) {
	lines := strings.SplitAfterN(text, "\n", -1)
	for i := 0; i < len(lines); i++ {
		DrawText(x, y+i, lines[i])
	}
}

func DrawBox(x, y, width, height int) {
	DrawText(x, y, TOPLEFT+strings.Repeat(HORIZONTAL, width+2)+TOPRIGHT)
	for i := 1; i < height+1; i++ {
		DrawText(x, y+i, VERTICAL+strings.Repeat(" ", width+2)+VERTICAL)
	}
	DrawText(x, y+height, BOTTOMLEFT+strings.Repeat(HORIZONTAL, width+2)+BOTTOMRIGHT)
}

func getDimensions(text string, options []string) (int, int) {
	width := len(text)
	for i := 0; i < len(options); i++ {
		if len(options[i]) > width {
			width = len(options[i])
		}
	}
	height := len(options) + 1
	return width, height
}

func drawOptions(x, y int, options []string, title string) {
	DrawText(x+2, y, title)
	y += 1
	for i := 0; i < len(options); i++ {
		DrawText(x+2, y+i, options[i])
	}
}

func drawSelection(x, y int, sel int) {
	DrawText(x+1, y+1+sel, ">")
}

func DrawMenu(x, y int, text string, options []string) int {
	selection := 0
	width, height := getDimensions(text, options)

loop:
	for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		DrawBox(x, y, width, height)
		drawSelection(x, y, selection)
		drawOptions(x, y, options, text)
		termbox.Flush()
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				break loop
			case termbox.KeyArrowDown:
				if selection < len(options)-1 {
					selection += 1
				}
			case termbox.KeyArrowUp:
				if selection > 0 {
					selection -= 1
				}
			}
		}
	}
	return selection
}
