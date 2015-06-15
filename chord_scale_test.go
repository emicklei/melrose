package melrose

import (
	"fmt"
)

func ExampleChord() {
	c := Chord(C(), Major)
	ebm := Chord(E(Flat), Minor)
	fmt.Println(c)
	fmt.Println(ebm)
	// Output:
	// (C E G)
	// (E♭ G♭ B♭)
}

func ExampleScale() {
	s := Scale(C(), Major)
	s2 := Scale(E(Flat), Major)
	s3 := Scale(G(Sharp), Minor)
	fmt.Println(s)
	fmt.Println(s2)
	fmt.Println(s3)
	fmt.Println(s3.PrintString(Sharp)) // TODO where do the spaces come from????
	// Output:
	// C D E F G A B C5
	// E♭ F G A♭ B♭ C5 D5 E♭5
	// G♯ B♭ B D♭5 E♭5 E5 G♭5 A♭5
	// G♯ A♯ B C♯5 D♯5 E5 F♯5 G♯5
}
