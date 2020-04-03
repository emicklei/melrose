package main

import (
	"testing"
)

func TestCompleteMe(t *testing.T) {
	defer func() { varStore.Delete("art") }()
	varStore.Put("art", "?")
	for i, each := range []struct {
		line          string
		pos           int
		head          string
		firstcomplete string
		tail          string
	}{
		{
			"se",
			0,
			"",
			"",
			"se",
		},
		{
			"",
			0,
			"",
			"",
			"",
		},
		{
			"a = seq",
			7,
			"a = seq",
			"uence('')",
			"",
		},
		{
			"a = seq[eqs]",
			4,
			"a = ",
			"",
			"seq[eqs]",
		},
		{
			"seq",
			3,
			"seq",
			"uence('')",
			"",
		},
		{
			"a",
			1,
			"a",
			"rt",
			"",
		},
	} {
		head, c, tail := completeMe(each.line, each.pos)
		if got, want := head, each.head; got != want {
			t.Errorf("%d: got [head=%v] want [%v]", i, got, want)
		}
		firstcomplete := ""
		if len(c) > 0 {
			firstcomplete = c[0]
		}
		if got, want := firstcomplete, each.firstcomplete; got != want {
			t.Errorf("%d: got [firstcomplete=%v] want [%v]", i, got, want)
		}
		if got, want := tail, each.tail; got != want {
			t.Errorf("%d: got [tail=%v] want [%v]", i, got, want)
		}
	}
}
