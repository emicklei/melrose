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
		s, err := p.ParseSequence()
		if err != nil {
			t.Error(err)
		}
		if got, want := s.Storex(), each.out; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
}
