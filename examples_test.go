package melrose

import (
	"fmt"
)

func ExampleChord() {
	// major triad
	c1, _ := ParseChord("C")
	// major sixth
	c2, _ := ParseChord("C:M6")
	// dominant seventh
	c3, _ := ParseChord("C:D7")
	// augmented triad
	c4, _ := ParseChord("C:A")
	// augmented triad seventh
	c5, _ := ParseChord("C:A7")
	// minor triad
	c6, _ := ParseChord("C:m")
	// minor sixth
	c7, _ := ParseChord("C:m6")
	// minor seventh
	c8, _ := ParseChord("C:m7")
	fmt.Println(c1.S())
	fmt.Println(c2.S())
	fmt.Println(c3.S())
	fmt.Println(c4.S())
	fmt.Println(c5.S())
	fmt.Println(c6.S())
	fmt.Println(c7.S())
	fmt.Println(c8.S())
	// Output:
	// [C E G]
	// [C E G]
	// C
	// C
	// C
	// [C E♭ G]
	// [C E♭ G]
	// [C E♭ G]
}

func ExampleChordInversion() {
	// major triad second inversion
	c1, _ := ParseChord("C:2")
	fmt.Println(c1.S())
	// Output:
	// [C E G]
}
