package core

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

// Note represents a musical note.
// Notations:
//
//			2.C#3 = half duration, pitch C, sharp, octave 3, velocity default
//			D     = quarter duration, pitch D, octave 4, no accidental
//	     8B_   = eighth duration, pitch B, octave 4, flat
//			=     = quarter rest
//	     -/+   = velocity number
//
// http://en.wikipedia.org/wiki/Musical_Note
type Note struct {
	Name       string // {C D E F G A B = ^ < >}
	Octave     int
	Accidental int  // -1 Flat, +1 Sharp, 0 Normal
	Dotted     bool // if true then fraction is increased by half
	Velocity   int  // 1..127

	fraction float32       // {0.03175,0.0625,0.125,0.25,0.5,1}
	duration time.Duration // if set then this overrides Dotted and fraction

	tied []Note // succeeding identical notes that are tied to this ; mostly empty
}

func (n Note) Equals(o Note) bool {
	return n.Name == o.Name &&
		n.Octave == o.Octave &&
		n.Accidental == o.Accidental &&
		n.Dotted == o.Dotted &&
		n.Velocity == o.Velocity &&
		n.fraction == o.fraction &&
		n.duration == o.duration &&
		n.HasEqualTied(o)
}

func (n Note) HasEqualTied(o Note) bool {
	if len(n.tied) != len(o.tied) {
		return false
	}
	for t := 0; t < len(n.tied); t++ {
		if !n.tied[t].Equals(o.tied[t]) {
			return false
		}
	}
	return true
}

func (n Note) Storex() string {
	return fmt.Sprintf("note('%s')", n.String())
}

// ToNote() is part of NoteConvertable
func (n Note) ToNote() (Note, error) {
	return n, nil
}

func (n Note) Fraction() float32 { return n.fraction }

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

	if !strings.Contains(allowedNoteNames, name) {
		return Rest4, fmt.Errorf("invalid note name [%s]:%s", allowedNoteNames, name)
	}
	switch frac {
	case 0.03175:
	case 0.0625:
	case 0.125:
	case 0.25:
	case 0.5:
	case 1:
	default:
		return Rest4, fmt.Errorf("invalid fraction [1,0.5,0.25,0.125,0.0625,0.03175]:%v", frac)
	}

	if accidental != 0 && accidental != -1 && accidental != 1 {
		return Rest4, fmt.Errorf("invalid accidental: %d", accidental)
	}
	return Note{Name: name, Octave: octave, fraction: frac, Accidental: accidental, Dotted: dot, Velocity: velocity}, nil
}

func MakeNote(name string, octave int, frac float32, accidental int, dot bool, velocity int) Note {
	return Note{Name: name, Octave: octave, fraction: frac, Accidental: accidental, Dotted: dot, Velocity: velocity}
}

func (n Note) IsRest() bool        { return Rest4.Name == n.Name }
func (n Note) IsPedalUp() bool     { return PedalUp.Name == n.Name }
func (n Note) IsPedalDown() bool   { return PedalDown.Name == n.Name }
func (n Note) IsPedalUpDown() bool { return PedalUpDown.Name == n.Name }
func (n Note) IsPedal() bool {
	return PedalUpDown.Name == n.Name || PedalDown.Name == n.Name || PedalUp.Name == n.Name
}

// DurationFactor is the actual duration time factor
// Only correct if n.duration is 0 and also for each tied note ; use DurationAt otherwise
func (n Note) DurationFactor() float32 {
	f := n.fraction
	if n.Dotted {
		f *= 1.5
	}
	for _, each := range n.tied {
		f += each.DurationFactor()
	}
	return f
}

func (n Note) DurationAt(bpm float64) time.Duration {
	if n.duration > 0 {
		sum := n.duration
		for _, each := range n.tied {
			sum += each.DurationAt(bpm)
		}
		return sum
	}
	return time.Duration(float32(WholeNoteDuration(bpm)) * n.DurationFactor())
}

func (n Note) S() Sequence {
	return BuildSequence([]Note{n})
}

func (n Note) WithDynamic(emphasis string) Note {
	return n.WithVelocity(ParseVelocity(emphasis))
}

func (n Note) WithoutDynamic() Note {
	return n.WithVelocity(Normal)
}

func (n Note) WithVelocity(v int) Note {
	n.Velocity = v
	if len(n.tied) == 0 {
		return n
	}
	// handle tied notes
	t := make([]Note, len(n.tied))
	for i := 0; i < len(n.tied); i++ {
		t[i] = n.tied[i].WithVelocity(v)
	}
	n.tied = t
	return n
}

func (n Note) WithFraction(f float32, dotted bool) Note {
	// TODO
	if f == 0.5*1.5 {
		n.fraction = 0.5
		n.Dotted = true
		return n
	}
	if f == 0.25*1.5 {
		n.fraction = 0.25
		n.Dotted = true
		return n
	}
	if f == 0.125*1.5 {
		n.fraction = 0.125
		n.Dotted = true
		return n
	}
	if f == 0.0625*1.5 {
		n.fraction = 0.0625
		n.Dotted = true
		return n
	}
	if f == 0.03175*1.5 {
		n.fraction = 0.03175
		n.Dotted = true
		return n
	}
	n.fraction = f
	n.Dotted = dotted
	if len(n.tied) == 0 {
		return n
	}
	// handle tied notes
	t := make([]Note, len(n.tied))
	for i := 0; i < len(n.tied); i++ {
		t[i] = n.tied[i].WithFraction(f, dotted)
	}
	n.tied = t
	return n
}

