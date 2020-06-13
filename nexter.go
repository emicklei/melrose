package melrose

import "fmt"

type Nexter struct {
	Target Valueable
}

var emptySequence = Sequence{}

// S is part of Sequenceable
func (n Nexter) S() Sequence {
	_ = n.value()
	return emptySequence
}

func (n Nexter) Storex() string {
	if st, ok := n.Target.(Storable); ok {
		return fmt.Sprintf("next(%s)", st.Storex())
	}
	return fmt.Sprintf("next(%v)", n.Target)
}

// Value is part of Valueable
// TODO
func (n Nexter) value() interface{} {
	v := n.Target.Value()
	if t, ok := v.(Nextable); ok {
		return t.Next()
	}
	return nil
}
