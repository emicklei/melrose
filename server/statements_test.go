package server

import "testing"

func Test_removeTrailingWithspace(t *testing.T) {
	src := `//
test
me


`
	s, l := removeTrailingWhitespace(src, 6)
	if got, want := l, 3; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := s, `//
test
me`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
