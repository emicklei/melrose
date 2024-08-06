package core

import (
	"testing"
)

func Test_formatParser_ParseSequence(t *testing.T) {
	for _, each := range []struct {
		in  string
		out string
	}{
		{"16.A♯++ .C♭-( A B ) C (D) ", "sequence('16.A#++ .C_- (A B) C D')"},
		{"8c#5-", "sequence('8C#5-')"},
		{" ", "sequence('')"},
		{"E♭ F G A♭ B♭ C D", "sequence('E_ F G A_ B_ C D')"},
	} {
		p := newFormatParser(each.in)
		s, err := p.parseSequence()
		if err != nil {
			t.Error(err)
		}
		if got, want := s.Storex(), each.out; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
}

func Test_formatParser_ParseChordProgression(t *testing.T) {
	for _, each := range []struct {
		root string
		in   string
		out  string
	}{
		// with duration
		{"C", "2.I", "sequence('(2.C 2.E 2.G)')"},
		{"C", "2ii", "sequence('(2D 2F 2A)')"},
		{"C", "8.V7", "sequence('(8.G 8.B 8.D5 8.F5)')"},
		// quarters
		{"C", "viidim", "sequence('(B D5 F5)')"},
		{"C", "I", "sequence('(C E G)')"},
		{"C", "Idim7", "sequence('(C E_ G_ A)')"}, // A -> B__
		{"C", "Im7", "sequence('(C E_ G B_)')"},
		{"C", "IM7", "sequence('(C E G B)')"},
		{"C", "ii", "sequence('(D F A)')"},
		{"C", "iii", "sequence('(E G B)')"},
		{"C", "IV", "sequence('(F A C5)')"},
		{"C", "V", "sequence('(G B D5)')"},
		{"C", "V7", "sequence('(G B D5 F5)')"}, // G/7
		{"C", "vi", "sequence('(A C5 E5)')"},
		{"C", "vii", "sequence('(B E_5 G_5)')"},
		{"C", "Imaj7", "sequence('(C E G B)')"},
	} {
		p := newFormatParser(each.in)
		sc, err := NewScale(each.root)
		if err != nil {
			t.Fatal(err)
		}
		cs, err := p.parseChordProgression(sc)
		if err != nil {
			t.Errorf("[%s] error:%v", each.in, err)
		}
		if got, want := cs[0].S().Storex(), each.out; got != want {
			t.Errorf("[%s] got [%v:%T] want [%v:%T]", each.in, got, got, want, want)
		}
	}
}

func Test_formatParser_ParseMultipleChordProgression(t *testing.T) {
	p := newFormatParser(` I  VI  II  V `)
	sc, err := NewScale("E")
	if err != nil {
		t.Fatal(err)
	}
	cs, err := p.parseChordProgression(sc)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(cs), 4; got != want {
		t.Fatalf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := cs[0].start.Name, "E"; got != want {
		t.Fatalf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := cs[3].start.Name, "B"; got != want {
		t.Fatalf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func Test_formatParser_ParseMultipleChordProgressionWitDuration(t *testing.T) {
	p := newFormatParser(` 2I++  4VI  8II  16V `)
	sc, err := NewScale("E")
	if err != nil {
		t.Fatal(err)
	}
	cs, err := p.parseChordProgression(sc)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(cs), 4; got != want {
		t.Fatalf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := cs[0].start.Name, "E"; got != want {
		t.Fatalf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := cs[0].start.Velocity, VelocityF; got != want {
		t.Fatalf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := cs[0].start.DurationFactor(), float32(0.5); got != want {
		t.Fatalf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := cs[1].start.DurationFactor(), float32(0.25); got != want {
		t.Fatalf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := cs[2].start.DurationFactor(), float32(0.125); got != want {
		t.Fatalf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := cs[3].start.DurationFactor(), float32(0.0625); got != want {
		t.Fatalf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func Test_formatParser_ParseNote(t *testing.T) {
	for i, each := range []struct {
		in  string
		out string
	}{
		{"16.A#++", "note('16.A#++')"},
		{"^", "note('^')"},
		{"<", "note('<')"},
		{">", "note('>')"},
		{"=", "note('=')"},
		{"1=", "note('1=')"},
		{"16.=", "note('16.=')"},
		{".=", "note('.=')"},
		{"4c-", "note('C-')"},
		{"8.f_", "note('8.F_')"},
		{"d8", "note('D8')"},
		{"8.d+8", "note('8.D8+')"},
		{".e_-", "note('.E_-')"},
		{"2e10", "note('2E10')"},
		{".2a", "note('2.A')"},
	} {
		p := newFormatParser(each.in)
		s, err := p.parseNote()
		if err != nil {
			t.Fatal(err)
		}
		if got, want := s.Storex(), each.out; got != want {
			t.Errorf("[%d:%s] got [%v:%T] want [%v:%T]", i, each.in, got, got, want, want)
		}
	}
}

func TestParseTiedNotes(t *testing.T) {
	for i, each := range []struct {
		in  string
		out string
	}{
		{"8c~4c", "note('8C~C')"},
		{"8c~4c~2c", "note('8C~C~2C')"},
	} {
		p := newFormatParser(each.in)
		s, err := p.parseNote()
		if err != nil {
			t.Fatal(err)
		}
		if got, want := s.Storex(), each.out; got != want {
			t.Errorf("[%d:%s] got [%v:%T] want [%v:%T]", i, each.in, got, got, want, want)
		}
	}
}

func Test_formatParser_ParseNoteError(t *testing.T) {
	for i, each := range []struct {
		in string
	}{
		{"11A"},
		{"X"},
		{"-1"},
		{"_"},
		{"aa"},
		{"A_A"},
		{"A_5_"},
		{"..C"},
		{"c~d"},
		{"~d"},
		{"e~~e"},
		{">~<"},
		{"<~>"},
		//{"4.4C"},
	} {
		p := newFormatParser(each.in)
		n, err := p.parseNote()
		if err == nil {
			t.Fatalf("%d:%s expected an error but got:%v", i, each.in, n)
		}
		t.Logf("%s = %v", each.in, err)
	}
}

func Test_formatParser_ParseChord(t *testing.T) {
	for i, each := range []struct {
		in      string
		isError bool
		seq     string
	}{
		{"1e/m7 1f#/m7", true, ""},
		{"1e", false, "chord('1E')"},
	} {
		p := newFormatParser(each.in)
		n, err := p.parseChord()
		if err == nil && each.isError {
			t.Fatalf("%d:%s expected an error but got:%v", i, each.in, n)
		}
		if err != nil && !each.isError {
			t.Fatalf("%d:%s not expected an error but got:%v", i, each.in, err)
		}
	}
}
