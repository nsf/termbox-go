package main

import "github.com/nsf/termbox-go"

import "fmt"
import "time"
import "math/rand"

// will display a lot of colored letters.
// you can change the background color using the arrow keys
// the foreground color will change randomly over time

var letters = []rune{'o', 'x', 'i', 'n', 'u', 's', ' '}

var color int

func main() {
	err := termbox.Init()

	if err != nil {
		fmt.Println(err)
		return
	}

	w, h := termbox.Size()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			termbox.SetChar(x, y, letters[rand.Intn(len(letters))])
		}
	}
	termbox.Flush()

	go bgthread()

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			if ev.Ch == 'q' || ev.Key == termbox.KeyEsc {
				break
			} else if ev.Ch == 'h' || ev.Key == termbox.KeyArrowLeft {
				color--
			} else if ev.Ch == 'l' || ev.Key == termbox.KeyArrowRight {
				color++
			}
			for color < 0 {
				color += 9
			}
			color %= 9
			fillbg(termbox.Attribute(color))
			termbox.Flush()
		}
	}

	termbox.Close()
}

func fillbg(bg termbox.Attribute) {
	w, h := termbox.Size()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			termbox.SetBg(x, y, bg)
		}
	}
}

func bgthread() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		w, h := termbox.Size()
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				termbox.SetFg(x, y, termbox.Attribute(rand.Intn(9)))
			}
		}
		termbox.Flush()
		<-ticker.C // wait for ticker
	}
}
