package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestNewDynamicMapper(t *testing.T) {
	l := core.MustParseSequence("A B")
	d := NewDynamicMap([]core.Sequenceable{l}, core.On("1:++,2:--"))
	if len(d.S().Notes) == 0 {
		t.Fail()
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
	d := NewDynamicMap([]core.Sequenceable{l}, core.On("2:o,1:++,2:--,1:++"))
	if got, want := d.S().Storex(), "sequence('B A++ B-- A++')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNewDynamicMapper_InvalidVelocity(t *testing.T) {
	l := core.MustParseSequence("A B")
	r := NewDynamicMap([]core.Sequenceable{l}, core.On("1:~")).S()
	if len(r.Notes) != 0 {
		t.Fail()
	}
}

func TestNewDynamicMapper_InvalidIndex(t *testing.T) {
	l := core.MustParseSequence("A B")
	r := NewDynamicMap([]core.Sequenceable{l}, core.On("-1:+++")).S()
	if len(r.Notes) != 0 {
		t.Fail()
	}
}

func TestDynamicMap_Replaced(t *testing.T) {
	l := core.MustParseSequence("A B")
	d := NewDynamicMap([]core.Sequenceable{l}, core.On("1:++,2:--"))
	r := d.S()
	if len(r.Notes) == 0 {
		t.Fail()
	}
	if core.IsIdenticalTo(d, l) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(d.Replaced(l, core.EmptySequence).(DynamicMap).Target[0], core.EmptySequence) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(d.Replaced(d, l), l) {
		t.Error("should be replaced by l")
	}
}

func TestDynamicMap_Invalid(t *testing.T) {
	l := core.MustParseSequence("A B")
	d := NewDynamicMap([]core.Sequenceable{l}, core.On("1:++,3:--"))
	if got, want := d.S().Storex(), "sequence('A++')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	NewDynamicMap([]core.Sequenceable{l}, core.On("a:b"))
	NewDynamicMap([]core.Sequenceable{l}, core.On("1:b"))
	NewDynamicMap([]core.Sequenceable{l}, core.On("a:++"))
}
