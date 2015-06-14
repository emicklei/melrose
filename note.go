package melrose

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Note represents a musical note.
// Notations:
// 		½C♯.3 = half+half duration, pitch C, sharp, octave 3
//		D     = quarter duration, pitch D, octave 4, no accidental
//      ⅛B♭   = eigth duration, pitch B, octave 4, flat
//		r     = quarter rest
// http://en.wikipedia.org/wiki/Musical_Note
type Note struct {
	Name       string  // {C D E F G A B r}
	Octave     int     // [0 .. 9]
	duration   float32 // {0.125,0.25,0.5,1}
	Accidental int     // -1 Flat, +1 Sharp, 0 Normal
	Dotted     bool    // if true then reported duration is increased by half
}

// Constructors

func C(modifiers ...int) Note {
	return NewNote("C").Modified(modifiers...)
}
func D(modifiers ...int) Note {
	return NewNote("D").Modified(modifiers...)
}
func E(modifiers ...int) Note {
	return NewNote("E").Modified(modifiers...)
}
func F(modifiers ...int) Note {
	return NewNote("F").Modified(modifiers...)
}
func G(modifiers ...int) Note {
	return NewNote("G").Modified(modifiers...)
}
func A(modifiers ...int) Note {
	return NewNote("A").Modified(modifiers...)
}
func B(modifiers ...int) Note {
	return NewNote("B").Modified(modifiers...)
}
func Rest(modifiers ...int) Note {
	return NewNote("r").Modified(modifiers...)
}

func NewNote(name string) Note {
	return newNote(name, 4, 0.25, 0, false)
}

// todo make it return error i.o panic
func newNote(name string, octave int, duration float32, accidental int, dot bool) Note {
	if len(name) != 1 {
		panic("note name too long [1]:" + name)
	}
	if !strings.Contains("ABCDEFGr", name) {
		panic("invalid note name [ABCDEFGr]:" + name)
	}
	if octave < 0 || octave > 9 {
		panic("invalid octave [0..9]:" + string(octave))
	}
	switch duration {
	case 0.125:
	case 0.25:
	case 0.5:
	case 1:
	default:
		panic(fmt.Sprintf("invalid duration [1,0.5,0.25,0.125]:%v", duration))
	}

	if accidental != 0 && accidental != -1 && accidental != 1 {
		panic("invalid accidental :" + string(accidental))
	}

	return Note{Name: name, Octave: octave, duration: duration, Accidental: accidental, Dotted: dot}
}

// Accessors

func (n Note) IsFlat() bool  { return n.Accidental == -1 }
func (n Note) IsSharp() bool { return n.Accidental == 1 }
func (n Note) IsRest() bool  { return "r" == n.Name }

func (n Note) Equals(other Note) bool {
	// quick check first
	if n.Octave != other.Octave {
		return false
	}
	// pitch independent check
	if n.DurationFactor() != other.DurationFactor() {
		return false
	}
	return n.MIDI() == other.MIDI()
}

func (n Note) DurationFactor() float32 {
	if n.Dotted {
		return n.duration + 0.5
	}
	return n.duration
}

func (n Note) Repeated(howMany int) Sequence {
	s := Sequence{}
	for i := 0; i < howMany; i++ {
		s = s.Append(n)
	}
	return s
}

func (n Note) Frequency() int {
	// http://en.wikipedia.org/wiki/Musical_Note
	// A4 == 440Hz (scientific pitch notation)
	panic("not implemented")
	return 440
}

// Modified applies modifiers on the Note and returns the new result
func (n Note) Modified(modifiers ...int) Note {
	modified := n
	for _, each := range modifiers {
		switch each {
		case Sharp:
			modified = modified.Sharp()
		case Flat:
			modified = modified.Flat()
		case Eight:
			modified = modified.Eight()
		case Half:
			modified = modified.Half()
		case Quarter:
			modified = modified.Quarter()
		case Whole:
			modified = modified.Whole()
		case Dot:
			modified = modified.Dot()
		}
	}
	return modified
}

// Pitch

func (n Note) Sharp() Note {
	return newNote(n.Name, n.Octave, n.duration, 1, n.Dotted)
}
func (n Note) Flat() Note {
	return newNote(n.Name, n.Octave, n.duration, -1, n.Dotted)
}

// Pitched creates a new Note with a pitch by a (positive or negative) number of semi tones
func (n Note) Pitched(howManySemitones int) Note {
	simple := MIDItoNote(n.MIDI() + howManySemitones)
	return newNote(simple.Name, simple.Octave, n.duration, simple.Accidental, n.Dotted)
}

// Major returns the note left or right on the Major Scale by an offset
func (n Note) Major(offset int) Note {
	// C=0
	nameIndex := strings.Index(NonRestNoteNames, n.Name)
	// semitones on the scale
	nameOffset := noteMidiOffsets[nameIndex]
	majors := offset % 7
	scales := offset / 7
	return n.Pitched(-nameOffset).Octaved(scales).Pitched(noteMidiOffsets[majors])
}

