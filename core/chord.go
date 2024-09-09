package core

import (
	"bytes"
	"fmt"
	"io"
	"log"
)

// https://en.wikipedia.org/wiki/Chord_(music)
type Chord struct {
	start     Note
	inversion int // Ground,Inversion1,Inversion2,Inversion3
	interval  int // Triad,Seventh,Sixth
	quality   int // Major,Minor,Dominant,Augmented,Diminished,Suspended2,Suspended4
}

func zeroChord() Chord {
	return Chord{start: N("C"), inversion: Ground, quality: Major, interval: Triad}
}

func (c Chord) Inspect(i Inspection) {
	i.Properties["sequence"] = c.S().String()
}

func (c Chord) WithInterval(i int) Chord {
	c.interval = i
	return c
}

func (c Chord) WithInversion(i int) Chord {
	c.inversion = i
	return c
}

func (c Chord) WithQuality(q int) Chord {
	c.quality = q
	return c
}

func (c Chord) WithVelocity(v int) Chord {
	c.start = c.start.WithVelocity(v)
	return c
}

func (c Chord) WithFraction(f float32, dotted bool) Chord {
	c.start = c.start.WithFraction(f, dotted)
	return c
}

func (c Chord) String() string {
	if c.start.IsRest() {
		return c.start.String()
	}
	var b bytes.Buffer
	fmt.Fprint(&b, c.start.String())
	endsWithSlash := false
	emitSeparator := func() {
		if !endsWithSlash {
			io.WriteString(&b, "/")
		}
		endsWithSlash = false
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
			io.WriteString(&b, "dim")
		case Augmented:
			emitSeparator()
			io.WriteString(&b, "aug") // OR +  TODO
		case Suspended2:
			emitSeparator()
			io.WriteString(&b, "sus2")
		case Suspended4:
			emitSeparator()
			io.WriteString(&b, "sus4")
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

// Replaced is part of Replaceable
func (c Chord) Replaced(from, to Sequenceable) Sequenceable {
	if IsIdenticalTo(from, c) {
		return to
	}
	return c
}

// Notes returns the Note values for this chord.
func (c Chord) Notes() []Note {
	notes := []Note{c.start}
	if c.start.IsRest() ||
		c.start.IsPedalDown() ||
		c.start.IsPedalUp() ||
		c.start.IsPedalUpDown() {
		return notes
	}
	var semitones []int
	if c.interval == Triad {
		if c.quality == Augmented {
			semitones = []int{4, 8}
		} else if c.quality == Diminished {
			semitones = []int{3, 6}
		} else if c.quality == Major {
			semitones = []int{4, 7}
		} else if c.quality == Minor {
			semitones = []int{3, 7}
		} else if c.quality == Suspended2 {
			semitones = []int{2, 7}
		} else if c.quality == Suspended4 {
			semitones = []int{5, 7}
		}
	}
	if c.interval == Seventh {
		if c.quality == Augmented {
			semitones = []int{4, 8, 10}
		} else if c.quality == Diminished {
			semitones = []int{3, 6, 9}
		} else if c.quality == Minor {
			semitones = []int{3, 7, 10}
		} else if c.quality == Major {
			semitones = []int{4, 7, 11}
		} else if c.quality == Septiem {
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
		// TODO handle inversion 3
	}
	return notes
}

// C/D7/2 = C dominant 7, 2nd inversion
func ParseChord(s string) (Chord, error) {
	return newFormatParser(s).parseChord()
}

func MustParseChord(s string) Chord {
	c, err := ParseChord(s)
	if err != nil {
		log.Fatal("ParseChord failed:", err)
	}
	return c
}
