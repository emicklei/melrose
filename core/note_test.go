package core

import (
	"fmt"
	"testing"
	"time"
)

var parsetests = []struct {
	in, name string
	octave   int
	dura     float32
	acc      int
	dot      bool
	vel      int
}{
	{"C", "C", 4, 0.25, 0, false, Normal},
	{"C3", "C", 3, 0.25, 0, false, Normal},
	{"C3", "C", 3, 0.25, 0, false, Normal},
	{"4C3", "C", 3, 0.25, 0, false, Normal},
	{"C#3", "C", 3, 0.25, 1, false, Normal},
	{"4.C#3", "C", 3, 0.25, 1, true, Normal},
	{"16C", "C", 4, 0.0625, 0, false, Normal},
	{"4C#", "C", 4, 0.25, 1, false, Normal},
	{"C#", "C", 4, 0.25, 1, false, Normal},
	{"B_", "B", 4, 0.25, -1, false, Normal}, //8
	{".F#9", "F", 9, 0.25, 1, true, Normal},
	{"1C", "C", 4, 1, 0, false, Normal},
	{"=", "=", 4, 0.25, 0, false, Normal},
	{"D++", "D", 4, 0.25, 0, false, VelocityF},
	{"D+", "D", 4, 0.25, 0, false, VelocityMF},
	{"D+++", "D", 4, 0.25, 0, false, VelocityFF},
	{"E-", "E", 4, 0.25, 0, false, VelocityMP},
	{"E--", "E", 4, 0.25, 0, false, VelocityP},
	{"E---", "E", 4, 0.25, 0, false, VelocityPP},
	{"Bo", "B", 4, 0.25, 0, false, Normal},
	{"<", "<", 0, 0, 0, false, 0},
	{"^", "^", 0, 0, 0, false, 0},
}

func TestParseNote(t *testing.T) {
	for i, each := range parsetests {
		n, err := ParseNote(each.in)
		if err != nil {
			t.Errorf("got [%v] for %s", err, each.in)
		}
		if n.Name != each.name {
			t.Fatal("name: line,exp,act", i, each.name, n.Name)
		}
		if n.Octave != each.octave {
			t.Fatal("oct: line,exp,act", i, each.octave, n.Octave)
		}
		if n.fraction != each.dura {
			t.Fatal("dur: line,exp,act", i, each.dura, n.fraction)
		}
		if n.Accidental != each.acc {
			t.Fatal("acc: line,exp,act", i, each.acc, n.Accidental)
		}
		if n.Dotted != each.dot {
			t.Fatal("dot: line,exp,act", i, each.dot, n.Dotted)
		}
		if n.Velocity != each.vel {
			t.Fatal("vel: line,exp,act", each.in, i, each.vel, n.Velocity)
		}
	}
}

var midi = []struct {
	note string
	nr   int
}{
	{"C", 60},
	{"C#", 61},
	{"C0", 12},
	{"C9", 120},
	{"G9", 127},
}

func TestMIDI(t *testing.T) {
	for _, each := range midi {
		n, _ := ParseNote(each.note)
		if n.MIDI() != each.nr {
			t.Error("line,exp,act", each.note, each.nr, n.MIDI())
		}
	}
}

func TestMIDIAll(t *testing.T) {
	for i := 12; i < 127; i++ {
		n, err := MIDItoNote(0.25, i, 1.0)
		if err != nil {
			t.Error(err)
		}
		m := n.MIDI()
		if m != i {
			t.Error("exp,act,note", i, m, n)
		}
	}
}

var pitchers = []struct {
	before string
	by     int
	after  string
}{
	{"C", 2, "D"},
	{"B", 1, "C5"},
	{"D", -2, "C"},
	{"C", 12, "C5"},
	{"C5", 2, "D5"},
	{"C#3", 0, "C#3"},
	{"C~2C", 2, "D~2D"},
}

func TestModifiedPitch(t *testing.T) {
	for i, each := range pitchers {
		n, _ := ParseNote(each.before)
		n = n.Pitched(each.by)
		if got, want := n.String(), each.after; got != want {
			t.Errorf("%d: got %v want %v", i, got, want)
		}
	}
}

