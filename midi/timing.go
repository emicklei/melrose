package midi

import (
	"math/rand"
	"time"
)

type TimingRandomOffset struct {
	NoteOn  int // maximum in ms
	NoteOff int // minimum in ms
	rnd     *rand.Rand
}

func newTimingOffset(on, off int) TimingRandomOffset {
	return TimingRandomOffset{
		NoteOn:  on,
		NoteOff: off,
		rnd:     rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (t TimingRandomOffset) NoteOffsets() (on int, off int) {
	return 0, 0
}
