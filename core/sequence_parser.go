package core

import (
	"fmt"
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
	return &formatParser{scanner: s}
}

func (f *formatParser) ParseSequence() (Sequence, error) {
	stm := new(sequenceSTM)
	for {
		ch := f.scanner.Scan()
		if ch == scanner.EOF {
			break
		}
		if err := stm.accept(f.scanner.TokenText()); err != nil {
			return EmptySequence, err
		}
	}
	stm.endNote()
	return Sequence{Notes: stm.groups}, nil
}

type sequenceSTM struct {
	groups  [][]Note
	ingroup bool
	group   []Note
	note    *noteSTM
}

type noteSTM struct {
	fraction   string
	name       string
	accidental string
	octave     string
	dynamic    string
}

func (s *noteSTM) accept(lit string) error {
	s.name += lit
	return nil
}

func (s *noteSTM) Note() (Note, error) {
	return ParseNote(s.name)
}

func (s *sequenceSTM) accept(lit string) error {
	switch {
	case " " == lit:
		s.endNote()
	case "(" == lit:
		if s.ingroup {
			return fmt.Errorf("unexpected (")
		}
		s.endNote()
		s.ingroup = true
	case ")" == lit:
		if !s.ingroup {
			return fmt.Errorf("unexpected (")
		}
		s.endNote()
		if len(s.group) > 0 {
			s.groups = append(s.groups, s.group)
			s.group = []Note{}
		}
		s.ingroup = false
	default:
		if s.note == nil {
			s.note = new(noteSTM)
		}
		return s.note.accept(lit)
	}
	return nil
}

func (s *sequenceSTM) endNote() {
	// pending note?
	if s.note == nil {
		return
	}
	// note complete
	n, _ := s.note.Note()
	s.group = append(s.group, n)
	if !s.ingroup {
		s.groups = append(s.groups, s.group)
		s.group = []Note{}
	}
	s.note = nil
}
