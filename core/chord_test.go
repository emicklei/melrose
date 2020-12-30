package core

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
			"C augmented",
			args{"C/aug"},
			Chord{start: N("C"), quality: Augmented, interval: Triad, inversion: Ground},
			"('(C E A♭)')",
			false,
		},
		{
			"A augmented",
			args{"a/+"},
			Chord{start: N("A"), quality: Augmented, interval: Triad, inversion: Ground},
			"('(A D♭5 F5)')",
			false,
		},
		{
			"C augmented seventh",
			args{"C/aug7"},
			Chord{start: N("C"), quality: Augmented, interval: Seventh, inversion: Ground},
			"('(C E A♭ B♭)')",
			false,
		},
		{
			"E♭ augmented seventh",
			args{"E♭/aug7"},
			Chord{start: N("E♭"), quality: Augmented, interval: Seventh, inversion: Ground},
			"('(E♭ G B D♭5)')",
			false,
		},
		{
			"C major",
			args{"C"},
			Chord{start: N("C"), quality: Major, interval: Triad, inversion: Ground},
			"('(C E G)')",
			false,
		},
		{
			"C diminished 7th",
			args{"C/dim7"},
			Chord{start: N("C"), quality: Diminished, interval: Seventh, inversion: Ground},
			"('(C E♭ G♭ A)')",
			false,
		},
		{
			"D diminished",
			args{"d/dim"},
			Chord{start: N("D"), quality: Diminished, interval: Triad, inversion: Ground},
			"('(D F A♭)')",
			false,
		},
		{
			"E diminished",
			args{"e/o"},
			Chord{start: N("E"), quality: Diminished, interval: Triad, inversion: Ground},
			"('(E G B♭)')",
			false,
		},
		{
			"C minor 7",
			args{"C/m7"},
			Chord{start: N("C"), quality: Minor, interval: Seventh, inversion: Ground},
			"('(C E♭ G B♭)')",
			false,
		},
		{
			"C major 7",
			args{"C/M7"},
			Chord{start: N("C"), quality: Major, interval: Seventh, inversion: Ground},
			"('(C E G B)')",
			false,
		},
		{
			"C 7",
			args{"C/7"},
			Chord{start: N("C"), quality: Dominant, interval: Seventh, inversion: Ground},
			"('(C E G B♭)')",
			false,
		},
		{
			"D 7",
			args{"D/7"},
			Chord{start: N("D"), quality: Dominant, interval: Seventh, inversion: Ground},
			"('(D G♭ A C5)')",
			false,
		},
		{
			"E 7",
			args{"E/7"},
			Chord{start: N("E"), quality: Dominant, interval: Seventh, inversion: Ground},
			"('(E A♭ B D5)')",
			false,
		},
		{
			"C major 6th 2nd inversion",
			args{"C/M6/2"},
			Chord{start: N("C"), quality: Major, interval: Sixth, inversion: Inversion2},
			// TODO
			"('C')",
			false,
		},
		{
			"C sharp major 1nd inversion",
			args{"C#/1"},
			Chord{start: N("C#"), quality: Major, interval: Triad, inversion: Inversion1},
			"('(F A♭ C♯5)')",
			false,
		},
		{
			"E minor 2nd inversion",
			args{"E/m/2"},
			Chord{start: N("E"), quality: Minor, interval: Triad, inversion: Inversion2},
			"('(B E5 G5)')",
			false,
		},
		{
			"Rest",
			args{"1="},
			Chord{start: N("1="), quality: Major, interval: Triad, inversion: Ground},
			"('1=')",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseChord(tt.args.s)
			s := strings.Replace(got.S().Storex(), "sequence", "", -1)
			if s != tt.seq {
				t.Errorf("ParseChord(%q) got %s, want %s", tt.args.s, s, tt.seq)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseChord(%q) error = %v, wantErr %v", tt.args.s, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseChord(%q) got %#v, want %#v", tt.args.s, got, tt.want)
			}
		})
	}
}
