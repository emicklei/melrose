package melrose

import (
	//"testing"
	"fmt"
)

func ExampleStretchBy() {
	s, _ := ParseSequence("C D E F")
	s = StretchBy{2}.Transform(s)
	fmt.Println(s)
	// Output:
	// ½C ½D ½E ½F
}

func ExamplePitchBy() {
	s, _ := ParseSequence("C D E F")
	s = PitchBy{2}.Transform(s)
	t := PitchBy{-4}.Transform(s)
	fmt.Println(s)
	fmt.Println(t)
	// Output:
	// D E G♭ G
	// B♭3 C D E♭
}

func ExampleGroupBy() {
	s, _ := ParseSequence("C D E F")
	s = GroupBy{[]int{2, 2}}.Transform(s)
	fmt.Println(s)

	t, _ := ParseSequence("C D E F G A B")
	t = GroupBy{[]int{3, 3, 1}}.Transform(t)
	fmt.Println(t)
	// Output:
	// (C D) (E F)
	// (C D E) (F G A) B
}

func ExampleRotateBy() {
	s, _ := ParseSequence("C D E F")
	s = RotateBy{Left, 1}.Transform(s)
	fmt.Println(s)

	t, _ := ParseSequence("C D E F G A B")
	t = RotateBy{Right, 3}.Transform(t)
	fmt.Println(t)

	u, _ := ParseSequence("C (D E F) G A B")
	u = RotateBy{Right, 2}.Transform(u)
	fmt.Println(u)
	// Output:
	// D E F C
	// G A B C D E F
	// A B C (D E F) G
}
