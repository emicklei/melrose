package core

import "fmt"

type SetBPM struct {
	bpm     Valueable
	control LoopController
}

func NewBPM(bpm Valueable, ctr LoopController) SetBPM {
	return SetBPM{bpm: bpm, control: ctr}
}

// S has the side effect of setting the BPM unless BPM is zero
func (s SetBPM) S() Sequence {
	f := Float(s.bpm)
	if f > 0.0 {
		s.control.SetBPM(float64(f))
	}
	return EmptySequence
}

// Inspect implements Inspectable
func (s SetBPM) Inspect(i Inspection) {
	i.Properties["bpm"] = fmt.Sprintf("%.2f", Float(s.bpm))
}

// Storex implements Storable
func (s SetBPM) Storex() string {
	return fmt.Sprintf("bpm(%v)", s.bpm)
}
