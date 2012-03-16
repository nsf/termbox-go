package termbox

import "syscall"

// public API

// Initializes termbox library. This function should be called before any other functions.
// After successful initialization, the library must be finalized using 'Shutdown' function.
//
// Example usage:
//      err := termbox.Init()
//      if err != nil {
//              panic(err.String())
//      }
//      defer termbox.Shutdown()
func Init() error {
	var err error

	in, err = syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	if err != nil {
		return err
	}
	out, err = syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil {
		return err
	}

	show_cursor(false)

	termw, termh = get_term_size(out)
	back_buffer.init(termw, termh)
	back_buffer.clear()

	attrsbuf = make([]word, termw*termh)
	wcharbuf = make([]wchar, termw*termh)

	err = get_console_mode(in, &orig_mode)
	if err != nil {
		return err
	}

	err = set_console_mode(in, enable_window_input)
	if err != nil {
		return err
	}

	go input_event_producer()

	return nil
}

// Finalizes termbox library, should be called after successful initialization
// when termbox's functionality isn't required anymore.
func Shutdown() {
	set_console_mode(in, orig_mode)
}

// Synchronizes the internal back buffer with the terminal.
func Present() {
	update_size_maybe()
	encode_attrs()
	encode_runes()
	write_console_output_attribute(out, attrsbuf, termw*termh, coord{0, 0}, nil)
	write_console_output_character(out, wcharbuf, termw*termh, coord{0, 0}, nil)
}

// Sets the position of the cursor. See also HideCursor().
func SetCursor(x, y int) {
	if is_cursor_hidden(cursor_x, cursor_y) && !is_cursor_hidden(x, y) {
		show_cursor(true)
	}

	if !is_cursor_hidden(cursor_x, cursor_y) && is_cursor_hidden(x, y) {
		show_cursor(false)
	}

	cursor_x, cursor_y = x, y
	if !is_cursor_hidden(cursor_x, cursor_y) {
		move_cursor(cursor_x, cursor_y)
	}
}

// The shortcut for SetCursor(-1, -1).
func HideCursor() {
	SetCursor(cursor_hidden, cursor_hidden)
}

// Puts the 'cell' into the internal back buffer at the specified position.
func PutCell(x, y int, cell *Cell) {
	if x < 0 || x >= back_buffer.width {
		return
	}
	if y < 0 || y >= back_buffer.height {
		return
	}

	back_buffer.cells[y*back_buffer.width+x] = *cell
}

// Changes cell's parameters in the internal back buffer at the specified
// position.
func ChangeCell(x, y int, ch rune, fg, bg Attribute) {
	var c = Cell{ch, fg, bg}
	PutCell(x, y, &c)
}

// Returns a slice of the termbox back buffer. You can get its dimensions using
// 'Size' function. The slice remains valid as long as no 'Clear' or 'Present'
// function calls were made after call to this function.
//
// The function is provided for performance reasons, normally it is suggested to
// use 'ChangeCell' or 'PutCell'.
func CellBuffer() []Cell {
	return back_buffer.cells
}

// Wait for an event and return it. This is a blocking function call.
func PollEvent() Event {
	return <-input_comm
}

// Returns the size of the internal back buffer (which is the same as
// terminal's window size in characters).
func Size() (int, int) {
	return termw, termh
}

// Clears the internal back buffer.
func Clear() {
	update_size_maybe()
	back_buffer.clear()
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

// Set attributes which are used for clearing the internal back buffer.
func SetClearAttributes(fg, bg Attribute) {
	foreground, background = fg, bg
}
