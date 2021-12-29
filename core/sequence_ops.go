package core

func (s Sequence) Pitched(semitones int) Sequence {
	if len(s.Notes) == 0 {
		return s
	}
	if semitones == 0 {
		return s
	}
	groups := [][]Note{}
	for _, group := range s.Notes {
		changed := []Note{}
		for _, each := range group {
			changed = append(changed, each.Pitched(semitones))
		}
		groups = append(groups, changed)
	}
	return Sequence{groups}
}

func (s Sequence) Stretched(f float32) Sequence {
	if len(s.Notes) == 0 {
		return s
	}
	groups := [][]Note{}
	for _, group := range s.Notes {
		changed := []Note{}
		for _, each := range group {
			changed = append(changed, each.Stretched(f))
		}
		groups = append(groups, changed)
	}
	return Sequence{groups}
}

func (s Sequence) Reversed() Sequence {
	if len(s.Notes) == 0 {
		return s
	}
	groups := [][]Note{}
	for c := len(s.Notes) - 1; c != -1; c-- {
		groups = append(groups, s.Notes[c])
	}
	return Sequence{groups}
}

func (s Sequence) Octaved(howMuch int) Sequence {
	if len(s.Notes) == 0 {
		return s
	}
	groups := [][]Note{}
	for _, group := range s.Notes {
		changed := []Note{}
		for _, each := range group {
			changed = append(changed, each.Octaved(howMuch))
		}
		groups = append(groups, changed)
	}
	return Sequence{groups}
}

// TODO create op for this
func (s Sequence) RotatedBy(howMany int) Sequence {
	if len(s.Notes) == 0 {
		return s
	}
	direction := Right
	if howMany < 0 {
		direction = Left
		howMany = -howMany
	}
	groups := s.Notes
	for c := 0; c < howMany; c++ {
		if direction == Left {
			first := groups[0]
			groups = append(groups[1:], first)
		} else {
			last := groups[len(groups)-1]
			groups = append([][]Note{last}, groups[:len(groups)-1]...)
		}
	}
	return Sequence{groups}
}

// Split return sequences with one-note groups. Merge would produce s again.
func (s Sequence) Split() []Sequence {
	longest := 0
	for _, each := range s.Notes {
		if len(each) > longest {
			longest = len(each)
		}
	}
	all := []Sequence{}
	for i := 0; i < longest; i++ {
		layer := []Note{}
		for j, each := range s.Notes {
			if i < len(each) {
				layer = append(layer, each[i])
			} else {
				// take rest duration from first sequence
				one := all[0].Notes[j][0]
				if !one.IsPedal() {
					layer = append(layer, one.ToRest())
				}
			}
		}
		all = append(all, BuildSequence(layer))
	}
	return all
}
