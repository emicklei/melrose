package melrose

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

// https://en.wikipedia.org/wiki/Chord_(music)
type Chord struct {
	start     Note
	inversion int // Ground,Inversion1,Inversion2,Inversion3
	interval  int // Triad,Seventh,Sixth
	quality   int // Major,Minor,Dominant,Augmented,Diminished
}

func zeroChord() Chord {
	return Chord{start: N("C"), inversion: Ground, quality: Major, interval: Triad}
}

func (c Chord) String() string {
	var b bytes.Buffer
	fmt.Fprint(&b, c.start.String())
	endsWithColon := false
	emitSeparator := func() {
		if !endsWithColon {
			io.WriteString(&b, "/")
		}
		endsWithColon = false
	}

	if c.quality != Major {
		switch c.quality {
		case Minor:
			emitSeparator()
			io.WriteString(&b, "m")
		case Major:
			emitSeparator()
			io.WriteString(&b, "M")
		case Diminished:
			emitSeparator()
			io.WriteString(&b, "o")
		case Dominant:
			emitSeparator()
		case Augmented:
			emitSeparator()
			io.WriteString(&b, "A")
		}
	}
	if c.interval != Triad {
		switch c.interval {
		case Sixth:
			emitSeparator()
			io.WriteString(&b, "6")
		case Seventh:
			emitSeparator()
			io.WriteString(&b, "7")
		}
	}
	if c.inversion != Ground {
		switch c.inversion {
		case Inversion1:
			emitSeparator()
			io.WriteString(&b, "1")
		case Inversion2:
			emitSeparator()
			io.WriteString(&b, "2")
		case Inversion3:
			emitSeparator()
			io.WriteString(&b, "3")
		}
	}
	return b.String()
}

// Storex implements Storable
func (c Chord) Storex() string {
	return fmt.Sprintf("chord('%s')", c.String())
}

// S converts a Chord into a Sequence
func (c Chord) S() Sequence {
	return Sequence{[][]Note{c.Notes()}}
}

// Notes returns the Note values for this chord.
func (c Chord) Notes() []Note {
	notes := []Note{c.start}
	var semitones []int
	if c.interval == Triad {
		if Major == c.quality {
			semitones = []int{4, 7}
		} else if Minor == c.quality {
			semitones = []int{3, 7}
		}
	}
	if c.interval == Seventh {
		if c.quality == Diminished {
			semitones = []int{3, 6, 9}
		} else if Minor == c.quality {
			semitones = []int{3, 7, 10}
		} else if Major == c.quality {
			semitones = []int{4, 7, 11}
		} else if Dominant == c.quality {
			semitones = []int{4, 7, 10}
		}
	}
	for _, each := range semitones {
		next := c.start.Pitched(each)
		notes = append(notes, next)
	}
	// apply inversion
	if c.interval == Triad {
		if c.inversion == Inversion1 {
			notes = append(notes, notes[0].Octaved(1))[1:]
		}
		if c.inversion == Inversion2 {
			notes = append(notes, notes[0].Octaved(1))[1:]
			notes = append(notes, notes[0].Octaved(1))[1:]
		}
	}
	return notes
}

var chordRegexp = regexp.MustCompile("([MmoA]?)([67]?)")

//  C/D7/2 = C dominant 7, 2nd inversion
func ParseChord(s string) (Chord, error) {
	if len(s) == 0 {
		return Chord{}, errors.New("illegal chord: missing note")
	}
	parts := strings.Split(s, "/")
	start, err := ParseNote(parts[0])
	if err != nil {
		return Chord{}, err
	}
	if len(parts) == 1 {
		z := zeroChord()
		z.start = start
		return z, nil
	}
	// parts > 1
	chord := Chord{start: start}
	chord.inversion = readInversion(parts[1])

	matches := chordRegexp.FindStringSubmatch(parts[1])
	if matches == nil {
		return Chord{}, fmt.Errorf("illegal chord: [%s]", s)
	}
	switch matches[1] {
	case "M":
		chord.quality = Major
	case "m":
		chord.quality = Minor
	case "o":
		chord.quality = Diminished
	case "A":
		chord.quality = Augmented
	default:
		chord.quality = Major
	}
	switch matches[2] {
	case "6":
		chord.interval = Sixth
	case "7":
		if len(matches[1]) == 0 {
			chord.quality = Dominant
		}
		chord.interval = Seventh
	default:
		chord.interval = Triad
	}

	// parts > 2
	if len(parts) > 2 {
		chord.inversion = readInversion(parts[2])
	}
	return chord, nil
}

func readInversion(s string) int {
	switch s {
	case "1":
		return Inversion1
	case "2":
		return Inversion2
	case "3":
		return Inversion3
	default:
		return Ground
	}
}

func MustParseChord(s string) Chord {
	c, err := ParseChord(s)
	if err != nil {
		log.Fatal("ParseChord failed:", err)
	}
	return c
}
