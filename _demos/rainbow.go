package main

import "../../termbox-go"

var x, y int

func print_str(fg, bg termbox.Attribute, s string) {
	for _, r := range s {
		termbox.SetCell(x, y, r, fg, bg)
		x++
	}
}

func print_cell(bg termbox.Attribute) {
	print_str(0, bg, "  ")
}

func vt100(c int) termbox.Attribute {
	return termbox.Attribute(c)
}

func rgb(r, g, b int) termbox.Attribute {
	return termbox.Attribute(r*36 + g*6 + b + 16)
}

func grey(v int) termbox.Attribute {
	return termbox.Attribute(v + 232)
}

func print_head(msg string) {
	x, y = 0, y+2
	print_str(15, 0, msg)
	x, y = 4, y+2
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	err = termbox.SetColorMode(termbox.ColorMode256)
	print_str(3, 0, "Termbox 256 Color Rainbow ")
	if err != nil {
		x, y = 0, y+1
		print_str(2, 0, err.Error())
	}

	print_head("ANSI Normal Intensity:")
	for i := 0; i < 8; i++ {
		print_cell(vt100(i))
	}

	print_head("ANSI Bright Intensity:")
	for i := 8; i < 16; i++ {
		print_cell(vt100(i))
	}

	print_head("6x6x6 Color Cube:")
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				print_cell(rgb(r, g, b))
			}
			print_str(0, 0, "  ")
		}
		x, y = 4, y+1
	}
	y --

	print_head("24 Value Greyscale:")
	for i := 0; i < 24; i++ {
		print_cell(grey(i))
	}

	for {
		termbox.Flush()
		if termbox.PollEvent().Type == termbox.EventKey {
			break
		}
	}
}
