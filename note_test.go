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
}{
	{"C", "C", 4, 0.25, 0, false},
	{"C3", "C", 3, 0.25, 0, false},
	{"C3", "C", 3, 0.25, 0, false},
	{"4C3", "C", 3, 0.25, 0, false},
	{"4C#3", "C", 3, 0.25, 1, false},
	{"4C#.3", "C", 3, 0.25, 1, true},
	{"4C#", "C", 4, 0.25, 1, false},
	{"C#", "C", 4, 0.25, 1, false},
	{"B_", "B", 4, 0.25, -1, false}, //8
	{"F#.9", "F", 9, 0.25, 1, true},
}

func TestParseNote(t *testing.T) {
	for i, each := range parsetests {
		n := ParseNote(each.in)
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
		n := ParseNote(each.note)
		if n.MIDI() != each.nr {
			t.Error("line,exp,act", each.note, each.nr, n.MIDI())
		}
	}
}

func TestMIDIAll(t *testing.T) {
	for i := 12; i < 127; i++ {
		n := MIDItoNote(i)
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
}

func TestModifiedPitch(t *testing.T) {
	for i, each := range pitchers {
		n := ParseNote(each.before).ModifiedPitch(each.by)
		if n.String() != each.after {
			t.Fatal("line,exp,act", i, each.after, n.String())
		}
	}
}

func ExampleParseNote() {
	fmt.Println(ParseNote("2C#3"))
	fmt.Println(ParseNote("2E_.2"))
	// Output:
	// ½C♯3
	// ½E♭.2
}

func ExampleParseNoteAsPrinted() {
	fmt.Println(ParseNote("½C♯"))
	fmt.Println(ParseNote("⅛B♭"))
	fmt.Println(ParseNote("¼D."))
	fmt.Println(ParseNote("E♭"))
	// Output:
	// ½C♯
	// ⅛B♭
	// D.
	// E♭
}

func ExampleC() {
	fmt.Println(C(Sharp), C())
	// Output:
	// C♯ C
}

func ExampleSharp() {
	fmt.Println(C().Sharp(), D().Sharp(), E().Sharp(), F().Sharp(), G().Sharp(), A().Sharp(), B().Sharp())
	// Output:
	// C♯ D♯ E♯ F♯ G♯ A♯ B♯
}

func ExampleFlat() {
	fmt.Println(B(Flat), B())
	// Output:
	// B♭ B
}

func ExampleBs() {
	fmt.Println(B(Sharp, Dot), B(Flat), B())
	// Output:
	// B♯. B♭ B
}

func ExampleCSharpOctave() {
	fmt.Println(C().Sharp().ModifiedOctave(-1))
	// Output:
	// C♯3
}

func ExampleSharpHalf() {
	fmt.Println(C(Sharp, Half))
	// Output:
	// ½C♯
}

func ExampleFlatEight() {
	fmt.Println(B(Flat, Eight))
	// Output:
	// ⅛B♭
}

func ExampleOctaveUp() {
	fmt.Println(C().ModifiedOctave(1), C().ModifiedOctave(-1))
	// Output:
	// C5 C3
}

func ExampleAdjecentName() {
	fmt.Println(C().AdjecentName(Left, 1))
	fmt.Println(C().AdjecentName(Right, 7))
	// Output:
	// B
	// C
}

func ExamplePrintString_Flat() {
	fmt.Println(C(Sharp).PrintString(Flat))
	// Output:
	// D♭
}

func ExamplePrintString_Sharp() {
	fmt.Println(E(Flat).PrintString(Sharp))
	// Output:
	// D♯
}

// Failures

func PanicDetector(t *testing.T) {
	if r := recover(); r != nil {
		t.Log("Good!, panic situation detected: ", r)
	} else {
		t.Fatal("Bummer!, should have panic-ed")
	}
}

func TestFailedNewNote_BadName(t *testing.T) {
	defer PanicDetector(t)
	// name string, octave int, duration float32, accidental int, dot bool
	t.Log(newNote("Z", 4, 0.5, 0, false))
}

func TestFailedNewNote_BadOctave(t *testing.T) {
	defer PanicDetector(t)
	// name string, octave int, duration float32, accidental int, dot bool
	t.Log(newNote("A", -1, 0.5, 0, false))
	t.Log(newNote("A", 10, 0.5, 0, false))
}

func TestFailedNewNote_BadDuration(t *testing.T) {
	defer PanicDetector(t)
	// name string, octave int, duration float32, accidental int, dot bool
	t.Log(newNote("A", 4, 2, 0, false))
	t.Log(newNote("A", 4, -1, 0, false))
}

func TestFailedNewNote_BadAccidental(t *testing.T) {
	defer PanicDetector(t)
	// name string, octave int, duration float32, accidental int, dot bool
	t.Log(newNote("A", 4, 0.25, -2, false))
	t.Log(newNote("A", 4, 0.25, 2, false))
}
