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
