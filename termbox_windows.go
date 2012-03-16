package termbox

import "syscall"
import "unsafe"
import "unicode/utf16"

type (
	wchar uint16
	short int16
	dword uint32
	word  uint16
	coord struct {
		x short
		y short
	}
	small_rect struct {
		left   short
		top    short
		right  short
		bottom short
	}
	console_screen_buffer_info struct {
		size                coord
		cursor_position     coord
		attributes          word
		window              small_rect
		maximum_window_size coord
	}
	console_cursor_info struct {
		size    dword
		visible int32
	}
	input_record struct {
		event_type word
		_          [2]byte
		event      [16]byte
	}
	key_event_record struct {
		key_down          int32
		repeat_count      word
		virtual_key_code  word
		virtual_scan_code word
		unicode_char      wchar
		control_key_state dword
	}
	window_buffer_size_record struct {
		size coord
	}
)

func (this coord) uintptr() uintptr {
	return uintptr(*(*int32)(unsafe.Pointer(&this)))
}

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	proc_get_console_screen_buffer_info = kernel32.NewProc("GetConsoleScreenBufferInfo")
	proc_write_console_output_character = kernel32.NewProc("WriteConsoleOutputCharacterW")
	proc_write_console_output_attribute = kernel32.NewProc("WriteConsoleOutputAttribute")
	proc_set_console_cursor_info        = kernel32.NewProc("SetConsoleCursorInfo")
	proc_set_console_cursor_position    = kernel32.NewProc("SetConsoleCursorPosition")
	proc_read_console_input             = kernel32.NewProc("ReadConsoleInputW")
	proc_get_console_mode               = kernel32.NewProc("GetConsoleMode")
	proc_set_console_mode               = kernel32.NewProc("SetConsoleMode")
)

