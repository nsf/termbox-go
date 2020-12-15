package main

import "os"
import "fmt"

import "github.com/nsf/termbox-go"
import "github.com/mattn/go-runewidth"

func main() {
	err := termbox.Init()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	termbox.SelectScreen(false)

	tbprint(2, 2, termbox.ColorRed, termbox.ColorDefault, "Hello terminal!")
	tbprint(2, 4, termbox.ColorBlue, termbox.ColorDefault, "Here is some text.")
	termbox.Flush()
	termbox.PollEvent()

	tbprint(2, 0, termbox.ColorBlue, termbox.ColorDefault, "This is what you'd see after flush")
	termbox.Flush()
	termbox.PollEvent()

	termbox.Close()
}

// This function is often useful:
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}
