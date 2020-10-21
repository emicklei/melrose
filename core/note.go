package core

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Note represents a musical note.
// Notations:
// 		½.C♯3 = half duration, pitch C, sharp, octave 3, velocity default (70)
//		D     = quarter duration, pitch D, octave 4, no accidental
//      ⅛B♭  = eigth duration, pitch B, octave 4, flat
//		=     = quarter rest
//      -/+   = velocity number
// http://en.wikipedia.org/wiki/Musical_Note
type Note struct {
	Name       string // {C D E F G A B = ^ < >}
	Octave     int    // [0 .. 9]
	Accidental int    // -1 Flat, +1 Sharp, 0 Normal
	Dotted     bool   // if true then fraction is increased by half
	Velocity   int    // 1..127

	fraction float32       // {0.0625,0.125,0.25,0.5,1}
	duration time.Duration // if set then this overrides Dotted and fraction
}

func (n Note) Storex() string {
	return fmt.Sprintf("note('%s')", n.String())
}

// ToNote() is part of NoteConvertable
func (n Note) ToNote() (Note, error) {
	return n, nil
}

func (n Note) ToRest() Note {
	return Note{
		Name:       "=",
		Octave:     n.Octave,
		Accidental: n.Accidental,
		Dotted:     n.Dotted,
		Velocity:   n.Velocity,
		fraction:   n.fraction,
		duration:   n.duration,
	}
}

// Replaced is part of Replaceable
func (n Note) Replaced(from, to Sequenceable) Sequenceable {
	if IsIdenticalTo(from, n) {
		return to
	}
	return n
}

var (
	Rest4        = Note{Name: "=", fraction: 0.25}
	PedalUpDown  = Note{Name: "^", fraction: 0}
	PedalDown    = Note{Name: ">", fraction: 0}
	PedalUp      = Note{Name: "<", fraction: 0}
	ZeroDuration = time.Duration(0)
)

const validNoteNames = "ABCDEFG=<^>"

func NewNote(name string, octave int, frac float32, accidental int, dot bool, velocity int) (Note, error) {
	if len(name) != 1 {
		return Rest4, fmt.Errorf("note must be one character, got [%s]", name)
	}
	// pedal check
	switch name {
	case "^":
		return PedalUpDown, nil
	case ">":
		return PedalDown, nil
	case "<":
		return PedalUp, nil
	}

	if !strings.Contains(validNoteNames, name) {
		return Rest4, fmt.Errorf("invalid note name [%s]:%s", validNoteNames, name)
	}
	if octave < 0 || octave > 9 {
		return Rest4, fmt.Errorf("invalid octave [0..9]: %d", octave)
	}
	switch frac {
	case 0.0625:
	case 0.125:
	case 0.25:
	case 0.5:
	case 1:
	default:
		return Rest4, fmt.Errorf("invalid fraction [1,0.5,0.25,0.125,0.0625]:%v", frac)
	}

	if accidental != 0 && accidental != -1 && accidental != 1 {
		return Rest4, fmt.Errorf("invalid accidental: %d", accidental)
	}

	return Note{Name: name, Octave: octave, fraction: frac, Accidental: accidental, Dotted: dot, Velocity: velocity}, nil
}

func (n Note) IsRest() bool        { return Rest4.Name == n.Name }
func (n Note) IsPedalUp() bool     { return PedalUp.Name == n.Name }
func (n Note) IsPedalDown() bool   { return PedalDown.Name == n.Name }
func (n Note) IsPedalUpDown() bool { return PedalUpDown.Name == n.Name }

// DurationFactor is the actual duration time factor
func (n Note) DurationFactor() float32 {
	if n.Dotted {
		return n.fraction * 1.5
	}
	return n.fraction
}

func (n Note) S() Sequence {
	return BuildSequence([]Note{n})
}

func (n Note) WithDynamic(emphasis string) Note {
	n.Velocity = ParseVelocity(emphasis)
	return n
}

func (n Note) WithVelocity(v int) Note {
	n.Velocity = v
	return n
}

func (n Note) WithFraction(f float32, dotted bool) Note {
	n.fraction = f
	n.Dotted = dotted
	return n
}

func (n Note) IsHearable() bool {
	return strings.IndexAny(n.Name, "ABCDEFG") != -1
}

// Conversion
// https://regoio.herokuapp.com/
var noteRegexp = regexp.MustCompile("([1]?[½¼⅛12468]?)(\\.?)([CDEFGAB=<^>])([#♯_♭]?)([0-9]?)([-+]?[-+]?[-+]?)")

// MustParseNote returns a Note by parsing the input. Panic if it fails.
func MustParseNote(input string) Note {
	n, err := ParseNote(input)
	if err != nil {
		panic("MustParseNote failed:" + err.Error())
	}
	return n
}

var N = MustParseNote

// ParseNote reads the format  <(inverse-)duration?>[CDEFGA=<^>]<accidental?><dot?><octave?>
func ParseNote(input string) (Note, error) {
	matches := noteRegexp.FindStringSubmatch(strings.ToUpper(input))
	if matches == nil {
		return Note{}, fmt.Errorf("illegal note: [%s]", input)
	}

	var fraction float32
	switch matches[1] {
	case "16":
		fraction = 0.0625
	case "⅛", "8":
		fraction = 0.125
	case "¼", "4":
		fraction = 0.25
	case "½", "2":
		fraction = 0.5
	case "1":
		fraction = 1
	default:
		fraction = 0.25 // quarter
	}

	dotted := matches[2] == "."

	// pedal
	switch matches[3] {
	case "^":
		return PedalUpDown, nil
	case "<":
		return PedalUp, nil
	case ">":
		return PedalDown, nil
	}

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
	return NewNote(matches[3], octave, fraction, accidental, dotted, velocity)
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
	case "0":
		velocity = Normal
	default:
		// invalid
		velocity = -1
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

func (n Note) NonFractionBasedDuration() (time.Duration, bool) {
	if n.duration > 0 {
		return n.duration, true
	}
	return ZeroDuration, false
}

func (n Note) durationf(encoded bool) string {
	switch n.fraction {
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
	i.Properties["length"] = n.DurationFactor()
	i.Properties["midi"] = n.MIDI()
	i.Properties["velocity"] = n.Velocity
	wholeNoteDuration := WholeNoteDuration(i.Context.Control().BPM())
	i.Properties["duration"] = time.Duration(float32(wholeNoteDuration) * n.DurationFactor())
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
	if n.IsPedalUp() {
		buf.WriteString(PedalUp.Name)
		return
	}
	if n.IsPedalDown() {
		buf.WriteString(PedalDown.Name)
		return
	}
	if n.IsPedalUpDown() {
		buf.WriteString(PedalUpDown.Name)
		return
	}

	if n.fraction != 0.25 {
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
