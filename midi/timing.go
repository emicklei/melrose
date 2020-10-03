package midi

import (
	"math/rand"
	"time"
)

type TimingModifier interface {
	NoteOn() time.Duration
	NoteOff() time.Duration
}

type IntRange struct {
	Min int
	Max int
}

func (i IntRange) Len() int {
	return i.Max - i.Min
}

type TimingRandomOffset struct {
	noteOn  IntRange // ms
	noteOff IntRange // ms
	rnd     *rand.Rand
}

func newTimingOffset(minOn, maxOn, minOff, maxOff int) TimingRandomOffset {
	return TimingRandomOffset{
		noteOn:  IntRange{Min: minOn, Max: maxOn},
		noteOff: IntRange{Min: minOff, Max: maxOff},
		rnd:     rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (t TimingRandomOffset) NoteOn() time.Duration {
	return time.Duration(t.noteOn.Min+t.rnd.Intn(t.noteOn.Len())) * time.Millisecond

}

func (t TimingRandomOffset) NoteOff() time.Duration {
	return time.Duration(t.noteOff.Min+t.rnd.Intn(t.noteOff.Len())) * time.Millisecond
}

type NoOffset struct{}

func (t NoOffset) NoteOff() time.Duration { return time.Duration(0) }
func (t NoOffset) NoteOn() time.Duration  { return time.Duration(0) }
func (t NoOffset) Offset() int            { return 0 }
