package core

type SequenceableList struct {
	Target []Sequenceable
}

func (l SequenceableList) S() Sequence {
	if len(l.Target) == 0 {
		return EmptySequence
	}
	joined := l.Target[0].S()
	for i := 1; i < len(l.Target); i++ {
		joined = joined.SequenceJoin(l.Target[i].S())
	}
	return joined
}
