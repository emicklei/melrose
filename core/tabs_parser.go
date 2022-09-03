package core

import (
	"fmt"
	"strconv"
	"strings"
)

var allowedTabnoteNames = "EADGeadg="

type tabnoteSTM struct {
	fraction   float32
	dotted     bool
	name       string
	accidental int
	fret       int
	velocity   string
}

func newTabNoteSTM() *tabnoteSTM {
	s := new(tabnoteSTM)
	s.reset()
	return s
}

func (s *tabnoteSTM) reset() {
	s.dotted = false
	s.fraction = 0.25
	s.name = ""
	s.fret = 0
	s.velocity = ""
}

func (s *tabnoteSTM) accept(lit string) error {
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
		if strings.ContainsAny(lit, allowedTabnoteNames) {
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
		// velocity
		if strings.ContainsAny(lit, "-o+") {
			s.velocity += lit
			return nil
		}
		// fret
		if i, err := strconv.Atoi(lit); err != nil {
			return fmt.Errorf("invalid fret, unexpected:%s", lit)
		} else {
			s.fret = i
		}
	}
	return nil
}

func (s *tabnoteSTM) currentNote() (TabNote, error) {
	vel := Normal
	if len(s.velocity) > 0 {
		vel = ParseVelocity(s.velocity)
		if vel == -1 {
			return TabNote{}, fmt.Errorf("invalid dynamic, unexpected:%s", s.velocity)
		}
	}
	return TabNote{
		Name:     s.name,
		Fret:     s.fret,
		Velocity: vel,
		fraction: s.fraction,
		Dotted:   s.dotted,
	}, nil
}

func (s *tabnoteSTM) note() (TabNote, error) {
	c, err := s.currentNote()
	if err != nil {
		return TabNote{}, err
	}
	return c, nil
}
