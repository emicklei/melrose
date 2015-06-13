package melrose

import "time"

type Playable interface {
	Play(Player, time.Duration)
}

type Player interface {
	PlayNote(Note, time.Duration)
	PlaySequence(Sequence, time.Duration)
}
