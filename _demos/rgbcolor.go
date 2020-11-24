package main

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// This example should demonstrate the functionality of full rgb-support,
// as well as the ability to combine rgb colors and (multiple) attributes.

var fgR uint8 = 150
var fgG uint8 = 100
var fgB uint8 = 50

var bgR uint8 = 50
var bgG uint8 = 100
var bgB uint8 = 150

var currentBold bool = true
var currentUnderline bool = false
var currentReverse bool = false
var currentRGB bool = true
var currentCursive bool = false
var currentHidden bool = false
var currentBlink bool = false
var currentDim bool = false

var boolLabel map[bool]string = make(map[bool]string)

const preview string = " Here is some example text "
const padding string = "                           "

const coldef = termbox.ColorDefault

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func redraw_all() {
	tbprint(20, 1, coldef, coldef, " - Current Settings - ")

	var r, g, b string
	r = fmt.Sprintf("%3d", fgR)
	g = fmt.Sprintf("%3d", fgG)
	b = fmt.Sprintf("%3d", fgB)
	tbprint(4, 3, coldef, coldef, "Foreground Red:")
	tbprint(5, 4, coldef, coldef, "[h] "+r+" [l]")
	tbprint(4, 5, coldef, coldef, "Foreground Green:")
	tbprint(5, 6, coldef, coldef, "[j] "+g+" [k]")
	tbprint(4, 7, coldef, coldef, "Foreground Blue:")
	tbprint(5, 8, coldef, coldef, "[u] "+b+" [i]")

	r = fmt.Sprintf("%3d", bgR)
	g = fmt.Sprintf("%3d", bgG)
	b = fmt.Sprintf("%3d", bgB)
	tbprint(23, 3, coldef, coldef, "Background Red:")
	tbprint(24, 4, coldef, coldef, "[H] "+r+" [L]")
	tbprint(23, 5, coldef, coldef, "Background Green:")
	tbprint(24, 6, coldef, coldef, "[J] "+g+" [K]")
	tbprint(23, 7, coldef, coldef, "Background Blue:")
	tbprint(24, 8, coldef, coldef, "[U] "+b+" [I]")

	var bold, ul, rev, rgb, cur, hid, blink, dim string
	bold = boolLabel[currentBold]
	ul = boolLabel[currentUnderline]
	rev = boolLabel[currentReverse]
	rgb = boolLabel[currentRGB]
	cur = boolLabel[currentCursive]
	hid = boolLabel[currentHidden]
	blink = boolLabel[currentBlink]
	dim = boolLabel[currentDim]

	tbprint(42, 3, coldef, coldef, "Bold:")
	tbprint(43, 4, coldef, coldef, bold+" [w]")
	tbprint(42, 5, coldef, coldef, "Underline:")
	tbprint(43, 6, coldef, coldef, ul+" [a]")
	tbprint(42, 7, coldef, coldef, "Reverse:")
	tbprint(43, 8, coldef, coldef, rev+" [s]")
	tbprint(42, 9, coldef, coldef, "Full RGB:")
	tbprint(43, 10, coldef, coldef, rgb+" [t]")
	tbprint(54, 3, coldef, coldef, "Cursive:")
	tbprint(55, 4, coldef, coldef, cur+" [d]")
	tbprint(54, 5, coldef, coldef, "Hidden:")
	tbprint(55, 6, coldef, coldef, hid+" [e]")
	tbprint(54, 7, coldef, coldef, "Blink:")
	tbprint(55, 8, coldef, coldef, blink+" [r]")
	tbprint(54, 9, coldef, coldef, "Dim:")
	tbprint(55, 10, coldef, coldef, dim+" [f]")

	tbprint(20, 12, coldef, coldef, "Quit with [q] or [ESC]")
	tbprint(6, 13, coldef, coldef, "Note that RGB may be incompatible with other modifiers")

	var fg, bg termbox.Attribute
	if currentRGB {
		termbox.SetOutputMode(termbox.OutputRGB)
		fg = termbox.RGBToAttribute(uint8(fgR), uint8(fgG), uint8(fgB))
		bg = termbox.RGBToAttribute(uint8(bgR), uint8(bgG), uint8(bgB))
	} else {
		termbox.SetOutputMode(termbox.OutputNormal)
		fg = termbox.ColorRed
		bg = termbox.ColorDefault
	}
	tfg := fg // tfg are the attributes that should be applied to the text
	if currentBold {
		tfg |= termbox.AttrBold
	}
	if currentUnderline {
		tfg |= termbox.AttrUnderline
	}
	if currentReverse {
		fg |= termbox.AttrReverse
		tfg |= termbox.AttrReverse
	}
	if currentCursive {
		tfg |= termbox.AttrCursive
	}
	if currentHidden {
		fg |= termbox.AttrHidden
		tfg |= termbox.AttrHidden
	}
	if currentBlink {
		fg |= termbox.AttrBlink
		tfg |= termbox.AttrBlink
	}
	if currentDim {
		fg |= termbox.AttrDim
		tfg |= termbox.AttrDim
	}
	tbprint(18, 15, fg, bg, padding)
	tbprint(18, 16, tfg, bg, preview)
	tbprint(18, 17, fg, bg, padding)

	termbox.Flush()
}

func main() {
	boolLabel[false] = "Off"
	boolLabel[true] = "On "

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	redraw_all()
mainloop:
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			default:
				switch ev.Ch {
				case 'q', 'Q':
					break mainloop
				case 'h':
					fgR--
				case 'l':
					fgR++
				case 'j':
					fgG--
				case 'k':
					fgG++
				case 'u':
					fgB--
				case 'i':
					fgB++
				case 'H':
					bgR--
				case 'L':
					bgR++
				case 'J':
					bgG--
				case 'K':
					bgG++
				case 'U':
					bgB--
				case 'I':
					bgB++
				case 'w', 'W':
					currentBold = !currentBold
				case 'a', 'A':
					currentUnderline = !currentUnderline
				case 's', 'S':
					currentReverse = !currentReverse
				case 't', 'T':
					currentRGB = !currentRGB
				case 'd', 'D':
					currentCursive = !currentCursive
				case 'e', 'E':
					currentHidden = !currentHidden
				case 'r', 'R':
					currentBlink = !currentBlink
				case 'f', 'F':
					currentDim = !currentDim
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redraw_all()
	}
}
