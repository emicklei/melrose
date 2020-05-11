package op

import "github.com/emicklei/melrose"

type NoteMap struct {
	Target  melrose.Valueable
	Indices []int
}

func NewNoteMapper(indices string, note melrose.Valueable) NoteMap {
	idx := []int{}
	for _, each := range parseIndices(indices) {
		idx = append(idx, each[0])
	}
	return NoteMap{
		Target:  note,
		Indices: idx,
	}
}

func (n NoteMap) S() melrose.Sequence {
	notelike, ok := n.Target.Value().(melrose.NoteConvertable)
	if !ok {
		// TODO warning here?
		return melrose.Sequence{}
	}
	max := 0
	min := 10000
	for _, each := range n.Indices {
		if each > max {
			max = each
		} else if each < min {
			min = each
		}
	}
	notes := make([]melrose.Note, max)
	for i := range notes {
		notes[i] = melrose.Rest4
	}
	note := notelike.ToNote()
	for _, each := range n.Indices {
		notes[each-1] = note
	}
	return melrose.BuildSequence(notes)
}

type NoteMerge struct {
	Target []melrose.Valueable
	Count  int
}

func NewNoteMerge(count int, maps []melrose.Valueable) NoteMerge {
	return NoteMerge{
		Count:  count,
		Target: maps,
	}
}

var restGroup = []melrose.Note{melrose.Rest4}

func (m NoteMerge) S() melrose.Sequence {
	groups := [][]melrose.Note{}
	for g := 1; g <= m.Count; g++ {
		group := []melrose.Note{}
		for _, eachMapVal := range m.Target {
			eachMap, ok := eachMapVal.Value().(NoteMap)
			if !ok {
				// TODO warning here?
				return melrose.Sequence{}
			}
			for _, eachIndex := range eachMap.Indices {
				if eachIndex == g {
					notelike, ok := eachMap.Target.Value().(melrose.NoteConvertable)
					if !ok {
						// TODO warning here?
						return melrose.Sequence{}
					}
					group = append(group, notelike.ToNote())
					break
				}
			}
		}
		if len(group) == 0 {
			groups = append(groups, restGroup)
		} else {
			groups = append(groups, group)
		}
	}
	return melrose.Sequence{Notes: groups}
}
