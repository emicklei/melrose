package melrose

import "time"

type Playable interface {
	Play(Player, time.Duration)
}

type Player interface {
	PlayNote(Note, time.Duration)
	PlaySequence(s Sequence, singleNoteDuration time.Duration)
}
