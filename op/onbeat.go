package op

import (
	"log"

	"github.com/emicklei/melrose"
)

type OnBeat struct {
	Beat        melrose.Valueable
	Target      melrose.Valueable
	Control     melrose.LoopController
	startBeat   int64
	beatsAtLast int64
	last        interface{}
}

func NewOnBeat(beats melrose.Valueable, target melrose.Valueable, control melrose.LoopController) *OnBeat {
	return &OnBeat{
		Beat:    beats,
		Target:  target,
		Control: control,
	}
}

func (o *OnBeat) Value() interface{} {
	beats, _ := o.Control.BeatsAndBars()
	// first time
	if o.last == nil {
		o.last = o.Target.Value()
		o.startBeat = beats
		o.beatsAtLast = beats
		return o.last
	}
	log.Println(o.startBeat, o.beatsAtLast, beats, o.last)
	myBeats := int64(melrose.Int(o.Beat))

	if (beats-o.startBeat)%myBeats == 0 {
		if o.beatsAtLast == beats {
			return o.last
		}
		o.last = o.Target.Value()
		o.beatsAtLast = beats
	}
	return o.last
}
