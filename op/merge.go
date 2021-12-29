package op

import (
	"bytes"
	"fmt"
	"math"

	"github.com/emicklei/melrose/core"
)

type Merge struct {
	Target []core.Sequenceable
}

func (m Merge) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "merge(")
	core.AppendStorexList(&b, true, m.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

// Replaced is part of Replaceable
func (m Merge) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(m, from) {
		return to
	}
	return Join{Target: replacedAll(m.Target, from, to)}
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
		shortest := float32(math.MaxFloat32)
		// first collect the shortest duration to advance
		emptyReaders := true
		for _, each := range readers {
			dur, ok := each.durationUntilNextNote(duration)
			if ok {
				// e.g. pedals do not have a duration
				hasDuration := dur > 0
				if !hasDuration {
					// consume and add it immediately ,reader could become empty
					ns := each.take()
					groups = append(groups, ns)
				} else {
					emptyReaders = false
					// update shortest
					if dur < shortest {
						shortest = dur
					}
				}
			}
		}
		done = true
		if !emptyReaders {
			group := []core.Note{}
			// collect notes from each reader that fit into shortest; could be a rest as filler
			for _, each := range readers {
				ns, ok := each.noteUpto(duration, shortest)
				if ok {
					done = false
					group = append(group, ns...)
				}
			}
			duration += shortest
			c := compactGroup(group)
			// do not add empty ones
			if len(c) > 0 {
				groups = append(groups, c)
			}
		}
	}
	return core.Sequence{Notes: groups}
}

// compactGroup returns a group without superfluous rest notes. Does not have pedals
func compactGroup(g []core.Note) (compacted []core.Note) {
	if len(g) <= 1 {
		return g
	}
	durationShortest := float32(math.MaxFloat32)
	shortestIsRest := false
	for _, each := range g {
		f := each.DurationFactor()
		if f < durationShortest {
			durationShortest = f
			shortestIsRest = each.IsRest()
		} else {
			// non-rest notes has prio
			if f == durationShortest && shortestIsRest {
				shortestIsRest = each.IsRest()
			}
		}
	}
	hasIncludedRest := false
	for _, each := range g {
		if each.IsRest() {
			if each.DurationFactor() > durationShortest {
				continue
			}
			if shortestIsRest && !hasIncludedRest {
				hasIncludedRest = true
			} else {
				continue
			}
		}
		compacted = append(compacted, each)
	}
	return
}

// sequenceReader is to read notes and keeping the total duration for each note read
type sequenceReader struct {
	index              int           // zero based
	sequence           core.Sequence // must be a single note group sequence
	durationAtLastNote float32       // total duration factor at the end of the note at the index
}

func (r *sequenceReader) noteUpto(duration float32, shortest float32) (list []core.Note, ok bool) {
	// first check duration, could be last note
	// too soon? then no note
	if duration < r.durationAtLastNote {
		diff := r.durationAtLastNote - duration
		if shortest < diff {
			diff = shortest
		}
		var rest core.Note
		switch diff {
		case 1.5:
			rest = core.MustParseNote("1.=")
		case 1.0:
			rest = core.MustParseNote("1=")
		case 0.75:
			rest = core.MustParseNote("2.=")
		case 0.5:
			rest = core.MustParseNote("2=")
		case 0.375:
			rest = core.MustParseNote(".=")
		case 0.25:
			rest = core.MustParseNote("=")
		case 0.1875:
			rest = core.MustParseNote("8.=")
		case 0.125:
			rest = core.MustParseNote("8=")
		case 0.09375:
			rest = core.MustParseNote("16.=")
		case 0.0625:
			rest = core.MustParseNote("16=")
		default:
			return list, false
		}
		return append(list, rest), true
	}
	// any notes left?
	if r.index == len(r.sequence.Notes) {
		return list, false
	}
	r.index++
	n := r.sequence.At(r.index - 1)
	// n is a single note group
	r.durationAtLastNote += n[0].DurationFactor()
	return n, true
}

func (r *sequenceReader) durationUntilNextNote(duration float32) (float32, bool) {
	// first check duration, could be last note
	if duration < r.durationAtLastNote {
		return r.durationAtLastNote - duration, true
	}
	// any notes left?
	if r.index == len(r.sequence.Notes) {
		return 0.0, false
	}
	n := r.sequence.At(r.index)
	// n is a single note group
	return n[0].DurationFactor(), true
}

func (r *sequenceReader) take() []core.Note {
	r.index++
	n := r.sequence.At(r.index - 1)
	// n is a single note group
	r.durationAtLastNote += n[0].DurationFactor()
	return n
}
