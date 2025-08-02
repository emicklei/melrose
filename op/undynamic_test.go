package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestUndynamic_S(t *testing.T) {
	s := core.MustParseSequence("C++ D-- E")
	u := Undynamic{Target: s}
	if got, want := u.S().Storex(), "sequence('C D E')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestUndynamic_Storex(t *testing.T) {
	s := core.MustParseSequence("C")
	u := Undynamic{Target: s}
	if got, want := u.Storex(), "undynamic(sequence('C'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	// TODO
	// u = Undynamic{Target: failingNoteConvertable{}}
	// if got, want := u.Storex(), ""; got != want {
	// 	t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	// }
}
