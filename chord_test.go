package melrose

import (
	"reflect"
	"testing"
)

// go test -timeout 30s github.com/emicklei/melrose -v -run "^(TestParseChord)$"
func TestParseChord(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Chord
		wantErr bool
	}{
		{
			"C major",
			args{"C"},
			Chord{start: N("C"), quality: Major, interval: Triad, inversion: Ground},
			false,
		},
		{
			"C augmented",
			args{"C:A"},
			Chord{start: N("C"), quality: Augmented, interval: Triad, inversion: Ground},
			false,
		},
		{
			"C minor 7",
			args{"C:m7"},
			Chord{start: N("C"), quality: Minor, interval: Seventh, inversion: Ground},
			false,
		},
		{
			"C major 6th 2nd inversion",
			args{"C:M6:2"},
			Chord{start: N("C"), quality: Major, interval: Sixth, inversion: Inversion2},
			false,
		},
		{
			"C sharp major 1nd inversion",
			args{"C#:1"},
			Chord{start: N("C#"), quality: Major, interval: Triad, inversion: Inversion1},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseChord(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseChord(%q) error = %v, wantErr %v", tt.args.s, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseChord(%q) = %#v, want %#v", tt.args.s, got, tt.want)
			}
		})
	}
}
