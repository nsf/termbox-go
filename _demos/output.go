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
