package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestNewDynamicMapper(t *testing.T) {
	l := core.MustParseSequence("A B")
	d := NewDynamicMapper([]core.Sequenceable{l}, "1:++,2:--")
	if got, want := d.Storex(), "dynamic('1:++,2:--',sequence('A B'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := d.S().Storex(), "sequence('A++ B--')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}

}
