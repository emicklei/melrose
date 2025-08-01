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
			"('(C E_ G)')",
			false,
		},
		{
			"C augmented",
			args{"C/aug"},
			"('(C E A_)')",
			false,
		},
		{
			"C +",
			args{"C/+"},
			"('(C E A_)')",
			false,
		},
		{
			"C diminished",
			args{"C/dim"},
			"('(C E_ G_)')",
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
			"('(C E G B_)')",
			false,
		},
		// {
		// 	"C minor major 7",
		// 	args{"C/mmaj7"},
		// 	"('(C E_ G B_)')",
		// 	false, // TODO
		// },
		{
			"C minor 7",
			args{"C/m7"},
			"('(C E_ G B_)')",
			false,
		},
		{
			"C augmented seventh",
			args{"C/aug7"},
			"('(C E A_ B_)')",
			false,
		},
		{
			"G Major second inv",
			args{"g/M/2"},
			"('(D5 G5 B5)')",
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
func TestChordWithInversion(t *testing.T) {
	c := zeroChord()
	c = c.WithInversion(Inversion1)
	if c.inversion != Inversion1 {
		t.Errorf("WithInversion failed, got %d, want %d", c.inversion, Inversion1)
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
			"('(A D_5 F5)')",
			false,
		},
		{
			"double chord",
			args{"1e/m7 1f#/m7"},
			"",
			true,
		},
		{
			"E_ augmented seventh",
			args{"E_/aug7"},
			"('(E_ G B D_5)')",
			false,
		},

		{
			"D diminished",
			args{"d/dim"},
			"('(D F A_)')",
			false,
		},
		{
			"E diminished",
			args{"e/o"},
			"('(E G B_)')",
			false,
		},
		{
			"D 7",
			args{"D/7"},
			"('(D G_ A C5)')",
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
			"('(E A_ B D5)')",
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
			"('(F A_ C#5)')",
			false,
		},
		{
			"E minor 2nd inversion",
			args{"E/m/2"},
			"('(B E5 G5)')",
			false,
		},
		{
			"E flat 3, 2nd inversion",
			args{"e_3/2"},
			"('(B_3 E_ G)')",
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
		break
	}
}
