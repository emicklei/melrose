package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestVolume_S(t *testing.T) {
	value := &core.ValueHolder{Any: 42}
	volume := Volume{
		Target: []core.Sequenceable{
			core.MustParseSequence("(C E) D"),
			core.MustParseNote("F"),
		},
		Value: value,
	}

	sequence := volume.S()
	if got, want := len(sequence.Notes), 3; got != want {
		t.Fatalf("got %d note groups, want %d", got, want)
	}
	if got, want := len(sequence.Notes[0]), 2; got != want {
		t.Fatalf("got %d notes in first group, want %d", got, want)
	}
	for _, group := range sequence.Notes {
		for _, note := range group {
			if got, want := note.Velocity, 42; got != want {
				t.Errorf("got velocity %d, want %d", got, want)
			}
		}
	}

	value.Any = 127
	if got, want := volume.S().Notes[0][0].Velocity, 127; got != want {
		t.Errorf("got updated velocity %d, want %d", got, want)
	}
}

func TestVolume_InvalidValue(t *testing.T) {
	sequence := core.MustParseSequence("C")
	for _, value := range []any{-1, 128, "loud"} {
		volume := Volume{Target: []core.Sequenceable{sequence}, Value: core.On(value)}
		if got := volume.S(); len(got.Notes) != 0 {
			t.Errorf("value %v produced %d note groups, want none", value, len(got.Notes))
		}
	}
}

func TestVolume_StorexAndReplaced(t *testing.T) {
	from := core.MustParseSequence("C D")
	to := core.MustParseSequence("E F")
	volume := Volume{Target: []core.Sequenceable{from}, Value: core.On(64)}

	if got, want := volume.Storex(), "volume(64,sequence('C D'))"; got != want {
		t.Errorf("got [%s], want [%s]", got, want)
	}
	if !core.IsIdenticalTo(volume.Replaced(from, to).(Volume).Target[0], to) {
		t.Error("target was not replaced")
	}
	if !core.IsIdenticalTo(volume.Replaced(volume, to), to) {
		t.Error("volume was not replaced")
	}
}

func TestVolume_EmptyTarget(t *testing.T) {
	if got := (Volume{Value: core.On(64)}).S(); len(got.Notes) != 0 {
		t.Errorf("got %d note groups, want none", len(got.Notes))
	}
}
