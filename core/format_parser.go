package core

import (
	"errors"
	"fmt"
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
}

const allowedNoteNames = "abcdefgABCDEFG=<^>"

func newNoteSTM() *noteSTM {
	return &noteSTM{
		fraction:   0.25,
		dotted:     false,
		name:       "",
		accidental: 0,
		octave:     4,
		velocity:   "",
	}
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
		case "⅛", "8":
			f = 0.125
		case "¼", "4":
			f = 0.25
		case "½", "2":
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
		case "♯":
			accidental = 1
		case "♭":
			accidental = -1
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
		if strings.ContainsAny(lit, "-+") {
			s.velocity += lit
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

func (s *noteSTM) note() (Note, error) {
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
	return NewNote(s.name, s.octave, s.fraction, s.accidental, s.dotted, vel)
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
