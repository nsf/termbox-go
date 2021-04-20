// +build !windows

package termbox

import "testing"

func TestTerminfoTerms(t *testing.T) {
	for _, term := range terms {
		t.Run(term.name, func(t *testing.T) {
			if len(term.funcs) != t_max_funcs {
				t.Errorf("want %d got %d terminfo entries", t_max_funcs, len(term.funcs))
			}
		})
	}
}
