package melrose

import "time"

type Playable interface {
	Play(Player, time.Duration)
}

type Player interface {
	PlayNote(Note)
	PlaySequence(s Sequence)
}
