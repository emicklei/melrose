package melrose

import "testing"

func TestJoinNoteAndSequence(t *testing.T) {
	n := N("C")
	s := S("D E")
	r := n.Join(s)
	if got, want := r.String(), "C D E"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestJoinSequenceAndNote(t *testing.T) {
	n := N("C")
	s := S("D E")
	r := s.Join(n)
	if got, want := r.String(), "D E C"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestJoinSequenceAndSequence(t *testing.T) {
	s1 := S("D E")
	s2 := S("F G")
	r := s1.Join(s2)
	if got, want := r.String(), "D E F G"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNoteAndNote(t *testing.T) {
	n1 := N("C")
	n2 := N("D")
	r := n1.Join(n2)
	if got, want := r.String(), "C D"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
