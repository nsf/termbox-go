package main

// import "fmt"
// import "github.com/nsf/termbox-go"
import "github.com/kubaroth/termbox-go"

import "time"

func draw(_x int, _y int) {
    w, h := termbox.Size()
    for y := 0; y < h; y++ {
        for x := 0; x < w; x++ {
            if _x == x && _y == y {
            termbox.SetCell(x, y, ' ', termbox.ColorDefault, 3)
            }
        }
    }
    termbox.Flush()
}

func main() {
    err := termbox.Init()
    if err != nil {
        panic(err)
    }
    defer termbox.Close()
    termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)  

loop:
    for {
        ev := termbox.PollEvent(); 

        if ev.Type == termbox.EventMouse && ev.Key == termbox.MouseLeft{
            draw(ev.MouseX, ev.MouseY)
            
        } else if  ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
            break loop
        }
        
        time.Sleep(5 * time.Millisecond)
    }
}