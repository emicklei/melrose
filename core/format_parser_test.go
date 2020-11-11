package core

import (
	"testing"
)

func Test_formatParser_ParseSequence(t *testing.T) {
	for _, each := range []struct {
		in  string
		out string
	}{
		{"16.A#++ .C_-( A B ) C (D) ", "sequence('16.A♯++ .C♭- (A B) C D')"},
		{"8c#5-", "sequence('⅛C♯5-')"},
		{" ", "sequence('')"},
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

func Test_formatParser_ParseNote(t *testing.T) {
	for i, each := range []struct {
		in  string
		out string
	}{
		{"16.A#++", "note('16.A♯++')"},
		{"^", "note('^')"},
		{"<", "note('<')"},
		{">", "note('>')"},
		{"=", "note('=')"},
		{"1=", "note('1=')"},
		{"16.=", "note('16.=')"},
		{".=", "note('.=')"},
		{"4c-", "note('C-')"},
		{"8.f_", "note('⅛.F♭')"},
		{"d8", "note('D8')"},
		{"8.d+8", "note('⅛.D8+')"},
		{".e_-", "note('.E♭-')"},
		{"2e10", "note('½E10')"},
		{".2a", "note('½.A')"},
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
