package main

import "os"
import "time"
import "fmt"

import "github.com/nsf/termbox-go"
import "github.com/mattn/go-runewidth"

func main() {
	err := termbox.Init()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tbprint(2, 2, termbox.ColorRed, termbox.ColorDefault, "Hello terminal!")
	termbox.Flush()

	time.Sleep(time.Second)
	termbox.Close()
}

// This function is often useful:
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}
