package melrose

import (
	"fmt"
)

func ExampleSequenceParse() {
	m := ParseSequence("C C4 4C4")
	fmt.Println(m)
	// Output:
	// C C C
}

func ExampleSequenceParseGroups() {
	m := ParseSequence("C (E G)")
	m2 := ParseSequence("C ( A )")
	m3 := ParseSequence("2C# (8D_ 8E_ 2F#)")
	m4 := ParseSequence("(C E)(D. F.)(E G)")
	canto := ParseSequence("B_ 8F 8D_5 8B_5 8F A_ 8E_ 8C5 8A_5 8E_")
	fmt.Println(m)
	fmt.Println(m2)
	fmt.Println(m3)
	fmt.Println(m4)
	fmt.Println(canto)
	// Output:
	// C (E G)
	// C A
	// ½C♯ (⅛D♭ ⅛E♭ ½F♯)
	// (C E) (D. F.) (E G)
	// B♭ ⅛F ⅛D♭5 ⅛B♭5 ⅛F A♭ ⅛E♭ ⅛C5 ⅛A♭5 ⅛E♭
}
