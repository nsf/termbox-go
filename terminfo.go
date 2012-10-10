// +build !windows
// This file contains a simple and incomplete implementation of the terminfo
// database. Information was taken from the ncurses manpages term(5) and
// terminfo(5). Currently, only the string capabilities for special keys and for
// functions without parameters are actually used. Colors are still done with
// ANSI escape sequences. Other special features that are not (yet?) supported
// are reading from ~/.terminfo, the TERMINFO_DIRS variable, Berkeley database
// format and extended capabilities.

package termbox

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"fmt"
	"os"
)

const (
	tiMagic = 0432
	tiHeaderLength = 12
)

func setup_term() (err error) {
	var data []byte
	var header []int16
	var strOffset, tableOffset int16

	term := os.Getenv("TERM")
	if term == "" {
		err = fmt.Errorf("termbox: TERM not set")
		return
	}

	// TODO: look in ~/.terminfo
	path := os.Getenv("TERMINFO")
	if path == "" {
		path = "/usr/share/terminfo"
	}
	path += "/" + term[0:1] + "/" + term

	data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	rd := bytes.NewReader(data)
	header = make([]int16, 6)
	// 0: magic number, 1: size of names section, 2: size of boolean section, 3:
	// size of numbers section (in integers), 4: size of the strings section (in
	// integers), 5: size of the string table

	err = binary.Read(rd, binary.LittleEndian, header)
	if err != nil {
		return
	}

	if header[2] % 2 != 0 {
		// old quirk to align everything on word boundaries
		header[2] += 1
	}
	strOffset = tiHeaderLength + header[1] + header[2] + 2 * header[3]
	tableOffset = strOffset + 2 * header[4]

	keys = make([]string, 0xFFFF - keyMax)
	for i, _ := range keys {
		keys[i], err = tiReadString(rd, strOffset + 2 * tiKeys[i], tableOffset)
		if err != nil {
			return
		}
	}
	funcs = make([]string, t_max_funcs)
	for i, _ := range funcs {
		funcs[i], err = tiReadString(rd, strOffset + 2 * tiFuncs[i], tableOffset)
		if err != nil {
			return
		}
	}
	err = nil
	return
}

func tiReadString(rd *bytes.Reader, strOff, table int16) (string, error) {
	var off int16

	_, err := rd.Seek(int64(strOff), 0)
	if err != nil {
		return "", err
	}
	err = binary.Read(rd, binary.LittleEndian, &off)
	if err != nil {
		return "", err
	}
	_, err = rd.Seek(int64(table + off), 0)
	if err != nil {
		return "", err
	}
	var bs []byte
	for {
		b, err := rd.ReadByte()
		if err != nil {
			return "", err
		}
		if b == byte(0x00) {
			break
		}
		bs = append(bs, b)
	}
	return string(bs), nil
}

// "Maps" the function constants from termbox.go to the number of the respective
// string capability in the terminfo file. Taken from (ncurses) term.h.
var tiFuncs = []int16{
	28, 40, 16, 13, 5, 39, 36, 27, 26, 34, 89, 88 }

// Same as above for the special keys.
var tiKeys = []int16{
	66, 68 /* apparently not a typo; 67 is F10 for whatever reason */, 69, 70,
	71, 72, 73, 74, 75, 67, 216, 217, 77, 59, 76, 164, 82, 81, 87, 61, 79, 83 }
