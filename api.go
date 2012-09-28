// +build !windows

package termbox

import "errors"
import "fmt"
import "os"
import "os/signal"
import "strings"
import "syscall"

// public API

// Initializes termbox library. This function should be called before any other functions.
// After successful initialization, the library must be finalized using 'Close' function.
//
// Example usage:
//      err := termbox.Init()
//      if err != nil {
//              panic(err)
//      }
//      defer termbox.Close()
func Init() error {
	// TODO: try os.Stdin and os.Stdout directly
	var err error

	// os.Create is confusing here, but it's just a shortcut for 'open'
	out, err = os.Create("/dev/tty")
	if err != nil {
		return err
	}
	in, err = os.Open("/dev/tty")
	if err != nil {
		return err
	}

	err = setup_term()
	if err != nil {
		return err
	}

	// we set two signal handlers, because input/output are not really
	// connected, but they both need to be aware of window size changes
	signal.Notify(sigwinch, syscall.SIGWINCH)

	err = tcgetattr(out.Fd(), &orig_tios)
	if err != nil {
		return err
	}

	tios := orig_tios
	tios.Iflag &^= syscall_IGNBRK | syscall_BRKINT | syscall_PARMRK |
		syscall_ISTRIP | syscall_INLCR | syscall_IGNCR |
		syscall_ICRNL | syscall_IXON
	tios.Oflag &^= syscall_OPOST
	tios.Lflag &^= syscall_ECHO | syscall_ECHONL | syscall_ICANON |
		syscall_ISIG | syscall_IEXTEN
	tios.Cflag &^= syscall_CSIZE | syscall_PARENB
	tios.Cflag |= syscall_CS8
	tios.Cc[syscall_VMIN] = 1
	tios.Cc[syscall_VTIME] = 0

	err = tcsetattr(out.Fd(), &tios)
	if err != nil {
		return err
	}

	out.WriteString(funcs[t_enter_ca])
	out.WriteString(funcs[t_enter_keypad])
	out.WriteString(funcs[t_hide_cursor])
	out.WriteString(funcs[t_clear_screen])

	termw, termh = get_term_size(out.Fd())
	back_buffer.init(termw, termh)
	front_buffer.init(termw, termh)
	back_buffer.clear()
	front_buffer.clear()

	go func() {
		buf := make([]byte, 128)
		for {
			n, err := in.Read(buf)
			input_comm <- input_event{buf[:n], err}
			ie := <-input_comm
			buf = ie.data[:128]
		}
	}()

	return nil
}

// used to construct palettes from 24-bit RGB values
type RGB struct{ R, G, B byte }

// used to load various color palettes for 256-color terminals
func SetColorPalette(p []RGB) {
	for n, c := range p {
		out.WriteString(fmt.Sprintf("\033]4;%v;rgb:%2x/%2x/%2x\x1b\\", n, c.R, c.G, c.B))
	}
}

// a preconfigured palette corresponding to XTERM's defaults
var Palette256 []RGB

func init() {
	var r, g, b byte

	// initialize the standard 256 color palette
	Palette256 = make([]RGB, 256)

	// this is the default xterm palette for the first 16 colors
	// pay attention to the blues, which are increased in luma 
	// to compensate for human insensitivity to cooler colors

	r, g, b = 205, 205, 238
	Palette256[0] = RGB{0, 0, 0}
	Palette256[1] = RGB{r, 0, 0}
	Palette256[2] = RGB{0, g, 0}
	Palette256[3] = RGB{r, g, 0}
	Palette256[4] = RGB{0, 0, b}
	Palette256[5] = RGB{r, 0, b}
	Palette256[6] = RGB{0, g, b}
	Palette256[7] = RGB{r, g, b}

	r, g, b = 255, 255, 255
	Palette256[8] = RGB{127, 127, 127}
	Palette256[9] = RGB{r, 0, 0}
	Palette256[10] = RGB{0, g, 0}
	Palette256[11] = RGB{r, g, 0}
	Palette256[12] = RGB{92, 92, b}
	Palette256[13] = RGB{r, 0, b}
	Palette256[14] = RGB{0, g, b}
	Palette256[15] = RGB{r, g, b}


	// next we establish a 6x6x6 color cube with no blue
	// correction -- also xterm common, to the point that
	// many users think this is hardcoded
	c := 16

	for r = 0; r < 6; r++ {
		rr := r * 40
		if r > 0 {
			rr += 55
		}
		for g = 0; g < 6; g++ {
			gg := g * 40
			if g > 0 {
				gg += 55
			}
			for b = 0; b < 6; b++ {
				bb := b * 40
				if b > 0 {
					bb += 55
				}
				Palette256[c] = RGB{rr, gg, bb}
				c++
			}
		}
	}

	// and, following the user assumptions, this is
	// a 24 color grey ramp
	var v byte = 8
	for g := 0; g < 24; g++ {
		v += 10
		Palette256[c] = RGB{v, v, v}
		c++
	}
}

