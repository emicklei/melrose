package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestFractionMap(t *testing.T) {
	pm := NewFractionMap(core.On(" 1:1 , 2:2, 3:4 "), core.MustParseSequence("c (e 4f) 8g"))
	if got, want := pm.S().Storex(), "sequence('1C (2E 2F) G')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestFractionMapNoColons(t *testing.T) {
	pm := NewFractionMap(core.On(" 1 , 2, 4 "), core.MustParseSequence("c (e 4f) 8g"))
	if got, want := pm.S().Storex(), "sequence('1C (2E 2F) G')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func Test_parseIndexFractions(t *testing.T) {
	m, err := parseIndexFractions("1:4  1:.2  1:16.")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := m[0].at, 1; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := m[0].dotted, false; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := m[0].inverseFraction, 4; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := m[1].inverseFraction, 2; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := m[1].dotted, true; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := m[2].inverseFraction, 16; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := m[2].dotted, true; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func Test_parseIndexFractions_Fails(t *testing.T) {
	for _, each := range []struct {
		input string
	}{
		{"2-2"},
		{"0:2"},
		{"1:3"},
		{"1:a"},
		{"a:a"},
	} {
		_, err := parseIndexFractions(each.input)
		if err == nil {
			t.Fail()
		} else {
			t.Log(err)
		}
	}
}
