package op

import (
	"github.com/emicklei/melrose/core"
)

type Merge struct {
	Target []core.Sequenceable
}

func (m Merge) S() core.Sequence {
	readers := []*sequenceReader{}
	for _, each := range m.Target {
		// ungroup each
		for _, other := range each.S().Split() {
			readers = append(readers, &sequenceReader{sequence: other})
		}
	}
	groups := [][]core.Note{}
	done := false
	duration := float32(0.0)
	for !done {
		// assume we are done, set false if one mergeable was found
		done = true
		group := []core.Note{}
		// shortest increments for duration
		shortest := float32(2.0)
		for _, each := range readers {
			var next []core.Note
			// assume no new note for this duration
			found := false
			nextDone := false
			// find the next group that can be merged
			for !nextDone {
				ns, ok := each.noteStartingAt(duration)
				if ok {
					// e.g. pedals do not have a duration
					hasDuration := ns[0].DurationFactor() > 0
					if hasDuration {
						next = ns
						found = true
						nextDone = true
						// update the shortest for increasing the total duration
						if df := ns[0].DurationFactor(); df < shortest {
							shortest = df
						}
					} else {
						// separate group on its own
						groups = append(groups, ns)
					}
				} else {
					// nothing left, no next found
					nextDone = true
				}
			}
			if found {
				// we found one so continue with other readers
				done = false
				group = append(group, next...)

			}
		}
		if !done {
			duration += shortest
			groups = append(groups, group)
		}
	}
	return core.Sequence{Notes: groups}
}

type sequenceReader struct {
	index                     int           // zero based
	sequence                  core.Sequence // must be a single note group sequence
	noteAtIndexEndsAtDuration float32       // total duration factor at the end of the note at the index
}

// noteStartingAt return the list and ok if the next note must start at duration of after it.
func (r *sequenceReader) noteStartingAt(duration float32) (list []core.Note, ok bool) {
	if r.index == len(r.sequence.Notes) {
		return list, false
	}
	// too soon? then no note
	if duration < r.noteAtIndexEndsAtDuration {
		return list, false
	}
	r.index++
	n := r.sequence.At(r.index - 1)
	// n is a single note group
	r.noteAtIndexEndsAtDuration += n[0].DurationFactor()
	return n, true
}
