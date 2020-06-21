package melrose

import "fmt"

// Nexter is an empty Sequence that has a sideeffect to call Value().Next() on its target when asked for the Sequence.
type Nexter struct {
	Target Valueable
}

// S is part of Sequenceable
func (n Nexter) S() Sequence {
	v := n.Target.Value()
	if t, ok := v.(Nextable); ok {
		t.Next()
	}
	return EmptySequence
}

// Storex is part of Storable
func (n Nexter) Storex() string {
	if st, ok := n.Target.(Storable); ok {
		return fmt.Sprintf("next(%s)", st.Storex())
	}
	return ""
}
