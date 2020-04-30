package m

import (
	"time"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/op"
)

func Sequence(notation string) melrose.Sequence {
	return melrose.MustParseSequence(notation)
}

func Repeat(times int, s melrose.Sequenceable) melrose.Repeat {
	return melrose.Repeat{Target: s, Times: times}
}

func Reverse(s melrose.Sequenceable) melrose.Reverse {
	return melrose.Reverse{Target: s}
}

func Pitch(semitones int, s melrose.Sequenceable) melrose.Pitch {
	return melrose.Pitch{Target: s, Semitones: melrose.On(semitones)}
}

func Parallel(s melrose.Sequenceable) melrose.Parallel {
	return melrose.Parallel{Target: s}
}

func Note(s string) melrose.Note { return melrose.MustParseNote(s) }

func Scale(s string) melrose.Scale { return melrose.MustParseScale(s) }

func Join(s ...melrose.Sequenceable) op.Join {
	return op.Join{Target: s}
}

func Play(a melrose.AudioDevice, s ...melrose.Sequenceable) {
	moment := time.Now()
	for _, each := range s {
		moment = a.Play(each, melrose.Context().LoopControl.BPM(), moment)
	}
}

func Go(a melrose.AudioDevice, s ...melrose.Sequenceable) {
	moment := time.Now()
	for _, each := range s {
		a.Play(each, melrose.Context().LoopControl.BPM(), moment)
	}
}

// IndexMap creates a IndexMapper from indices.
// Example of indices: "1 (2 3 4) 5 (6 7)". One-based indexes.
func IndexMap(indices string, s melrose.Sequenceable) melrose.IndexMapper {
	return melrose.NewIndexMapper(s, indices)
}

// Chord creates a new Chord by parsing the input. See Chord for the syntax.
func Chord(s string) melrose.Chord {
	return melrose.MustParseChord(s)
}

// Serial returns a new object that serialises all the notes of the argument.
func Serial(s ...melrose.Sequenceable) melrose.Serial {
	return melrose.Serial{Target: s}
}

func Channel(nr int, s melrose.Sequenceable) melrose.ChannelSelector {
	return melrose.ChannelSelector{Target: s, Number: melrose.On(nr)}
}

// Loop returns a new loop for playing a sequence. It is not started.
func Loop(s melrose.Sequenceable) *melrose.Loop {
	return melrose.NewLoop(s)
}
