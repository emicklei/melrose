package dsl

import (
	"reflect"
	"testing"
)

func TestCleanStatements(t *testing.T) {
	source := `// ignored
a = note('c') // ignored; too
b = sequence(
    'd; e'
)
c = 3`
	want := []string{
		"// ignored",
		"a = note('c') // ignored; too",
		"b = sequence( 'd; e' )",
		"c = 3",
	}

	got, err := CleanStatements(source)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v want %#v", got, want)
	}
}

func TestFormat(t *testing.T) {
	source := "a = note('c'); b = note(\"d//e\")\n"
	want := "a = note('c')\nb = note(\"d//e\")"

	got, err := Format(source)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestCleanStatements_FirstTab(t *testing.T) {
	if _, err := CleanStatements("\ta = 1"); err == nil {
		t.Fatal("expected an error for a first-line continuation")
	}
}
