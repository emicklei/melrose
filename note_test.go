package melrose

import (
	"fmt"
	"testing"
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
	{"D++", "D", 4, 0.25, 0, false, Forte},
	{"D+", "D", 4, 0.25, 0, false, MezzoForte},
	{"D+++", "D", 4, 0.25, 0, false, Fortissimo},
	{"E-", "E", 4, 0.25, 0, false, MezzoPiano},
	{"E--", "E", 4, 0.25, 0, false, Piano},
	{"E---", "E", 4, 0.25, 0, false, Pianissimo},
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
		if n.duration != each.dura {
			t.Fatal("dur: line,exp,act", i, each.dura, n.duration)
		}
		if n.Accidental != each.acc {
			t.Fatal("acc: line,exp,act", i, each.acc, n.Accidental)
		}
		if n.Dotted != each.dot {
			t.Fatal("dot: line,exp,act", i, each.dot, n.Dotted)
		}
		if n.Velocity != each.vel {
			t.Fatal("vel: line,exp,act", i, each.vel, n.Velocity)
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
	{"B9", 131},
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
		n := MIDItoNote(i, 1.0)
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
	{"C#3", 0, "C♯3"},
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
	fmt.Println(n1)
	fmt.Println(n2)
	fmt.Println(n3)
	// Output:
	// ½C♯3
	// ½.E♭2
	// .F♭2++
}

func ExampleParseNoteAsPrinted() {
	n1, _ := ParseNote("½C♯")
	n2, _ := ParseNote("⅛B♭")
	n3, _ := ParseNote("¼.D")
	n4, _ := ParseNote("E♭")
	fmt.Println(n1)
	fmt.Println(n2)
	fmt.Println(n3)
	fmt.Println(n4)
	// Output:
	// ½C♯
	// ⅛B♭
	// .D
	// E♭
}

// Failures

func TestFailedNewNote_BadName(t *testing.T) {
	// name string, octave int, duration float32, accidental int, dot bool
	if _, err := NewNote("Z", 4, 0.5, 0, false, 1.0); err == nil {
		t.Fail()
	}
}

func TestFailedNewNote_BadOctave(t *testing.T) {
	// name string, octave int, duration float32, accidental int, dot bool
	if _, err := NewNote("A", -1, 0.5, 0, false, 1); err == nil {
		t.Fail()
	}
	if _, err := NewNote("A", 10, 0.5, 0, false, 1); err == nil {
		t.Fail()
	}
}

func TestFailedNewNote_BadDuration(t *testing.T) {
	// name string, octave int, duration float32, accidental int, dot bool
	if _, err := NewNote("A", 4, 2, 0, false, 1); err == nil {
		t.Fail()
	}
	if _, err := NewNote("A", 4, -1, 0, false, 1); err == nil {
		t.Fail()
	}
}

func TestFailedNewNote_BadAccidental(t *testing.T) {
	// name string, octave int, duration float32, accidental int, dot bool
	if _, err := NewNote("A", 4, 0.25, -2, false, 1); err == nil {
		t.Error(err)
	}
	if _, err := NewNote("A", 4, 0.25, 2, false, 1); err == nil {
		t.Error(err)
	}
}

func TestNote_Storex(t *testing.T) {
	n, _ := NewNote("A", 4, 0.25, 1, false, 1)
	if got, want := n.Storex(), `note('A♯')`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNoteLength(t *testing.T) {
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
	} {
		n := MustParseNote(each.note)
		if got, want := n.Length(), each.length; got != want {
			t.Errorf("got [%s] [%v:%T] want [%v:%T]", each.note, got, got, want, want)
		}
	}
}
