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
	{"1C", "C", 4, 1, 0, false},
}

func TestParseNote(t *testing.T) {
	for i, each := range parsetests {
		n, _ := ParseNote(each.in)
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

func TestAdjecentName(t *testing.T) {

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
	{"C5", 2, "D5"},
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

func TestMajorOffset(t *testing.T) {
	for i, each := range []struct {
		before string
		by     int
		after  string
	}{
		{"C", 1, "D"},
		{"C", 7, "C5"},
		{"C", 3, "F"},
		{"C", 10, "F5"},
	} {
		n, _ := ParseNote(each.before)
		n = n.Major(each.by)
		if got, want := n.String(), each.after; got != want {
			t.Errorf("%d: got %v want %v", i, got, want)
		}
	}
}

func ExampleParseNote() {
	n1, _ := ParseNote("2C#3")
	n2, _ := ParseNote("2E_.2")
	fmt.Println(n1)
	fmt.Println(n2)
	// Output:
	// ½C♯3
	// ½E♭.2
}

func ExampleParseNoteAsPrinted() {
	n1, _ := ParseNote("½C♯")
	n2, _ := ParseNote("⅛B♭")
	n3, _ := ParseNote("¼D.")
	n4, _ := ParseNote("E♭")
	fmt.Println(n1)
	fmt.Println(n2)
	fmt.Println(n3)
	fmt.Println(n4)
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
	fmt.Println(C().Sharp().Octaved(-1))
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
	fmt.Println(C().Octaved(1), C().Octaved(-1))
	// Output:
	// C5 C3
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
