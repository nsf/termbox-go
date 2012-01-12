// Termbox library provides facilities for terminal input/output manipulation
// in a pseudo-GUI style.
package termbox

// #include "termbox.h"
import "C"
import (
	"unsafe"
	"errors"
)

// This type represents termbox event. 'Mod', 'Key' and 'Ch' fields are valid
// if 'Type' is EVENT_KEY. 'W' and 'H' are valid if 'Type' is EVENT_RESIZE.
type Event struct {
	Type uint8  // one of EVENT_ constants
	Mod  uint8  // one of MOD_ constants or 0
	Key  uint16 // one of KEY_ constants, invalid if 'Ch' is not 0
	Ch   rune   // a unicode character
	W    int32  // width of the screen
	H    int32  // height of the screen
}

// A cell, single conceptual entity on the screen. The screen is basically a 2d
// array of cells. 'Ch' is a unicode character, 'Fg' and 'Bg' are foreground
// and background attributes respectively.
type Cell struct {
	Ch rune
	Fg uint16
	Bg uint16
}

type struct_tb_event_ptr *C.struct_tb_event
type struct_tb_cell_ptr *C.struct_tb_cell

// Key constants, see Event.Key field.
const (
	KEY_F1          = (0xFFFF - 0)
	KEY_F2          = (0xFFFF - 1)
	KEY_F3          = (0xFFFF - 2)
	KEY_F4          = (0xFFFF - 3)
	KEY_F5          = (0xFFFF - 4)
	KEY_F6          = (0xFFFF - 5)
	KEY_F7          = (0xFFFF - 6)
	KEY_F8          = (0xFFFF - 7)
	KEY_F9          = (0xFFFF - 8)
	KEY_F10         = (0xFFFF - 9)
	KEY_F11         = (0xFFFF - 10)
	KEY_F12         = (0xFFFF - 11)
	KEY_INSERT      = (0xFFFF - 12)
	KEY_DELETE      = (0xFFFF - 13)
	KEY_HOME        = (0xFFFF - 14)
	KEY_END         = (0xFFFF - 15)
	KEY_PGUP        = (0xFFFF - 16)
	KEY_PGDN        = (0xFFFF - 17)
	KEY_ARROW_UP    = (0xFFFF - 18)
	KEY_ARROW_DOWN  = (0xFFFF - 19)
	KEY_ARROW_LEFT  = (0xFFFF - 20)
	KEY_ARROW_RIGHT = (0xFFFF - 21)

	KEY_CTRL_TILDE       = 0x00
	KEY_CTRL_2           = 0x00
	KEY_CTRL_A           = 0x01
	KEY_CTRL_B           = 0x02
	KEY_CTRL_C           = 0x03
	KEY_CTRL_D           = 0x04
	KEY_CTRL_E           = 0x05
	KEY_CTRL_F           = 0x06
	KEY_CTRL_G           = 0x07
	KEY_BACKSPACE        = 0x08
	KEY_CTRL_H           = 0x08
	KEY_TAB              = 0x09
	KEY_CTRL_I           = 0x09
	KEY_CTRL_J           = 0x0A
	KEY_CTRL_K           = 0x0B
	KEY_CTRL_L           = 0x0C
	KEY_ENTER            = 0x0D
	KEY_CTRL_M           = 0x0D
	KEY_CTRL_N           = 0x0E
	KEY_CTRL_O           = 0x0F
	KEY_CTRL_P           = 0x10
	KEY_CTRL_Q           = 0x11
	KEY_CTRL_R           = 0x12
	KEY_CTRL_S           = 0x13
	KEY_CTRL_T           = 0x14
	KEY_CTRL_U           = 0x15
	KEY_CTRL_V           = 0x16
	KEY_CTRL_W           = 0x17
	KEY_CTRL_X           = 0x18
	KEY_CTRL_Y           = 0x19
	KEY_CTRL_Z           = 0x1A
	KEY_ESC              = 0x1B
	KEY_CTRL_LSQ_BRACKET = 0x1B
	KEY_CTRL_3           = 0x1B
	KEY_CTRL_4           = 0x1C
	KEY_CTRL_BACKSLASH   = 0x1C
	KEY_CTRL_5           = 0x1D
	KEY_CTRL_RSQ_BRACKET = 0x1D
	KEY_CTRL_6           = 0x1E
	KEY_CTRL_7           = 0x1F
	KEY_CTRL_SLASH       = 0x1F
	KEY_CTRL_UNDERSCORE  = 0x1F
	KEY_SPACE            = 0x20
	KEY_BACKSPACE2       = 0x7F
	KEY_CTRL_8           = 0x7F
)

// Alt modifier constant, see Event.Mod field and SetInputMode function.
const MOD_ALT = 0x01

// Cell attributes, it is possible to use multiple attributes by combining them
// using bitwise OR ('|'). Although, colors cannot be combined. But you can
// combine attributes and a single color.
const (
	BLACK   = 0x00
	RED     = 0x01
	GREEN   = 0x02
	YELLOW  = 0x03
	BLUE    = 0x04
	MAGENTA = 0x05
	CYAN    = 0x06
	WHITE   = 0x07

	BOLD      = 0x10
	UNDERLINE = 0x20
)

