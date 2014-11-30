package main

import "github.com/nsf/termbox-go"

const chars = "nnnnnnnnnbbbbbbbbbuuuuuuuuuBBBBBBBBB"

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
	print_combinations_table(1, 1, []termbox.Attribute{
		0,
		termbox.AttrBold,
	})
	print_combinations_table(2+len(chars), 1, []termbox.Attribute{
		termbox.AttrReverse,
	})
	print_wide(2+len(chars), 11, hello_world)
	termbox.Flush()

	termbox.SetOutputMode(termbox.OutputGrayscale)
	var x, y, c int
	for x, y = 0, 23; x < 24; x++ {
		termbox.SetCell(x, y, '@', termbox.Attribute(x), 0)
		termbox.SetCell(x+25, y, ' ', 0, termbox.Attribute(x))
	}
	termbox.Flush()

	termbox.SetOutputMode(termbox.Output216)
	y++
	for c, x = 0, 0; c < 216; c, x = c+1, x+1 {
		if x % 24 == 0 {
			x = 0
			y++
		}
		termbox.SetCell(x, y, '@', termbox.Attribute(c), 0)
		termbox.SetCell(x+25, y, ' ', 0, termbox.Attribute(c))
	}
	termbox.Flush()

	termbox.SetOutputMode(termbox.Output256)
	y++
	for c, x = 0, 0; c < 256; c, x = c+1, x+1 {
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
	termbox.Flush()

	termbox.SetOutputMode(termbox.OutputNormal)
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
			}
		case termbox.EventResize:
			draw_all()
		}
	}
}