func ExampleParseNote() {
	n1, _ := ParseNote("2C#3")
	n2, _ := ParseNote("2.E_2")
	n3, _ := ParseNote(".F_2++")
	e1, _ := ParseNote("2C#")
	e2, _ := ParseNote("8B_")
	e3, _ := ParseNote("4.D")
	e4, _ := ParseNote("E_")
	t1, _ := ParseNote("2c~4c")
	fmt.Println(e1)
	fmt.Println(e2)
	fmt.Println(e3)
	fmt.Println(e4)
	fmt.Println(n1)
	fmt.Println(n2)
	fmt.Println(n3)
	fmt.Println(t1)
	// Output:
	// 2C#
	// 8B_
	// .D
	// E_
	// 2C#3
	// 2.E_2
	// .F_2++
	// 2C~C
}

// Failures

func TestNote_Storex(t *testing.T) {
	n, _ := NewNote("A", 4, 0.25, 1, false, 1)
	if got, want := n.Storex(), `note('A#-----')`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNoteDurationFactor(t *testing.T) {
	for _, each := range []struct {
		note   string
		length float32
	}{
		{"=", 0.25},
		{".=", 0.375},
		{"c", 0.25},
		{".c", 0.375},
		{"2c", 0.5},
		{"2.c", 0.75},
		{"1c", 1.0},
		{"1.c", 1.5},
		{"4c", 0.25},
		{"4.c", 0.375},
		{"8c", 0.125},
		{"16c", 0.0625},
		{">", 0},
		{"^", 0},
		{"<", 0},
		{"2c~8c", 0.5 + 0.125},
	} {
		n := MustParseNote(each.note)
		if got, want := n.DurationFactor(), each.length; got != want {
			t.Errorf("got [%s] [%v:%T] want [%v:%T]", each.note, got, got, want, want)
		}
	}
}

func TestNoteWithDynamic(t *testing.T) {
	for _, each := range []struct {
		in      string
		dynamic string
		out     string
	}{
		{"c", "-", "note('C-')"},
		{"2.c#2", "--", "note('2.C#2--')"},
		{"e~2e", "+", "note('E+~2E+')"},
	} {
		nin := MustParseNote(each.in)
		before := nin.Storex()
		nout := nin.WithDynamic(each.dynamic)
		after := nin.Storex()
		if got, want := after, before; got != want {
			t.Errorf("got [%v:%T] want unchanged [%v:%T]", got, got, want, want)
		}
		if got, want := nout.Storex(), each.out; got != want {
			t.Errorf("got [%v:%T] want changed [%v:%T]", got, got, want, want)
		}
	}
}

func TestQuantizeNote(t *testing.T) {
	bpm := 120.0
	w := WholeNoteDuration(bpm)
	w16 := w / 16
	if got, want := w16, time.Duration(125)*time.Millisecond; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestQuantizedFractions(t *testing.T) {
	s := MustParseSequence("16c 16.c 8c 8.c 4c 4.c 2c 2.c 1c 1.c 16c~2c")
	for _, each := range s.Notes {
		t.Log(each[0].String(), each[0].DurationFactor())
	}
}

func TestQuantizeFraction(t *testing.T) {
	tests := []struct {
		name         string
		df           float32
		wantFraction float32
		wantDotted   bool
		wantOk       bool
	}{
		{
			"1/32",
			1.0 / 32.0,
			0.0,
			false,
			false,
		},
		{
			"1/16",
			1.0 / 16.0,
			0.0625,
			false,
			true,
		},
		{
			"3/32",
			3.0 / 32.0,
			0.09375,
			true,
			true,
		},
		{
			"1/8",
			1.0 / 8.0,
			0.125,
			false,
			true,
		},
		{
			"1.75",
			7.0 / 4.0,
			1.5,
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFraction, gotDotted, gotOk := QuantizeFraction(tt.df)
			if gotFraction != tt.wantFraction {
				t.Errorf("QuantizeFraction() gotFraction = %v, want %v", gotFraction, tt.wantFraction)
			}
			if gotDotted != tt.wantDotted {
				t.Errorf("QuantizeFraction() gotDotted = %v, want %v", gotDotted, tt.wantDotted)
			}
			if gotOk != tt.wantOk {
				t.Errorf("QuantizeFraction() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
