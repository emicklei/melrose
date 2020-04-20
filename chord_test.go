package melrose

import (
	"reflect"
	"strings"
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
		seq     string
		wantErr bool
	}{
		{
			"C major",
			args{"C"},
			Chord{start: N("C"), quality: Major, interval: Triad, inversion: Ground},
			"('[C E G]')",
			false,
		},
		{
			"C diminished 7th",
			args{"C/d7"},
			Chord{start: N("C"), quality: Diminished, interval: Seventh, inversion: Ground},
			"('[C E♭ G♭ A]')",
			false,
		},
		{
			"C augmented",
			args{"C/A"},
			Chord{start: N("C"), quality: Augmented, interval: Triad, inversion: Ground},
			// TODO
			"('C')",
			false,
		},
		{
			"C minor 7",
			args{"C/m7"},
			Chord{start: N("C"), quality: Minor, interval: Seventh, inversion: Ground},
			"('[C E♭ G B♭]')",
			false,
		},
		{
			"C major 7",
			args{"C/7"},
			Chord{start: N("C"), quality: Major, interval: Seventh, inversion: Ground},
			"('[C E G B]')",
			false,
		},
		{
			"C major 6th 2nd inversion",
			args{"C/M6/2"},
			Chord{start: N("C"), quality: Major, interval: Sixth, inversion: Inversion2},
			// TODO
			"('C6')",
			false,
		},
		{
			"C sharp major 1nd inversion",
			args{"C#/1"},
			Chord{start: N("C#"), quality: Major, interval: Triad, inversion: Inversion1},
			"('[F A♭ C♯5]')",
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
			s := strings.Replace(got.S().Storex(), "sequence", "", -1)
			if s != tt.seq {
				t.Errorf("ParseChord(%q) = %s, want %s", tt.args.s, s, tt.seq)
			}
		})
	}
}
