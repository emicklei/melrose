package core

import (
	"strings"
	"testing"
)

func TestParseChord_C(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		seq     string
		wantErr bool
	}{
		// Triads
		{
			"C major",
			args{"C"},
			"('(C E G)')",
			false,
		},
		{
			"C minor",
			args{"C/m"},
			"('(C E♭ G)')",
			false,
		},
		{
			"C augmented",
			args{"C/aug"},
			"('(C E A♭)')",
			false,
		},
		{
			"C +",
			args{"C/+"},
			"('(C E A♭)')",
			false,
		},
		{
			"C diminished",
			args{"C/dim"},
			"('(C E♭ G♭)')",
			false,
		},
		{
			"C sus2",
			args{"C/sus2"},
			"('(C D G)')",
			false,
		},
		{
			"C sus4",
			args{"C/sus4"},
			"('(C F G)')",
			false,
		},
		// Seventh
		{
			"C major 7",
			args{"C/maj7"},
			"('(C E G B)')",
			false,
		},
		{
			"C dominant 7",
			args{"C/7"},
			"('(C E G B♭)')",
			false,
		},
		// {
		// 	"C minor major 7",
		// 	args{"C/mmaj7"},
		// 	"('(C E♭ G B♭)')",
		// 	false, // TODO
		// },
		{
			"C minor 7",
			args{"C/m7"},
			"('(C E♭ G B♭)')",
			false,
		},
		{
			"C augmented seventh",
			args{"C/aug7"},
			"('(C E A♭ B♭)')",
			false,
		},
		// Fifth
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseChord(tt.args.s)
			s := strings.Replace(got.S().Storex(), "sequence", "", -1)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseChord(%q) error = %v, wantErr %v", tt.args.s, err, tt.wantErr)
				return
			}
			if s != tt.seq {
				t.Errorf("ParseChord(%q) got %s, want %s", tt.args.s, s, tt.seq)
			}
		})
	}
}

// https://muzieknotatie.nl/wiki/Akkoordsymbool
// go test -timeout 30s github.com/emicklei/melrose -v -run "^(TestParseChord)$"
func TestParseChord(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		seq     string
		wantErr bool
	}{

		{
			"A augmented",
			args{"a/+"},
			"('(A D♭5 F5)')",
			false,
		},
		{
			"E♭ augmented seventh",
			args{"E♭/aug7"},
			"('(E♭ G B D♭5)')",
			false,
		},

		{
			"D diminished",
			args{"d/dim"},
			"('(D F A♭)')",
			false,
		},
		{
			"E diminished",
			args{"e/o"},
			"('(E G B♭)')",
			false,
		},
		{
			"D 7",
			args{"D/7"},
			"('(D G♭ A C5)')",
			false,
		},

		{
			"G 7",
			args{"G/7"},
			"('(G B D5 F5)')",
			false,
		},
		{
			"E 7",
			args{"E/7"},
			"('(E A♭ B D5)')",
			false,
		},
		{
			"C major 6th 2nd inversion",
			args{"C/maj6/2"},
			// TODO
			"('C')",
			false,
		},
		{
			"C sharp major 1nd inversion",
			args{"C#/1"},
			"('(F A♭ C♯5)')",
			false,
		},
		{
			"E minor 2nd inversion",
			args{"E/m/2"},
			"('(B E5 G5)')",
			false,
		},
		{
			"Rest",
			args{"1="},
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
		})
	}
}
