package core

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/emicklei/melrose/notify"
)

// Note represents a musical note.
// Notations:
// 		½C♯.3 = half duration, pitch C, sharp, octave 3, velocity default (70)
//		D     = quarter duration, pitch D, octave 4, no accidental
//      ⅛B♭  = eigth duration, pitch B, octave 4, flat
//		=     = quarter rest
//      -/+   = velocity number
// http://en.wikipedia.org/wiki/Musical_Note
type Note struct {
	Name       string // {C D E F G A B = }
	Octave     int    // [0 .. 9]
	Accidental int    // -1 Flat, +1 Sharp, 0 Normal
	Dotted     bool   // if true then reported duration is increased by half
	Velocity   int    // 1..127

	duration float32 // {0.0625,0.125,0.25,0.5,1}
}

func (n Note) Storex() string {
	return fmt.Sprintf("note('%s')", n.String())
}

// ToNote() is part of NoteConvertable
func (n Note) ToNote() Note {
	return n
}

// Replaced is part of Replaceable
func (n Note) Replaced(from, to Sequenceable) Sequenceable {
	if IsIdenticalTo(from, n) {
		return to
	}
	return n
}

var (
	Rest4 = Note{Name: "=", duration: 0.25}
)

var rest = Note{Name: "="}

func NewNote(name string, octave int, duration float32, accidental int, dot bool, velocity int) (Note, error) {
	if len(name) != 1 {
		return rest, fmt.Errorf("note must be one character, got [%s]", name)
	}
	if !strings.Contains("ABCDEFG=", name) {
		return rest, fmt.Errorf("invalid note name [ABCDEFG=]:" + name)
	}
	if octave < 0 || octave > 9 {
		return rest, fmt.Errorf("invalid octave [0..9]:" + string(octave))
	}
	switch duration {
	case 0.0625:
	case 0.125:
	case 0.25:
	case 0.5:
	case 1:
	default:
		return rest, fmt.Errorf("invalid duration [1,0.5,0.25,0.125,0.0625]:%v\n", duration)
	}

	if accidental != 0 && accidental != -1 && accidental != 1 {
		return rest, fmt.Errorf("invalid accidental :" + string(accidental))
	}

	return Note{Name: name, Octave: octave, duration: duration, Accidental: accidental, Dotted: dot, Velocity: velocity}, nil
}

func (n Note) IsRest() bool { return "=" == n.Name }

func (n Note) Length() float32 {
	if n.Dotted {
		return n.duration * 1.5
	}
	return n.duration
}

func (n Note) S() Sequence {
	return BuildSequence([]Note{n})
}

// TODO rename
func (n Note) ModifiedVelocity(velo int) Note {
	nn, _ := NewNote(n.Name, n.Octave, n.duration, n.Accidental, n.Dotted, velo)
	return nn
}

func (n Note) WithDuration(dur float64) Note {
	var duration float32
	switch dur {
	case 16:
		duration = 0.0625
	case 8:
		duration = 0.125
	case 4:
		duration = 0.25
	case 2:
		duration = 0.5
	case 1:
		duration = 1
	case 0.5:
		duration = n.duration / 2.0
	case 0.25:
		duration = n.duration / 4.0
	case 0.125:
		duration = n.duration / 8.0
	case 0.0625:
		duration = n.duration / 16.0
	default:
		notify.Panic(fmt.Errorf("cannot create note with duration [%f]", dur))
	}
	// shortest
	if duration < 0.0625 {
		duration = 0.0625
	}
	nn, err := NewNote(n.Name, n.Octave, duration, n.Accidental, n.Dotted, n.Velocity)
	if err != nil {
		notify.Panic(fmt.Errorf("cannot create note with duration [%f] because:%v", dur, err))
	}
	return nn
}

// Conversion

var noteRegexp = regexp.MustCompile("([1]?[½¼⅛12468]?)(\\.?)([CDEFGAB=])([#♯_♭]?)([0-9]?)([-+]?[-+]?[-+]?)")

// MustParseNote returns a Note by parsing the input. Panic if it fails.
func MustParseNote(input string) Note {
	n, err := ParseNote(input)
	if err != nil {
		panic("MustParseNote failed:" + err.Error())
	}
	return n
}

var N = MustParseNote

// ParseNote reads the format  <(inverse-)duration?>[CDEFGA=]<accidental?><dot?><octave?>
func ParseNote(input string) (Note, error) {
	matches := noteRegexp.FindStringSubmatch(strings.ToUpper(input))
	if matches == nil {
		return Note{}, fmt.Errorf("illegal note: [%s]", input)
	}

	var duration float32
	switch matches[1] {
	case "16":
		duration = 0.0625
	case "⅛", "8":
		duration = 0.125
	case "¼", "4":
		duration = 0.25
	case "½", "2":
		duration = 0.5
	case "1":
		duration = 1
	default:
		duration = 0.25 // quarter
	}

	dotted := matches[2] == "."

	var accidental int
	switch matches[4] {
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

	octave := 4
	if len(matches[5]) > 0 {
		i, err := strconv.Atoi(matches[5])
		if err != nil {
			return Note{}, fmt.Errorf("illegal octave: %s", matches[5])

		}
		octave = i
	}
	var velocity = Normal
	if len(matches[6]) > 0 {
		velocity = ParseVelocity(matches[6])
	}
	return NewNote(matches[3], octave, duration, accidental, dotted, velocity)
}

func ParseVelocity(plusmin string) (velocity int) {
	switch plusmin {
	case "-":
		velocity = MezzoPiano
	case "--":
		velocity = Piano
	case "---":
		velocity = Pianissimo
	case "----":
		velocity = Pianississimo
	case "+":
		velocity = MezzoForte
	case "++":
		velocity = Forte
	case "+++":
		velocity = Fortissimo
	case "++++":
		velocity = Fortississimo
	default:
		// 0
		velocity = 72
	}
	return
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
	case 0.0625:
		return "16"
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
	case 1.0:
		return "1"
	}
	return ""
}

func (n Note) Inspect(i Inspection) {
	i.Properties["length"] = n.Length()
	i.Properties["midi"] = n.MIDI()
	i.Properties["velocity"] = n.Velocity
}

func (n Note) String() string {
	return n.PrintString(PrintAsSpecified)
}

func (n Note) PrintString(sharpOrFlatKey int) string {
	var buf bytes.Buffer
	n.printOn(&buf, sharpOrFlatKey)
	return buf.String()
}

func (n Note) printOn(buf *bytes.Buffer, sharpOrFlatKey int) {
	if n.duration != 0.25 {
		buf.WriteString(n.durationf(false))
	}

	if n.Dotted {
		buf.WriteString(".")
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
	if n.Octave != 4 {
		fmt.Fprintf(buf, "%d", n.Octave)
	}
	if n.Velocity != 72 {
		switch n.Velocity {
		case Pianissimo:
			io.WriteString(buf, "---")
		case Piano:
			io.WriteString(buf, "--")
		case MezzoPiano:
			io.WriteString(buf, "-")
		case MezzoForte:
			io.WriteString(buf, "+")
		case Forte:
			io.WriteString(buf, "++")
		case Fortissimo:
			io.WriteString(buf, "+++")
		}
	}
}