// instructs termbox to switch to either ColorMode16 or ColorMode256 
func SetColorMode(cm ColorMode) error {
	switch cm {
	case ColorMode16:
		color_mode = cm
		return nil
	case ColorMode256:
		// let it fall through, we need to examine $TERM
	default:
		return errors.New("termbox: invalid color mode requested")
	}

	term := os.Getenv("TERM")
	switch {
	case term == "":
		return errors.New("termbox: TERM environment variable not set")
	case strings.Index(term, "256") == -1:
		return errors.New("termbox: TERM does not contain \"256\"")
	}

	// this is the common palette expected by xterm-256 hackers; it is 
	// NOT the only possible one, and a SetColorPalette command might
	// be in order..

	color_mode = cm
	SetColorPalette(Palette256)
	return nil
}

// Finalizes termbox library, should be called after successful initialization
// when termbox's functionality isn't required anymore.
func Close() {
	out.WriteString(funcs[t_show_cursor])
	out.WriteString(funcs[t_sgr0])
	out.WriteString(funcs[t_clear_screen])
	out.WriteString(funcs[t_exit_ca])
	out.WriteString(funcs[t_exit_keypad])
	tcsetattr(out.Fd(), &orig_tios)

	// I don't close them, becase on darwin a file descriptor which is
	// blocked in one thread in a read call, gets blocked here as well and
	// that prevents termbox from shutting down without getting additional
	// input. Honestly there are issues which prevent multiple Init/Close
	// calls within the same program anyway, so, let's just leave them open,
	// OS will clean them up for us. Although correct behaviour will be
	// implemented one day.

	/*
		out.Close()
		in.Close()
	*/
}

// Synchronizes the internal back buffer with the terminal.
func Flush() error {
	// invalidate cursor position
	lastx = coord_invalid
	lasty = coord_invalid

	update_size_maybe()

	for y := 0; y < front_buffer.height; y++ {
		line_offset := y * front_buffer.width
		for x := 0; x < front_buffer.width; x++ {
			cell_offset := line_offset + x
			back := &back_buffer.cells[cell_offset]
			front := &front_buffer.cells[cell_offset]
			if *back == *front {
				continue
			}
			send_attr(back.Fg, back.Bg)
			send_char(x, y, back.Ch)
			*front = *back
		}
	}
	if !is_cursor_hidden(cursor_x, cursor_y) {
		write_cursor(cursor_x, cursor_y)
	}
	return flush()
}

// Sets the position of the cursor. See also HideCursor().
func SetCursor(x, y int) {
	if is_cursor_hidden(cursor_x, cursor_y) && !is_cursor_hidden(x, y) {
		outbuf.WriteString(funcs[t_show_cursor])
	}

	if !is_cursor_hidden(cursor_x, cursor_y) && is_cursor_hidden(x, y) {
		outbuf.WriteString(funcs[t_hide_cursor])
	}

	cursor_x, cursor_y = x, y
	if !is_cursor_hidden(cursor_x, cursor_y) {
		write_cursor(cursor_x, cursor_y)
	}
}

// The shortcut for SetCursor(-1, -1).
func HideCursor() {
	SetCursor(cursor_hidden, cursor_hidden)
}

// Changes cell's parameters in the internal back buffer at the specified
// position.
func SetCell(x, y int, ch rune, fg, bg Attribute) {
	if x < 0 || x >= back_buffer.width {
		return
	}
	if y < 0 || y >= back_buffer.height {
		return
	}

	back_buffer.cells[y*back_buffer.width+x] = Cell{ch, fg, bg}
}

// Returns a slice into the termbox's back buffer. You can get its dimensions
// using 'Size' function. The slice remains valid as long as no 'Clear' or
// 'Flush' function calls were made after call to this function.
func CellBuffer() []Cell {
	return back_buffer.cells
}

// Wait for an event and return it. This is a blocking function call.
func PollEvent() Event {
	var event Event

	// try to extract event from input buffer, return on success
	event.Type = EventKey
	if extract_event(&event) {
		return event
	}

	for {
		select {
		case ev := <-input_comm:
			if ev.err != nil {
				return Event{Type: EventError, Err: ev.err}
			}

			inbuf = append(inbuf, ev.data...)
			input_comm <- ev
			if extract_event(&event) {
				return event
			}
		case <-sigwinch:
			event.Type = EventResize
			event.Width, event.Height = get_term_size(out.Fd())
			return event
		}
	}
	panic("unreachable")
}

// Returns the size of the internal back buffer (which is the same as
// terminal's window size in characters).
func Size() (int, int) {
	return termw, termh
}

// Clears the internal back buffer.
func Clear(fg, bg Attribute) error {
	foreground, background = fg, bg
	err := update_size_maybe()
	back_buffer.clear()
	return err
}

// Sets termbox input mode. Termbox has two input modes:
//
// 1. Esc input mode. When ESC sequence is in the buffer and it doesn't match
// any known sequence. ESC means KeyEsc.
//
// 2. Alt input mode. When ESC sequence is in the buffer and it doesn't match
// any known sequence. ESC enables ModAlt modifier for the next keyboard event.
//
// If 'mode' is InputCurrent, returns the current input mode. See also Input*
// constants.
func SetInputMode(mode InputMode) InputMode {
	if mode != InputCurrent {
		input_mode = mode
	}
	return input_mode
}
