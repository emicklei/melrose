package melrose

import (
	"fmt"
)

func ExampleChord() {
	c := C().Chord()
	ebm := E(Flat).Chord(Minor)
	fmt.Println(c.S())
	fmt.Println(ebm.S())
	// Output:
	// (C E G)
	// (E♭ G♭ B♭)
}

func ExampleScale() {
	s := C().Scale()
	s2 := E(Flat).Scale()
	s3 := G(Sharp).Scale(Minor)
	fmt.Println(s.S())
	fmt.Println(s2.S())
	fmt.Println(s3.S())
	fmt.Println(s3.S().PrintString(Sharp))
	// Output:
	// C D E F G A B
	// E♭ F G A♭ B♭ C5 D5
	// G♯ B♭ B D♭5 E♭5 E5 G♭5
	// G♯ A♯ B C♯5 D♯5 E5 F♯5
}