func (n Note) WithTiedNote(t Note) Note {
	n.tied = append(n.tied, t)
	return n
}

func (n Note) IsHearable() bool {
	return strings.ContainsAny(n.Name, "ABCDEFG")
}

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
	return newFormatParser(input).parseNote()
}

func ParseVelocity(plusmin string) (velocity int) {
	switch plusmin {
	case "--":
		velocity = VelocityP
	case "---":
		velocity = VelocityPP
	case "----":
		velocity = VelocityPPP
	case "-----":
		velocity = VelocityPPPP
	case "++":
		velocity = VelocityF
	case "+++":
		velocity = VelocityFF
	case "++++":
		velocity = VelocityFFF
	case "+++++":
		velocity = VelocityFFFF
	case "o":
		velocity = Normal
	case "-":
		velocity = VelocityMP
	case "+":
		velocity = VelocityMF
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
			return "_"
		}
	}
	if n.Accidental == 1 {
		if encoded {
			return "#"
		} else {
			return "#"
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

func FractionToString(f float32) string {
	switch f {
	case 0.03175:
		return "32"
	case 0.0625:
		return "16"
	case 0.125:
		return "8"
	case 0.25:
		return "4"
	case 0.5:
		return "2"
	case 1.0:
		return "1"
	}
	return ""
}

func (n Note) CheckTieableTo(t Note) error {
	if n.Name != t.Name {
		return fmt.Errorf("note name mismatch, got [%s] want [%s]", t.Name, n.Name)
	}
	if n.Octave != t.Octave {
		return fmt.Errorf("note octave mismatch, got [%d] want [%d]", t.Octave, n.Octave)
	}
	if n.Accidental != t.Accidental {
		return fmt.Errorf("note accidental mismatch, got [%d] want [%d]", t.Accidental, n.Accidental)
	}
	if n.Velocity != t.Velocity {
		return fmt.Errorf("note velocity mismatch, got [%d] want [%d]", t.Velocity, n.Velocity)
	}
	return nil
}

func (n Note) Inspect(i Inspection) {
	i.Properties["length"] = n.DurationFactor()
	i.Properties["midi"] = n.MIDI()
	i.Properties["velocity"] = n.Velocity
	i.Properties["duration"] = n.DurationAt(i.Context.Control().BPM())
}

func (n Note) String() string {
	var buf bytes.Buffer
	n.printOn(&buf, PrintAsSpecified)
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
		buf.WriteString(FractionToString(n.fraction))
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
		buf.WriteString("#")
	} else if Flat == sharpOrFlatKey && n.Accidental == 1 { // want Flat, specified in Sharp
		buf.WriteString(n.Pitched(1).Name)
		buf.WriteString("_")
	} else { // PrintAsSpecified
		buf.WriteString(n.Name)
		if n.Accidental != 0 {
			buf.WriteString(n.accidentalf(false))
		}
	}
	if n.Octave != 4 {
		fmt.Fprintf(buf, "%d", n.Octave)
	}
	if n.Velocity != Normal {
		io.WriteString(buf, VelocityToDynamic(n.Velocity))
	}
	if len(n.tied) > 0 {
		for _, each := range n.tied {
			io.WriteString(buf, "~")
			each.printOn(buf, sharpOrFlatKey)
		}
	}
}

func VelocityToDynamic(v int) string {
	if v == Normal {
		return ""
	}
	switch {
	case v <= VelocityPPPP:
		return "-----"
	case v <= VelocityPPP:
		return "----"
	case v <= VelocityPP:
		return "---"
	case v <= VelocityP:
		return "--"
	case v <= VelocityMP:
		return "-"
	case v <= Normal:
	case v <= VelocityMF:
		return "+"
	case v <= VelocityF:
		return "++"
	case v <= VelocityFF:
		return "+++"
	case v <= VelocityFFF:
		return "++++"
	case v > VelocityFFF:
		return "+++++"
	}
	return ""
}

var fractionRanges = []struct {
	fraction float32
	dotted   bool
}{
	{0.03175, false}, // 1/32
	{0.03175 * 1.5, true},
	{0.0625, false}, // 1/16
	{0.09375, true},
	{0.125, false},
	{0.1875, true},
	{0.25, false},
	{0.375, true},
	{0.5, false},
	{0.75, true},
	{1.0, false},
	{1.5, true},
	{2.0, false}, // non-exist
}

func QuantizeFraction(durationFactor float32) (fraction float32, dotted bool, ok bool) {
	last := float32(0.0)
	for i := 0; i < len(fractionRanges); i++ {
		next := fractionRanges[i]
		halfway := (last + next.fraction) / 2.0
		if durationFactor <= halfway {
			if i == 0 {
				return 0.0, false, false
			}
			prev := fractionRanges[i-1]
			return prev.fraction, prev.dotted, true
		}
		last = next.fraction
	}
	return 0.0, false, false
}