func get_console_screen_buffer_info(h syscall.Handle, info *console_screen_buffer_info) (err error) {
	r0, _, e1 := syscall.Syscall(proc_get_console_screen_buffer_info.Addr(),
		2, uintptr(h), uintptr(unsafe.Pointer(info)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func write_console_output_character(h syscall.Handle, chars []wchar, n int, pos coord, written *dword) (err error) {
	r0, _, e1 := syscall.Syscall6(proc_write_console_output_character.Addr(),
		5, uintptr(h), uintptr(unsafe.Pointer(&chars[0])), uintptr(n), pos.uintptr(),
		uintptr(unsafe.Pointer(written)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func write_console_output_attribute(h syscall.Handle, attrs []word, n int, pos coord, written *dword) (err error) {
	r0, _, e1 := syscall.Syscall6(proc_write_console_output_attribute.Addr(),
		5, uintptr(h), uintptr(unsafe.Pointer(&attrs[0])), uintptr(n), pos.uintptr(),
		uintptr(unsafe.Pointer(written)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func set_console_cursor_info(h syscall.Handle, info *console_cursor_info) (err error) {
	r0, _, e1 := syscall.Syscall(proc_set_console_cursor_info.Addr(),
		2, uintptr(h), uintptr(unsafe.Pointer(info)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func set_console_cursor_position(h syscall.Handle, pos coord) (err error) {
	r0, _, e1 := syscall.Syscall(proc_set_console_cursor_position.Addr(),
		2, uintptr(h), pos.uintptr(), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func read_console_input(h syscall.Handle, record *input_record) (err error) {
	var read dword // required, it fails without it
	r0, _, e1 := syscall.Syscall6(proc_read_console_input.Addr(),
		4, uintptr(h), uintptr(unsafe.Pointer(record)), 1, uintptr(unsafe.Pointer(&read)), 0, 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func get_console_mode(h syscall.Handle, mode *dword) (err error) {
	r0, _, e1 := syscall.Syscall(proc_get_console_mode.Addr(),
		2, uintptr(h), uintptr(unsafe.Pointer(mode)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func set_console_mode(h syscall.Handle, mode dword) (err error) {
	r0, _, e1 := syscall.Syscall(proc_set_console_mode.Addr(),
		2, uintptr(h), uintptr(mode), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

var (
	orig_mode   dword
	back_buffer cellbuf
	termw       int
	termh       int
	input_mode  = InputEsc
	cursor_x    = cursor_hidden
	cursor_y    = cursor_hidden
	foreground  = ColorDefault
	background  = ColorDefault
	in          syscall.Handle
	out         syscall.Handle
	attrsbuf    []word
	wcharbuf    []wchar
	input_comm  = make(chan Event)
)

func get_term_size(out syscall.Handle) (int, int) {
	var info console_screen_buffer_info
	err := get_console_screen_buffer_info(out, &info)
	if err != nil {
		panic(err)
	}
	return int(info.size.x), int(info.size.y)
}

func update_size_maybe() {
	w, h := get_term_size(out)
	if w != termw || h != termh {
		termw, termh = w, h
		back_buffer.resize(termw, termh)

		size := termw * termh
		if cap(attrsbuf) < size {
			attrsbuf = make([]word, size)
		}
		if cap(wcharbuf) < size {
			wcharbuf = make([]wchar, size)
		}
	}
}

var color_table_bg = []word{
	0, // black
	background_red,
	background_green,
	background_red | background_green, // yellow
	background_blue,
	background_red | background_blue,                    // magenta
	background_green | background_blue,                  // cyan
	background_red | background_blue | background_green, // white
	0,                                                   // default (black)
}

var color_table_fg = []word{
	0,
	foreground_red,
	foreground_green,
	foreground_red | foreground_green, // yellow
	foreground_blue,
	foreground_red | foreground_blue,                    // magenta
	foreground_green | foreground_blue,                  // cyan
	foreground_red | foreground_blue | foreground_green, // white
	foreground_red | foreground_blue | foreground_green, // default (white)
}

// encodes all attributes in the back buffer to winapi format and places them to
// attrsbuf
func encode_attrs() {
	n := len(back_buffer.cells)
	if cap(attrsbuf) < n {
		attrsbuf = make([]word, n)
	} else {
		attrsbuf = attrsbuf[:n]
	}

	for i, cell := range back_buffer.cells {
		attr := color_table_fg[cell.Fg&0x0F] |
			color_table_bg[cell.Bg&0x0F]
		if cell.Fg&AttrBold != 0 {
			attr |= foreground_intensity
		}
		if cell.Bg&AttrBold != 0 {
			attr |= background_intensity
		}

		attrsbuf[i] = attr
	}
}

const (
	replacement_char = '\uFFFD'
	max_rune         = '\U0010FFFF'
	surr1            = 0xd800
	surr2            = 0xdc00
	surr3            = 0xe000
	surr_self        = 0x10000
)

// encodes all runes in the back buffer to utf16 and places them to wcharbuf
func encode_runes() {
	n := len(back_buffer.cells)
	for _, cell := range back_buffer.cells {
		if cell.Ch >= surr_self {
			n++
		}
	}

	if cap(wcharbuf) < n {
		wcharbuf = make([]wchar, n)
	} else {
		wcharbuf = wcharbuf[:n]
	}

	n = 0
	for _, cell := range back_buffer.cells {
		v := cell.Ch
		switch {
		case v < 0, surr1 <= v && v < surr3, v > max_rune:
			v = replacement_char
			fallthrough
		case v < surr_self:
			wcharbuf[n] = wchar(v)
			n++
		default:
			r1, r2 := utf16.EncodeRune(v)
			wcharbuf[n] = wchar(r1)
			wcharbuf[n+1] = wchar(r2)
			n += 2
		}
	}
}

func move_cursor(x, y int) {
	err := set_console_cursor_position(out, coord{short(x), short(y)})
	if err != nil {
		panic(err)
	}
}

func show_cursor(visible bool) {
	var v int32
	if visible {
		v = 1
	}

	var info console_cursor_info
	info.size = 100
	info.visible = v
	err := set_console_cursor_info(out, &info)
	if err != nil {
		panic(err)
	}
}

func key_event_record_to_event(r *key_event_record) (Event, bool) {
	if r.key_down == 0 {
		return Event{}, false
	}

	e := Event{Type: EventKey}
	if input_mode == InputAlt {
		if r.control_key_state&(left_alt_pressed|right_alt_pressed) != 0 {
			e.Mod = ModAlt
		}
	}

	if r.virtual_key_code >= vk_f1 && r.virtual_key_code <= vk_f12 {
		switch r.virtual_key_code {
		case vk_f1:
			e.Key = KeyF1
		case vk_f2:
			e.Key = KeyF2
		case vk_f3:
			e.Key = KeyF3
		case vk_f4:
			e.Key = KeyF4
		case vk_f5:
			e.Key = KeyF5
		case vk_f6:
			e.Key = KeyF6
		case vk_f7:
			e.Key = KeyF7
		case vk_f8:
			e.Key = KeyF8
		case vk_f9:
			e.Key = KeyF9
		case vk_f10:
			e.Key = KeyF10
		case vk_f11:
			e.Key = KeyF11
		case vk_f12:
			e.Key = KeyF12
		default:
			panic("unreachable")
		}

		return e, true
	}

	if r.virtual_key_code <= vk_delete {
		switch r.virtual_key_code {
		case vk_insert:
			e.Key = KeyInsert
		case vk_delete:
			e.Key = KeyDelete
		case vk_home:
			e.Key = KeyHome
		case vk_end:
			e.Key = KeyEnd
		case vk_pgup:
			e.Key = KeyPgup
		case vk_pgdn:
			e.Key = KeyPgdn
		case vk_arrow_up:
			e.Key = KeyArrowUp
		case vk_arrow_down:
			e.Key = KeyArrowDown
		case vk_arrow_left:
			e.Key = KeyArrowLeft
		case vk_arrow_right:
			e.Key = KeyArrowRight
		case vk_backspace:
			e.Key = KeyBackspace
		case vk_tab:
			e.Key = KeyTab
		case vk_enter:
			e.Key = KeyEnter
		case vk_esc:
			e.Key = KeyEsc
		case vk_space:
			e.Key = KeySpace
		default:
			goto keep_matching
		}

		return e, true
	}

keep_matching:
	if r.control_key_state&(left_ctrl_pressed|right_ctrl_pressed) != 0 {
		if Key(r.unicode_char) >= KeyCtrlA && Key(r.unicode_char) <= KeyCtrlZ {
			e.Key = Key(r.unicode_char)
			return e, true
		}
	}

	if r.unicode_char != 0 {
		e.Ch = rune(r.unicode_char)
		return e, true
	}

	return Event{}, false
}

func input_event_producer() {
	var r input_record
	var err error
	for {
		err = read_console_input(in, &r)
		if err != nil {
			panic(err)
		}

		switch r.event_type {
		case key_event:
			kr := (*key_event_record)(unsafe.Pointer(&r.event))
			ev, ok := key_event_record_to_event(kr)
			if ok {
				for i := 0; i < int(kr.repeat_count); i++ {
					input_comm <- ev
				}
			}
		case window_buffer_size_event:
			sr := *(*window_buffer_size_record)(unsafe.Pointer(&r.event))
			input_comm <- Event{
				Type:   EventResize,
				Width:  int(sr.size.x),
				Height: int(sr.size.y),
			}
		}
	}
}
