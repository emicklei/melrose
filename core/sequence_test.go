package core

import (
	"testing"
)

func TestParseSequenceBlocks(t *testing.T) {
	t.Skip()
	ParseSequence(" 8.[C (C D)]#5++ ") // 8.C#5++ (8.C#5++ 8.D#5++)
}

func TestParseSequence(t *testing.T) {
	for _, each := range []struct {
		in  string
		out string
	}{
		{"C (E G)", "C (E G)"},
		{"C ( A )", "C A"},
		{"2C# (8D_ E_ F#)", "2C# (8D_ E_ F#)"},
		{"(C E)(.D F)(E G)", "(C E) (.D F) (E G)"},
		{"B_ 8F 8D_5 8B_5 8F A_ 8E_ 8C5 8A_5 8E_", "B_ 8F 8D_5 8B_5 8F A_ 8E_ 8C5 8A_5 8E_"},
		{"> c d e ^ ( c d e ) <", "> C D E ^ (C D E) <"},
		{"< = ^ > ^ = < ^ = ^ >", "< = ^ > ^ = < ^ = ^ >"},
		//{"(c e g)~(2C 2E 2G)", "(C E G)~(2C 2E 2G)"},
	} {
		sin, err := ParseSequence(each.in)
		if err != nil {
			t.Error(err)
		} else {
			if got, want := sin.String(), each.out; got != want {
				t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
			}
		}
	}
}

func TestSequence_Storex(t *testing.T) {
	m, _ := ParseSequence("C (E G)")
	if got, want := m.Storex(), `sequence('C (E G)')`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestSequence_Duration(t *testing.T) {
	m, _ := ParseSequence("C (E G)")
	if got, want := m.DurationFactor(), 0.5; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	m, _ = ParseSequence("e5 d#5 2.c#5")
	if got, want := m.DurationFactor(), 1.25; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestSequenceLength(t *testing.T) {
	m, _ := ParseSequence("C (E G)")
	if got, want := m.Duration(120).Seconds(), 1.0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestSequencePitchlane(t *testing.T) {
	m, _ := ParseSequence("1C (8E#++ G)")
	if got, want := m.W(), "1C:0 (8E#++:5 G:7)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestSequence_Split(t *testing.T) {
	s := MustParseSequence("2C (4D 4E) (16F 16G 16A)")
	m := s.Split()
	if len(m) != 3 {
		t.Fatal()
	}
	if got, want := m[0].Storex(), "sequence('2C D 16F')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := m[1].Storex(), "sequence('2= E 16G')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := m[2].Storex(), "sequence('2= = 16A')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestSequence_SplitPedals(t *testing.T) {
	//t.Skip()
	s := MustParseSequence("> (4D 4E) <")
	m := s.Split()
	if got, want := m[0].Storex(), "sequence('> D <')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := m[1].Storex(), "sequence('E')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestSequence_RestSequence(t *testing.T) {
	s := RestSequence(2, 4)
	if got, want := len(s.Notes), 8; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := s.Storex(), "sequence('= = = = = = = =')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
