package core

import "testing"

func TestImplements(t *testing.T) {
	for _, each := range []struct {
		source          interface{}
		notSequenceable bool
		notStorable     bool
	}{
		{source: Note{}},
		{source: Chord{}},
		{source: Scale{}},
		{source: Progression{}},
		{source: new(Track)},
		{source: new(MultiTrack), notSequenceable: true},
	} {
		if !each.notSequenceable {
			if _, ok := each.source.(Sequenceable); !ok {
				t.Errorf("%T does not implement Sequenceable", each.source)
			}
		}
		if !each.notStorable {
			if _, ok := each.source.(Storable); !ok {
				t.Errorf("%T does not implement Sequenceable", each.source)
			}
		}
	}
}
