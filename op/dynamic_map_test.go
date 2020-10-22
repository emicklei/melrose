package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestNewDynamicMapper(t *testing.T) {
	l := core.MustParseSequence("A B")
	d, err := NewDynamicMap([]core.Sequenceable{l}, "1:++,2:--")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := d.Storex(), "dynamicmap('1:++,2:--',sequence('A B'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := d.S().Storex(), "sequence('A++ B--')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}

}

func TestNewDynamicMapper_DuplicateAndChangeOrder(t *testing.T) {
	l := core.MustParseSequence("A B")
	d, err := NewDynamicMap([]core.Sequenceable{l}, "2:o,1:++,2:--,1:++")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := d.S().Storex(), "sequence('B A++ B-- A++')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNewDynamicMapper_InvalidVelocity(t *testing.T) {
	l := core.MustParseSequence("A B")
	_, err := NewDynamicMap([]core.Sequenceable{l}, "1:~")
	if err == nil {
		t.Fail()
	}
	t.Log(err)
}

func TestNewDynamicMapper_InvalidIndex(t *testing.T) {
	l := core.MustParseSequence("A B")
	_, err := NewDynamicMap([]core.Sequenceable{l}, "-1:+++")
	if err == nil {
		t.Fail()
	}
	t.Log(err)
}
