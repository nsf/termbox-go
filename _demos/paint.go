package main

// import "fmt"
// import "github.com/nsf/termbox-go"
import "github.com/kubaroth/termbox-go"

import "time"

func draw(_x int, _y int, clr termbox.Attribute) {
    w, h := termbox.Size()
    for y := 0; y < h; y++ {
        for x := 0; x < w; x++ {
            if _x == x && _y == y {
            termbox.SetCell(x, y, ' ', termbox.ColorDefault, clr)
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
        // Draw
        if ev.Type == termbox.EventMouse && ev.Key == termbox.MouseLeft || ev.Key == termbox.MouseDown{
            draw(ev.MouseX, ev.MouseY, 3)
        // Erase
        } else if ev.Type == termbox.EventMouse && ev.Key == termbox.MouseAltLeft || ev.Key == termbox.MouseAltDown{
            draw(ev.MouseX, ev.MouseY, 0)
        // Exit
        } else if  ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
            break loop
        }
        
        time.Sleep(5 * time.Millisecond)
    }
}