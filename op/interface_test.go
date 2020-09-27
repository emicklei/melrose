package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestImplements(t *testing.T) {
	for _, each := range []struct {
		source          interface{}
		notSequenceable bool
		notStorable     bool
	}{
		{source: Fraction{}},
		{source: Join{}},
		{source: JoinMap{}},
		{source: NoteMap{}},
		{source: Dynamic{}},
	} {
		if !each.notSequenceable {
			if _, ok := each.source.(core.Sequenceable); !ok {
				t.Errorf("%T does not implement Sequenceable", each.source)
			}
		}
		if !each.notStorable {
			if _, ok := each.source.(core.Storable); !ok {
				t.Errorf("%T does not implement Storable", each.source)
			}
		}
	}
}

func storex(s interface{}) string {
	if st, ok := s.(core.Storable); ok {
		return st.Storex()
	}
	return ""
}
