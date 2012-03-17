package main

import "github.com/nsf/termbox-go"
import "math/rand"
import "time"

func draw() {
	w, h := termbox.Size()
	termbox.Clear()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			termbox.ChangeCell(x, y, ' ', termbox.ColorWhite,
				termbox.Attribute(rand.Int() % 8))
		}
	}
	termbox.Present()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Shutdown()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			ev := termbox.PollEvent()
			event_queue <- ev
		}
	}()

	draw()
loop:
	for {
		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
		default:
			draw()
			time.Sleep(10 * time.Millisecond)
		}
	}
}
