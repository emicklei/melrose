package core

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"
)

type formatParser struct {
	scanner *scanner.Scanner
}

func newFormatParser(src string) *formatParser {
	s := new(scanner.Scanner)
	s.Init(strings.NewReader(src))
	s.Whitespace ^= 1 << ' '
	s.Mode = scanner.ScanChars | scanner.ScanInts
	return &formatParser{scanner: s}
}

func (f *formatParser) parseNote() (Note, error) {
	var err error
	// capture scan errors
	f.scanner.Error = func(s *scanner.Scanner, m string) {
		err = errors.New(m)
	}
	stm := newNoteSTM()
	for {
		ch := f.scanner.Scan()
		if err != nil {
			return Rest4, err
		}
		if ch == scanner.EOF {
			break
		}
		if err := stm.accept(f.scanner.TokenText()); err != nil {
			return Rest4, err
		}
	}
	return stm.note()
}

func (f *formatParser) parseSequence() (Sequence, error) {
	var err error
	// capture scan errors
	f.scanner.Error = func(s *scanner.Scanner, m string) {
		err = errors.New(m)
	}
	stm := new(sequenceSTM)
	for {
		ch := f.scanner.Scan()
		if err != nil {
			return EmptySequence, err
		}
		if ch == scanner.EOF {
			break
		}
		if err := stm.accept(f.scanner.TokenText()); err != nil {
			return EmptySequence, err
		}
	}
	stm.endNote()
	return stm.sequence()
}

func (f *formatParser) parseChordProgression(s Scale) ([]Chord, error) {
	var err error
	// capture scan errors
	f.scanner.Error = func(s *scanner.Scanner, m string) {
		err = errors.New(m)
	}
	f.scanner.Mode = scanner.ScanIdents | scanner.ScanInts
	f.scanner.Whitespace = 1 << ' '
	stm := new(chordprogressionSTM)
	stm.scale = s
	for {
		ch := f.scanner.Scan()
		if err != nil {
			return []Chord{}, err
		}
		if ch == scanner.EOF {
			break
		}
		if err := stm.accept(f.scanner.TokenText()); err != nil {
			return []Chord{}, err
		}
		stm.endChord()
	}
	return stm.chords, nil
}

type sequenceSTM struct {
	groups  [][]Note
	ingroup bool
	group   []Note
	note    *noteSTM
}

type noteSTM struct {
	fraction   float32
	dotted     bool
	name       string
	accidental int
	octave     int
	velocity   string
	tied       []Note
}

type chordprogressionSTM struct {
	scale    Scale
	chords   []Chord
	index    int
	interval int
	quality  int
}

var romanChordRegex = regexp.MustCompile("([iIvV]{1,3})([Mmaj]{0,3})([dim]{0,3})(7?)")

func (s *chordprogressionSTM) accept(lit string) error {
	if lit == " " {
		s.endChord()
		return nil
	}
	matches := romanChordRegex.FindStringSubmatch(lit)
	if matches == nil {
		return fmt.Errorf("illegal chord: %s", lit)
	}
	switch matches[1] {
	case "I", "i":
		s.index = 1
	case "II", "ii":
		s.index = 2
	case "III", "iii":
		s.index = 3
	case "IV", "iv":
		s.index = 4
	case "V", "v":
		s.index = 5
	case "VI", "vi":
		s.index = 6
	case "VII", "vii":
		s.index = 7
	default:
		return fmt.Errorf("illegal roman chord: [%s]", lit)
	}
	if maj := matches[2]; len(maj) > 0 {
		if maj == "maj" {
			s.quality = Major
		}
		if maj == "m" {
			s.quality = Minor
		}
		if maj == "M" {
			s.quality = Major
		}
	}
	if dim := matches[3]; dim == "dim" {
		s.quality = Diminished
	}
	if seventh := matches[4]; seventh == "7" {
		if s.index == 5 {
			s.quality = Septiem
		}
		s.interval = Seventh
	}
	return nil
}

func (s *chordprogressionSTM) endChord() {
	ch := s.scale.ChordAt(s.index)
	if s.interval > 0 {
		ch = ch.WithInterval(s.interval)
	}
	if s.quality > 0 {
		ch = ch.WithQuality(s.quality)
	}
	s.chords = append(s.chords, ch)
}

const allowedNoteNames = "abcdefgABCDEFG=<^>"

func newNoteSTM() *noteSTM {
	s := new(noteSTM)
	s.reset()
	return s
}

func (s *noteSTM) reset() {
	s.accidental = 0
	s.dotted = false
	s.fraction = 0.25
	s.name = ""
	s.octave = 4
	s.velocity = ""
}

