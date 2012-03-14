package termbox

import "fmt"
import "os"
import "os/signal"
import "syscall"

// public API

type (
	InputMode int
	EventType uint8
	Modifier  uint8
	Key       uint16
	Attribute uint16
)

// This type represents a termbox event. 'Mod', 'Key' and 'Ch' fields are valid
// if 'Type' is EventKey. 'W' and 'H' are valid if 'Type' is EventResize.
type Event struct {
	Type   EventType // one of Event* constants
	Mod    Modifier  // one of Mod* constants or 0
	Key    Key       // one of Key* constants, invalid if 'Ch' is not 0
	Ch     rune      // a unicode character
	Width  int       // width of the screen
	Height int       // height of the screen
}

// A cell, single conceptual entity on the screen. The screen is basically a 2d
// array of cells. 'Ch' is a unicode character, 'Fg' and 'Bg' are foreground
// and background attributes respectively.
type Cell struct {
	Ch rune
	Fg Attribute
	Bg Attribute
}

// Key constants, see Event.Key field.
const (
	KeyF1 Key = 0xFFFF - iota
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyInsert
	KeyDelete
	KeyHome
	KeyEnd
	KeyPgup
	KeyPgdn
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
)

const (
	KeyCtrlTilde      Key = 0x00
	KeyCtrl2          Key = 0x00
	KeyCtrlA          Key = 0x01
	KeyCtrlB          Key = 0x02
	KeyCtrlC          Key = 0x03
	KeyCtrlD          Key = 0x04
	KeyCtrlE          Key = 0x05
	KeyCtrlF          Key = 0x06
	KeyCtrlG          Key = 0x07
	KeyBackspace      Key = 0x08
	KeyCtrlH          Key = 0x08
	KeyTab            Key = 0x09
	KeyCtrlI          Key = 0x09
	KeyCtrlJ          Key = 0x0A
	KeyCtrlK          Key = 0x0B
	KeyCtrlL          Key = 0x0C
	KeyEnter          Key = 0x0D
	KeyCtrlM          Key = 0x0D
	KeyCtrlN          Key = 0x0E
	KeyCtrlO          Key = 0x0F
	KeyCtrlP          Key = 0x10
	KeyCtrlQ          Key = 0x11
	KeyCtrlR          Key = 0x12
	KeyCtrlS          Key = 0x13
	KeyCtrlT          Key = 0x14
	KeyCtrlU          Key = 0x15
	KeyCtrlV          Key = 0x16
	KeyCtrlW          Key = 0x17
	KeyCtrlX          Key = 0x18
	KeyCtrlY          Key = 0x19
	KeyCtrlZ          Key = 0x1A
	KeyEsc            Key = 0x1B
	KeyCtrlLsqBracket Key = 0x1B
	KeyCtrl3          Key = 0x1B
	KeyCtrl4          Key = 0x1C
	KeyCtrlBackslash  Key = 0x1C
	KeyCtrl5          Key = 0x1D
	KeyCtrlRsqBracket Key = 0x1D
	KeyCtrl6          Key = 0x1E
	KeyCtrl7          Key = 0x1F
	KeyCtrlSlash      Key = 0x1F
	KeyCtrlUnderscore Key = 0x1F
	KeySpace          Key = 0x20
	KeyBackspace2     Key = 0x7F
	KeyCtrl8          Key = 0x7F
)

// Alt modifier constant, see Event.Mod field and SetInputMode function.
const (
	ModAlt Modifier = 0x01
)

// Cell attributes, it is possible to use multiple attributes by combining them
// using bitwise OR ('|'). Although, colors cannot be combined. But you can
// combine attributes and a single color.
const (
	ColorBlack Attribute = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorDefault
)

const (
	AttrBold      Attribute = 0x10
	AttrUnderline Attribute = 0x20
)

// Input mode. See SelectInputMode function.
const (
	InputCurrent InputMode = iota
	InputEsc
	InputAlt
)

// Event type. See Event.Type field.
const (
	EventNone EventType = iota
	EventKey
	EventResize
	EventError
)

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
	signal.Notify(sigwinch_input, syscall.SIGWINCH)
	signal.Notify(sigwinch_draw, syscall.SIGWINCH)

	err = tcgetattr(out.Fd(), &orig_tios)
	if err != nil {
		return err
	}

	tios := orig_tios
	tios.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK |
		syscall.ISTRIP | syscall.INLCR | syscall.IGNCR |
		syscall.ICRNL | syscall.IXON
	tios.Oflag &^= syscall.OPOST
	tios.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON |
		syscall.ISIG | syscall.IEXTEN
	tios.Cflag &^= syscall.CSIZE | syscall.PARENB
	tios.Cflag |= syscall.CS8
	tios.Cc[syscall.VMIN] = 1
	tios.Cc[syscall.VTIME] = 0

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
			n, _ := in.Read(buf)
			input_comm <- buf[:n]
			buf = (<-input_comm)[:128]
		}
	}()

	return nil
}

func Shutdown() {
	out.WriteString(funcs[t_show_cursor])
	out.WriteString(funcs[t_sgr0])
	out.WriteString(funcs[t_clear_screen])
	out.WriteString(funcs[t_exit_ca])
	out.WriteString(funcs[t_exit_keypad])
	tcsetattr(out.Fd(), &orig_tios)

	out.Close()
	in.Close()
}

func Present() {
	// invalidate cursor position
	lastx = coord_invalid
	lasty = coord_invalid

	select {
	case <-sigwinch_draw:
		update_size()
	default:
	}

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
		fmt.Fprintf(&outbuf, funcs[t_move_cursor], cursor_y+1, cursor_x+1)
	}
	flush()
}

func SetCursor(x, y int) {
	if is_cursor_hidden(cursor_x, cursor_y) && !is_cursor_hidden(x, y) {
		outbuf.WriteString(funcs[t_show_cursor])
	}

	if !is_cursor_hidden(cursor_x, cursor_y) && is_cursor_hidden(x, y) {
		outbuf.WriteString(funcs[t_hide_cursor])
	}

	cursor_x, cursor_y = x, y
	if !is_cursor_hidden(cursor_x, cursor_y) {
		fmt.Fprintf(&outbuf, funcs[t_move_cursor], cursor_y+1, cursor_y+1)
	}
}

func PutCell(x, y int, cell *Cell) {
	if x < 0 || x >= back_buffer.width {
		return
	}
	if y < 0 || y >= back_buffer.height {
		return
	}

	back_buffer.cells[y*back_buffer.width+x] = *cell
}

func ChangeCell(x, y int, ch rune, fg, bg Attribute) {
	var c = Cell{ch, fg, bg}
	PutCell(x, y, &c)
}

// TODO: func Blit

func PollEvent(event *Event) EventType {
	*event = Event{}

	// try to extract event from input buffer, return on success
	event.Type = EventKey
	if extract_event(event) {
		return EventKey
	}

	for {
		select {
		case data := <-input_comm:
			inbuf = append(inbuf, data...)
			input_comm <- data
			if extract_event(event) {
				return EventKey
			}
		case <-sigwinch_input:
			event.Type = EventResize
			event.Width, event.Height = get_term_size(out.Fd())
			return EventResize
		}
	}
	panic("unreachable")
}

func Size() (int, int) {
	return termw, termh
}

func Clear() {
	select {
	case <-sigwinch_draw:
		update_size()
	default:
	}
	back_buffer.clear()
}

func SetInputMode(mode InputMode) InputMode {
	if mode != InputCurrent {
		input_mode = mode
	}
	return input_mode
}

func SetClearAttributes(fg, bg Attribute) {
	foreground, background = fg, bg
}
