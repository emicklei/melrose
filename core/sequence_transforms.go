package core

func (s Sequence) NoFractions() Sequence {
	notes := [][]Note{}
	for _, g := range s.Notes {
		ng := []Note{}
		for _, n := range g {
			ng = append(ng, n.WithFraction(0.25, false))
		}
		notes = append(notes, ng)
	}
	return Sequence{Notes: notes}
}
func (s Sequence) NoDynamics() Sequence {
	notes := [][]Note{}
	for _, g := range s.Notes {
		ng := []Note{}
		for _, n := range g {
			ng = append(ng, n.WithoutDynamic())
		}
		notes = append(notes, ng)
	}
	return Sequence{Notes: notes}
}
func (s Sequence) NoRests() Sequence {
	notes := [][]Note{}
	for _, g := range s.Notes {
		ng := []Note{}
		for _, n := range g {
			if !n.IsRest() {
				ng = append(ng, n)
			}
		}
		if len(ng) > 0 {
			notes = append(notes, ng)
		}
	}
	return Sequence{Notes: notes}
}
