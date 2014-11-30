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
		for y := 0; y < 24; y++ {
			for x:= 0; x < 24; x++ {
				termbox.SetCell(x, y, 'n',
					termbox.Attribute(x),
					termbox.Attribute(y))
				termbox.SetCell(x+25, y, 'b',
					termbox.Attribute(x) | termbox.AttrBold,
					termbox.Attribute(23 - y))
				termbox.SetCell(x+50, y, 'u',
					termbox.Attribute(x) | termbox.AttrUnderline,
					termbox.Attribute(y))
			}
		}

	case termbox.Output216:
		for r := 0; r < 6; r++ {
			for g := 0; g < 6; g++ {
				for b := 0; b < 6; b++ {
					y := r
					x := g + 6 * b
					c1 := termbox.Attribute(r*36 + g*6 + b)
					bg := termbox.Attribute(g*36 + b*6 + r)
					c2 := termbox.Attribute(b*36 + r*6 + g)
					bc1 := c1 | termbox.AttrBold
					uc1 := c1 | termbox.AttrUnderline
					bc2 := c2 | termbox.AttrBold
					uc2 := c2 | termbox.AttrUnderline
					termbox.SetCell(x, y, 'n', c1, bg)
					termbox.SetCell(x, y + 6, 'b', bc1, bg)
					termbox.SetCell(x, y + 12, 'u', uc1, bg)
					termbox.SetCell(x, y + 18, 'B', bc1 | uc1, bg)
					termbox.SetCell(x + 37, y, 'n', c2, bg)
					termbox.SetCell(x + 37, y + 6, 'b', bc2, bg)
					termbox.SetCell(x + 37, y + 12, 'u', uc2, bg)
					termbox.SetCell(x + 37, y + 18, 'B', bc2 | uc2, bg)
				}
			}
		}

	case termbox.Output256:
		for y := 0; y < 4; y++ {
			for x := 0; x < 8; x++ {
				for z := 0; z < 8; z++ {
					bg := termbox.Attribute(y * 64 + x * 8 + z)
					c1 := termbox.Attribute(255 - y*64 - x*8 - z)
					c2 := termbox.Attribute(y*64 + z*8 + x)
					c3 := termbox.Attribute(255 - y*64 - z*8 - x)
					c4 := termbox.Attribute(y*64 + x*4 + z*4)
					termbox.SetCell(z + 8*x, y, ' ', 0, bg)
					termbox.SetCell(z + 8*x, y+5, 'n', c4, bg)
					termbox.SetCell(z + 8*x, y+10, 'b', c2, bg)
					termbox.SetCell(z + 8*x, y+15, 'u', c3, bg)
					termbox.SetCell(z + 8*x, y+20, 'B', c1, bg)
				}
			}
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