func (n Note) Octaved(howmuch int) Note {
	return newNote(n.Name, n.Octave+howmuch, n.duration, n.Accidental, n.Dotted)
}

// Duration

func (n Note) Eight() Note {
	return newNote(n.Name, n.Octave, 0.125, n.Accidental, n.Dotted)
}

func (n Note) Quarter() Note {
	return newNote(n.Name, n.Octave, 0.25, n.Accidental, n.Dotted)
}

func (n Note) Half() Note {
	return newNote(n.Name, n.Octave, 0.5, n.Accidental, n.Dotted)
}

func (n Note) Whole() Note {
	return newNote(n.Name, n.Octave, 1, n.Accidental, false)
}

func (n Note) Dot() Note {
	return newNote(n.Name, n.Octave, n.duration, n.Accidental, true)
}

func (n Note) ModifiedDuration(by float32) Note {
	return newNote(n.Name, n.Octave, n.duration+by, n.Accidental, n.Dotted)
}

// Conversion

var noteRegexp = regexp.MustCompile("([½¼⅛1248]?)([CDEFGABr])([#♯_♭]?)(\\.?)([0-9]?)")

// ParseNote reads the format  <(inverse-)duration?>[CDEFGABr]<accidental?><dot?><octave?>
func ParseNote(input string) Note {
	matches := noteRegexp.FindStringSubmatch(strings.ToUpper(input))
	if matches == nil {
		panic("illegal note:" + input)
	}

	var duration float32
	switch matches[1] {
	case "⅛":
		duration = 0.125
	case "8":
		duration = 0.125
	case "¼":
		duration = 0.25
	case "4":
		duration = 0.25
	case "½":
		duration = 0.5
	case "2":
		duration = 0.5
	case "1":
		duration = 1
	default:
		duration = 0.25
	}

	var accidental int
	switch matches[3] {
	case "#":
		accidental = 1
	case "♯":
		accidental = 1
	case "♭":
		accidental = -1
	case "_":
		accidental = -1
	default:
		accidental = 0
	}

	dotted := matches[4] == "."

	octave := 4
	if len(matches[5]) > 0 {
		i, err := strconv.Atoi(matches[5])
		if err != nil {
			panic("illegal octave:" + matches[5])
		}
		octave = i
	}

	return newNote(matches[2], octave, duration, accidental, dotted)
}

// Formatting

func (n Note) accidentalf(encoded bool) string {
	if n.Accidental == -1 {
		if encoded {
			return "b"
		} else {
			return "♭"
		}
	}
	if n.Accidental == 1 {
		if encoded {
			return "#"
		} else {
			return "♯"
		}
	}
	return ""
}

func (n Note) durationf(encoded bool) string {
	switch n.duration {
	case 0.125:
		if encoded {
			return "8"
		} else {
			return "⅛"
		}
	case 0.25:
		if encoded {
			return "4"
		} else {
			return "¼"
		}
	case 0.5:
		if encoded {
			return "2"
		} else {
			return "½"
		}
	}
	return ""
}

func (n Note) EncodeOn(buf *bytes.Buffer, sharpOrFlatKey int) {
	buf.WriteString(n.durationf(true))
	buf.WriteString(n.Name)
	if !n.IsRest() {
		if n.Accidental != 0 {
			buf.WriteString(n.accidentalf(true))
		}
		if n.Dotted {
			buf.WriteString(".")
		}
		if n.Octave != 4 {
			buf.WriteString(fmt.Sprintf("%d", n.Octave))
		}
	}
}

func (n Note) String() string {
	return n.PrintString(PrintAsSpecified)
}

func (n Note) PrintString(sharpOrFlatKey int) string {
	var buf bytes.Buffer
	n.PrintOn(&buf, sharpOrFlatKey)
	return buf.String()
}

func (n Note) PrintOn(buf *bytes.Buffer, sharpOrFlatKey int) {
	if n.duration != 0.25 {
		buf.WriteString(n.durationf(false))
	}
	if n.IsRest() {
		buf.WriteString(n.Name)
		return
	}
	if Sharp == sharpOrFlatKey && n.Accidental == -1 { // want Sharp, specified in Flat
		buf.WriteString(n.Pitched(-1).Name)
		buf.WriteString("♯")
	} else if Flat == sharpOrFlatKey && n.Accidental == 1 { // want Flat, specified in Sharp
		buf.WriteString(n.Pitched(1).Name)
		buf.WriteString("♭")
	} else { // PrintAsSpecified
		buf.WriteString(n.Name)
		if n.Accidental != 0 {
			buf.WriteString(n.accidentalf(false))
		}
	}
	if n.Dotted {
		buf.WriteString(".")
	}
	if n.Octave != 4 {
		buf.WriteString(fmt.Sprintf("%d", n.Octave))
	}
}

func (n Note) Play(p Player, t time.Duration) {
	p.PlayNote(n, t)
}
