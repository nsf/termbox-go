package main

import "github.com/nsf/termbox-go"

const chars = "nnnnnnnnnbbbbbbbbbuuuuuuuuuBBBBBBBBB"
var output_mode = termbox.OutputNormal

func next_char(current int) int {
	current++
	if current >= len(chars) {
		return 0
	}
	return current
}

func print_combinations_table(sx, sy int, attrs []termbox.Attribute) {
	var bg termbox.Attribute
	current_char := 0
	y := sy

	all_attrs := []termbox.Attribute{
		0,
		termbox.AttrBold,
		termbox.AttrUnderline,
		termbox.AttrBold | termbox.AttrUnderline,
	}

	draw_line := func() {
		x := sx
		for _, a := range all_attrs {
			for c := termbox.ColorDefault; c <= termbox.ColorWhite; c++ {
				fg := a | c
				termbox.SetCell(x, y, rune(chars[current_char]), fg, bg)
				current_char = next_char(current_char)
				x++
			}
		}
	}

	for _, a := range attrs {
		for c := termbox.ColorDefault; c <= termbox.ColorWhite; c++ {
			bg = a | c
			draw_line()
			y++
		}
	}
}

func print_wide(x, y int, s string) {
	red := false
	for _, r := range s {
		c := termbox.ColorDefault
		if red {
			c = termbox.ColorRed
		}
		termbox.SetCell(x, y, r, termbox.ColorDefault, c)
		x += 2
		red = !red
	}
}

const hello_world = "こんにちは世界"

func draw_all() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	switch output_mode {

	case termbox.OutputNormal:
		print_combinations_table(1, 1, []termbox.Attribute{
			0,
			termbox.AttrBold,
		})
		print_combinations_table(2+len(chars), 1, []termbox.Attribute{
			termbox.AttrReverse,
		})
		print_wide(2+len(chars), 11, hello_world)

	case termbox.OutputGrayscale:
		for x, y := 0, 0; x < 24; x++ {
			termbox.SetCell(x, y, '@', termbox.Attribute(x), 0)
			termbox.SetCell(x+25, y, ' ', 0, termbox.Attribute(x))
		}

	case termbox.Output216:
		for x, y, c := 0, 0, 0; c < 216; c, x = c+1, x+1 {
			if x % 24 == 0 {
				x = 0
				y++
			}
			termbox.SetCell(x, y, '@', termbox.Attribute(c), 0)
			termbox.SetCell(x+25, y, ' ', 0, termbox.Attribute(c))
		}

	case termbox.Output256:
		for x, y, c := 0, 0, 0; c < 256; c, x = c+1, x+1 {
			if x % 24 == 0 {
				x = 0
				y++
			}
			if y & 1 != 0 {
				termbox.SetCell(x, y, '+',
						termbox.Attribute(c) | termbox.AttrUnderline, 0)
			} else {
				termbox.SetCell(x, y, '+', termbox.Attribute(c), 0)
			}
			termbox.SetCell(x+25, y, ' ', 0, termbox.Attribute(c))
		}

	}

	termbox.Flush()
}

var available_modes = []termbox.OutputMode {
	termbox.OutputNormal,
	termbox.OutputGrayscale,
	termbox.Output216,
	termbox.Output256,
}

var output_mode_index = 0

func switch_output_mode(direction int) {
	output_mode_index += direction
	if output_mode_index < 0 {
		output_mode_index = len(available_modes) - 1
	} else if output_mode_index >= len(available_modes) {
		output_mode_index = 0
	}
	output_mode = available_modes[output_mode_index]
	termbox.SetOutputMode(output_mode)
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	draw_all()
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			case termbox.KeyArrowUp, termbox.KeyArrowRight:
				switch_output_mode(1)
				draw_all()
			case termbox.KeyArrowDown, termbox.KeyArrowLeft:
				switch_output_mode(-1)
				draw_all()
			}
		case termbox.EventResize:
			draw_all()
		}
	}
}
