package core

import (
	"reflect"
	"testing"
)

func TestMIDI_ToNote(t *testing.T) {
	type fields struct {
		duration Valueable
		number   Valueable
		velocity Valueable
	}
	tests := []struct {
		name   string
		fields fields
		want   Note
	}{
		{
			"F+",
			fields{On(0.25), On(65), On(80)},
			MustParseNote("F+"),
		},
		{
			"½G3+",
			fields{On(0.5), On(55), On(80)},
			MustParseNote("½G3+"),
		},
		{
			"16E♭5",
			fields{On(16), On(75), On(Normal)},
			MustParseNote("16E♭5"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MIDI{
				duration: tt.fields.duration,
				number:   tt.fields.number,
				velocity: tt.fields.velocity,
			}
			if got := m.ToNote(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MIDI.ToNote() = %v, want %v", got, tt.want)
			}
		})
	}
}
