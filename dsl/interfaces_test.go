package dsl

import (
	"testing"

	"github.com/emicklei/melrose/control"
	"github.com/emicklei/melrose/core"
)

func TestImplements(t *testing.T) {
	for _, each := range []struct {
		source          interface{}
		notSequenceable bool
		notStorable     bool
	}{
		{source: core.Note{}},
		{source: core.Chord{}},
		{source: core.Scale{}},
		{source: core.ChordSequence{}},
		{source: new(core.Track)},
		{source: new(core.MultiTrack), notSequenceable: true},
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

func TestImplementsPlayable(t *testing.T) {
	for _, each := range []struct {
		source          interface{}
		notSequenceable bool
		notStorable     bool
	}{
		{source: new(core.Loop)},
		{source: new(control.Listen)},
	} {
		if !each.notSequenceable {
			if _, ok := each.source.(core.Playable); !ok {
				t.Errorf("%T does not implement Playable", each.source)
			}
		}
	}
}
