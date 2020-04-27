package melrose

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Note represents a musical note.
// Notations:
// 		½C♯.3 = half duration, pitch C, sharp, octave 3, velocity default (70)
//		D     = quarter duration, pitch D, octave 4, no accidental
//      ⅛B♭  = eigth duration, pitch B, octave 4, flat
//		=     = quarter rest
//      -     = velocity -20%
//      --    = velocity -40%
//      ---   = velocity -80%
//      +     = velocity +20%
//      ++    = velocity +40%
//      +++   = velocity +80%
// http://en.wikipedia.org/wiki/Musical_Note
type Note struct {
	Name       string // {C D E F G A B = }
	Octave     int    // [0 .. 9]
	Accidental int    // -1 Flat, +1 Sharp, 0 Normal
	Dotted     bool   // if true then reported duration is increased by half

	duration       float32 // {0.0625,0.125,0.25,0.5,1}
	velocityFactor float32 // {0.8,0.6,0.2,1.0,1.2,1.6,1.8}
}

func (n Note) Storex() string {
	return fmt.Sprintf("note('%s')", n.String())
}

var rest = Note{Name: "="}

func NewNote(name string, octave int, duration float32, accidental int, dot bool, velocityFactor float32) (Note, error) {
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

	return Note{Name: name, Octave: octave, duration: duration, Accidental: accidental, Dotted: dot, velocityFactor: velocityFactor}, nil
}

func (n Note) IsRest() bool { return "=" == n.Name }

func (n Note) DurationFactor() float32 {
	if n.Dotted {
		return n.duration * 1.5
	}
	return n.duration
}

func (n Note) VelocityFactor() float32 {
	if n.velocityFactor == 0 {
		return 1.0
	}
	return n.velocityFactor
}

func (n Note) S() Sequence {
	return BuildSequence([]Note{n})
}

func (n Note) ModifiedVelocity(velo float32) Note {
	nn, _ := NewNote(n.Name, n.Octave, n.duration, n.Accidental, n.Dotted, velo)
	return nn
}

// Conversion

var noteRegexp = regexp.MustCompile("([1]?[½¼⅛12468]?)([CDEFGAB=])([#♯_♭]?)(\\.?)([0-9]?)([-+]?[-+]?[-+]?)")

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
			return Note{}, fmt.Errorf("illegal octave: %s", matches[5])

		}
		octave = i
	}
	var velocity float32 = 1.0 // audio decides
	if len(matches[6]) > 0 {
		velocity = ParseVelocity(matches[6])
	}
	return NewNote(matches[2], octave, duration, accidental, dotted, velocity)
}

func ParseVelocity(plusmin string) (velocity float32) {
	// 0.8,0.6,0.2,1.2,1.6,1.8
	switch plusmin {
	case "-":
		velocity = F_MezzoPiano
	case "--":
		velocity = F_Piano
	case "---":
		velocity = F_Pianissimo
	case "+":
		velocity = F_MezzoForte
	case "++":
		velocity = F_Forte
	case "+++":
		velocity = F_Fortissimo
	default:
		// 0
		velocity = 1.0
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
		fmt.Fprintf(buf, "%d", n.Octave)
	}
	if n.velocityFactor != 1.0 {
		switch n.velocityFactor {
		case F_Pianissimo:
			io.WriteString(buf, "---")
		case F_Piano:
			io.WriteString(buf, "--")
		case F_MezzoPiano:
			io.WriteString(buf, "-")
		case F_MezzoForte:
			io.WriteString(buf, "+")
		case F_Forte:
			io.WriteString(buf, "++")
		case F_Fortissimo:
			io.WriteString(buf, "+++")
		}
	}
}
