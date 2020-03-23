package melrose

import (
	//"testing"
	"fmt"
)

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
	s = s.RotatedBy(-1)
	fmt.Println(s)

	t, _ := ParseSequence("C D E F G A B")
	t = t.RotatedBy(3)
	fmt.Println(t)

	u, _ := ParseSequence("C (D E F) G A B")
	u = u.RotatedBy(2)
	fmt.Println(u)
	// Output:
	// D E F C
	// G A B C D E F
	// A B C (D E F) G
}