func (s *noteSTM) accept(lit string) error {
	if len(lit) == 0 {
		return nil
	}
	if strings.HasSuffix(lit, ".") {
		// without dot
		if err := s.accept(lit[0 : len(lit)-1]); err != nil {
			return err
		}
		lit = "."
		// proceed
	}
	if len(s.name) == 0 {
		if strings.ContainsAny(lit, allowedNoteNames) {
			if len(lit) != 1 {
				return fmt.Errorf("invalid note name, must be one character, got:%s", lit)
			}
			s.name = strings.ToUpper(lit)
			return nil
		}
		// fraction or dotted
		if lit == "." {
			if s.dotted {
				return fmt.Errorf("duration already known, got:%s", lit)
			}
			s.dotted = true
			return nil
		}
		var f float32
		switch lit {
		case "16":
			f = 0.0625
		case "8":
			f = 0.125
		case "4":
			f = 0.25
		case "2":
			f = 0.5
		case "1":
			f = 1
		default:
			return fmt.Errorf("invalid fraction or illegal note name, got:%s", lit)
		}
		if s.fraction != 0.25 {
			return fmt.Errorf("fraction already known, got:%s", lit)
		}
		s.fraction = f
		return nil
	} else {
		// name is set
		if strings.ContainsAny(lit, allowedNoteNames) {
			return fmt.Errorf("name already known, got:%s", lit)
		}
		// accidental
		var accidental int = 0
		switch lit {
		case "#":
			accidental = 1
		case "_":
			accidental = -1
		}
		if accidental != 0 {
			if s.accidental != 0 {
				return fmt.Errorf("accidental already known, unexpected:%s", lit)
			}
			s.accidental = accidental
			return nil
		}
		// velocity
		if strings.ContainsAny(lit, "-o+") {
			s.velocity += lit
			return nil
		}
		// tie
		if lit == "~" {
			n, err := s.currentNote()
			if err != nil {
				return err
			}
			s.tied = append(s.tied, n)
			s.reset()
			return nil
		}
		// octave
		if i, err := strconv.Atoi(lit); err != nil {
			return fmt.Errorf("invalid octave, unexpected:%s", lit)
		} else {
			s.octave = i
		}
	}
	return nil
}

func (s *noteSTM) currentNote() (Note, error) {
	// pedal
	switch s.name {
	case "^":
		return PedalUpDown, nil
	case "<":
		return PedalUp, nil
	case ">":
		return PedalDown, nil
	}
	vel := Normal
	if len(s.velocity) > 0 {
		vel = ParseVelocity(s.velocity)
		if vel == -1 {
			return Rest4, fmt.Errorf("invalid dynamic, unexpected:%s", s.velocity)
		}
	}
	return MakeNote(s.name, s.octave, s.fraction, s.accidental, s.dotted, vel), nil
}

func (s *noteSTM) note() (Note, error) {
	c, err := s.currentNote()
	if err != nil {
		return Rest4, err
	}
	// handle tied onces
	if len(s.tied) == 0 {
		return c, nil
	}
	// must be identical notes.
	here := s.tied[0]
	if err := here.CheckTieableTo(c); err != nil {
		return Rest4, err
	}
	for i := 1; i < len(s.tied); i++ {
		each := s.tied[i]
		if err := here.CheckTieableTo(each); err != nil {
			return Rest4, err
		}
		here = here.WithTiedNote(each)
	}
	return here.WithTiedNote(c), nil
}

func (s *sequenceSTM) accept(lit string) error {
	if len(lit) == 0 {
		return nil
	}
	switch {
	case " " == lit:
		if err := s.endNote(); err != nil {
			return err
		}
	case "(" == lit:
		if s.ingroup {
			return fmt.Errorf("unexpected (")
		}
		if err := s.endNote(); err != nil {
			return err
		}
		s.ingroup = true
	case ")" == lit:
		if !s.ingroup {
			return fmt.Errorf("unexpected (")
		}
		if err := s.endNote(); err != nil {
			return err
		}
		if len(s.group) > 0 {
			s.groups = append(s.groups, s.group)
			s.group = []Note{}
		}
		s.ingroup = false
	default:
		if s.note == nil {
			s.note = newNoteSTM()
		}
		return s.note.accept(lit)
	}
	return nil
}

func (s *sequenceSTM) endNote() error {
	// pending note?
	if s.note == nil {
		return nil
	}
	// note complete
	n, err := s.note.note()
	if err != nil {
		return err
	}
	s.group = append(s.group, n)
	if !s.ingroup {
		s.groups = append(s.groups, s.group)
		s.group = []Note{}
	}
	s.note = nil
	return nil
}

func (s *sequenceSTM) sequence() (Sequence, error) {
	return Sequence{Notes: s.groups}, nil
}
