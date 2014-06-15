package termbox

import "testing"

func TestKeyStringer(t *testing.T) {
	if KeyF1.String() != "F1" {
		t.Error("Key F1 should be named 'F1'")
	}

	k, ok := StringToKey("F1")
	if !ok || k != KeyF1 {
		t.Error("'F1' should result in Key F1")
	}

	if Key(0x1234).String() != UnknownStringIdentifier {
		t.Error("Unknown Key should be named as such")
	}
}

func TestColorStringer(t *testing.T) {
	if ColorBlack.String() != "black" {
		t.Error("Black should be 'black'")
	}

	c, ok := StringToColor("black")
	if !ok || c != ColorBlack {
		t.Error("'black' should be Black")
	}

	if (ColorWhite | AttrBold).String() != "white" {
		t.Error("Bold White should be 'white'")
	}

	if Attribute(0xFFFF).String() != UnknownStringIdentifier {
		t.Error("Unknown Color should be named as such")
	}
}
