package melrose

import "sync"

type Functions struct{}

func (f Functions) Sequence(notation string) Sequenceable {
	return MustParseSequence(notation)
}

func (f Functions) Repeat(times int, s Sequenceable) Sequenceable {
	return Repeat{Target: s, Times: times}
}

func (f Functions) Pitch(semitones int, s Sequenceable) Sequenceable {
	return Pitch{Target: s, Semitones: semitones}
}

func (f Functions) Parallel(s Sequenceable) Sequenceable {
	return Parallel{Target: s}
}

func (f Functions) Note(s string) Note { return MustParseNote(s) }

func (f Functions) Join(s ...Sequenceable) Sequenceable {
	return Join{List: s}
}

func (f Functions) Go(a AudioDevice, s ...Sequenceable) {
	wg := new(sync.WaitGroup)
	for _, each := range s {
		wg.Add(1)
		go func(p Sequenceable) {
			a.Play(p, true)
			wg.Done()
		}(each)
	}
	wg.Wait()
}

// IndexMap creates a IndexMapper from indices.
// Example of indices: "1 (2 3 4) 5 (6 7)". One-based indexes.
func (f Functions) IndexMap(s Sequenceable, indices string) Sequenceable {
	return IndexMapper{Target: s, Indices: parseIndices(indices)}
}

// Chord creates a new Chord by parsing the input. See Chord for the syntax.
func (f Functions) Chord(s string) Chord {
	return MustParseChord(s)
}

// Serial returns a new object that serialises all the notes of the argument.
func (f Functions) Serial(s Sequenceable) Sequenceable {
	return Serial{Target: s}
}