// Special coordinate for SetCursor. If you call:
//	SetCursor(HIDE_CURSOR, HIDE_CURSOR)
// This function call hides the cursor.
const HIDE_CURSOR = -1

// Input mode. See SelectInputMode function.
const (
	INPUT_ESC = 1
	INPUT_ALT = 2
)

// Event type. See Event.Type field.
const (
	EVENT_KEY    = 1
	EVENT_RESIZE = 2
)

// Initializes termbox library. This function should be called before any other functions.
// After successful initialization, the library must be finalized using 'Shutdown' function.
//
// Example usage:
//	err := termbox.Init()
//	if err != nil {
//		panic(err.String())
//	}
//	defer termbox.Shutdown()
func Init() error {
	switch int(C.tb_init()) {
	case -3:
		return errors.New("Pipe trap error")
	case -2:
		return errors.New("Failed to open /dev/tty")
	case -1:
		return errors.New("Unsupported terminal")
	}
	return nil
}

// Finalizes termbox library, should be called after successful initialization
// when termbox's functionality isn't required anymore.
func Shutdown() {
	C.tb_shutdown()
}

// Changes cell's parameters in the internal back buffer at the specified
// position.
func ChangeCell(x int, y int, ch rune, fg uint16, bg uint16) {
	C.tb_change_cell(C.uint(x), C.uint(y), C.uint32_t(ch), C.uint16_t(fg), C.uint16_t(bg))
}

// Puts the 'cell' into the internal back buffer at the specified position.
func PutCell(x, y int, cell *Cell) {
	C.tb_put_cell(C.uint(x), C.uint(y), struct_tb_cell_ptr(unsafe.Pointer(cell)))
}

// 'Blit' function copies the 'cells' buffer to the internal back buffer at the
// position specified by 'x' and 'y'. Blit doesn't perform any kind of cuts and
// if contents of the cells buffer cannot be placed without crossing back
// buffer's boundaries, the operation is discarded. Parameter 'w' must be > 0,
// otherwise it will cause "division by zero" panic.
//
// The width and the height of the 'cells' buffer are calculated that way:
//	w := w
//	h := len(cells) / w
func Blit(x, y, w int, cells []Cell) {
	h := len(cells) / w
	C.tb_blit(C.uint(x), C.uint(y), C.uint(w), C.uint(h), struct_tb_cell_ptr(unsafe.Pointer(&cells[0])))
}

// Synchronizes the internal back buffer with the terminal.
func Present() {
	C.tb_present()
}

// Clears the internal back buffer.
func Clear() {
	C.tb_clear()
}

// Wait for an event. This is a blocking function call. If an error occurs,
// returns -1. Otherwise the return value is one of EVENT_ consts.
func PollEvent(e *Event) int {
	return int(C.tb_poll_event(struct_tb_event_ptr(unsafe.Pointer(e))))
}

// Wait for an event 'timeout' milliseconds. If no event occurs, returns 0. If
// an error occurs, returns -1. Otherwise the return value is one of EVENT_
// consts.
func PeekEvent(e *Event, timeout int) int {
	return int(C.tb_peek_event(struct_tb_event_ptr(unsafe.Pointer(e)), C.uint(timeout)))
}

// Returns the width of the internal back buffer (which is the same as
// terminal's window width in characters).
func Width() int {
	return int(C.tb_width())
}

// Returns the height of the internal back buffer (which is the same as
// terminal's window height in characters).
func Height() int {
	return int(C.tb_height())
}

// Sets the position of the cursor. See also HIDE_CURSOR and HideCursor().
func SetCursor(x int, y int) {
	C.tb_set_cursor(C.int(x), C.int(y))
}

// The shortcut for SetCursor(HIDE_CURSOR, HIDE_CURSOR).
func HideCursor() {
	C.tb_set_cursor(HIDE_CURSOR, HIDE_CURSOR)
}

// Selects termbox input mode. Termbox has two input modes:
//
// 1. ESC input mode. When ESC sequence is in the buffer and it doesn't
// match any known sequence. ESC means KEY_ESC.
//
// 2. ALT input mode. When ESC sequence is in the buffer and it doesn't match
// any known sequence. ESC enables MOD_ALT modifier for the next keyboard
// event.
//
// If 'mode' is 0, returns the current input mode. See also INPUT_ constants.
//
// Note: INPUT_ALT mode may not work with PeekEvent.
func SelectInputMode(mode int) {
	C.tb_select_input_mode(C.int(mode))
}

// Shortcut for termbox.PollEvent(e).
func (e *Event) Poll() int {
	return PollEvent(e)
}

// Shortcut for termbox.PeekEvent(e, timeout).
func (e *Event) Peek(timeout int) int {
	return PeekEvent(e, timeout)
}
