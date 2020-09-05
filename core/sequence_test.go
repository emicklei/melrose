package core

import (
	"fmt"
	"testing"
)

func ExampleParseSequence() {
	m, _ := ParseSequence("C (E G)")
	m2, _ := ParseSequence("C ( A )")
	m3, _ := ParseSequence("2C# (8D_ E_ F#)")
	m4, _ := ParseSequence("(C E)(.D F)(E G)")
	canto, _ := ParseSequence("B_ 8F 8D_5 8B_5 8F A_ 8E_ 8C5 8A_5 8E_")
	pedal, _ := ParseSequence("~ c d e ^ ( c d e ) ~")
	fmt.Println(m)
	fmt.Println(m2)
	fmt.Println(m3)
	fmt.Println(m4)
	fmt.Println(canto)
	fmt.Println(pedal)
	// Output:
	// C (E G)
	// C A
	// ½C♯ (⅛D♭ ⅛E♭ ⅛F♯)
	// (C E) (.D .F) (E G)
	// B♭ ⅛F ⅛D♭5 ⅛B♭5 ⅛F A♭ ⅛E♭ ⅛C5 ⅛A♭5 ⅛E♭
	// ~ C D E ^ (C D E) ~
}

func TestSequence_Storex(t *testing.T) {
	m, _ := ParseSequence("C (E G)")
	if got, want := m.Storex(), `sequence('C (E G)')`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestSequence_Duration(t *testing.T) {
	m, _ := ParseSequence("C (E G)")
	if got, want := m.NoteLength(), 0.5; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	m, _ = ParseSequence("e5 d#5 2.c#5")
	if got, want := m.NoteLength(), 1.25; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
