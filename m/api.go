package m

import (
	"sync"

	"github.com/emicklei/melrose"
)

func Sequence(notation string) melrose.Sequence {
	return melrose.MustParseSequence(notation)
}

func Repeat(times int, s melrose.Sequenceable) melrose.Repeat {
	return melrose.Repeat{Target: s, Times: times}
}

func Pitch(semitones int, s melrose.Sequenceable) melrose.Pitch {
	return melrose.Pitch{Target: s, Semitones: melrose.On(semitones)}
}

func Parallel(s melrose.Sequenceable) melrose.Parallel {
	return melrose.Parallel{Target: s}
}

func Note(s string) melrose.Note { return melrose.MustParseNote(s) }

func Scale(s string) melrose.Scale { return melrose.MustParseScale(s) }

func Join(s ...melrose.Sequenceable) melrose.Join {
	return melrose.Join{List: s}
}

func Go(a melrose.AudioDevice, s ...melrose.Sequenceable) {
	wg := new(sync.WaitGroup)
	for _, each := range s {
		wg.Add(1)
		go func(p melrose.Sequenceable) {
			a.Play(p)
			wg.Done()
		}(each)
	}
	wg.Wait()
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
func Serial(s melrose.Sequenceable) melrose.Serial {
	return melrose.Serial{Target: s}
}

func Channel(nr int, s melrose.Sequenceable) melrose.ChannelSelector {
	return melrose.ChannelSelector{Target: s, Number: melrose.On(nr)}
}
