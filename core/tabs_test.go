package core

import (
	"reflect"
	"testing"
)

func TestParseTabNote(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name     string
		args     args
		want     TabNote
		wantMIDI int
		wantErr  bool
	}{
		{"open e", args{"e"}, TabNote{Name: "E", Velocity: Normal, fraction: 0.25}, 40, false},
		{"e3", args{"e3"}, TabNote{Name: "E", Fret: 3, Velocity: Normal, fraction: 0.25}, 43, false},
		{"a24", args{"a24"}, TabNote{Name: "A", Fret: 24, Velocity: Normal, fraction: 0.25}, 69, false},
		// errors
		{"c2", args{"c2"}, TabNote{Name: "", Fret: 0}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTabNote(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTabNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTabNote() = %v, want %v", got, tt.want)
			}
			if m := got.ToNote().MIDI(); m != tt.wantMIDI {
				t.Errorf("ToNote().MIDI() = %v, want %v", m, tt.wantMIDI)
			}
		})
	}
}
