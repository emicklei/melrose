package midi

import (
	"time"

	"github.com/emicklei/melrose/notify"
	"github.com/rakyll/portmidi"
)

type tracingMIDIStream struct {
	out     *portmidi.Stream
	notesOn map[int64]time.Time
}

func tracingMIDIStreamOn(out *portmidi.Stream) tracingMIDIStream {
	return tracingMIDIStream{
		out:     out,
		notesOn: map[int64]time.Time{},
	}
}

func (t tracingMIDIStream) WriteShort(status int64, data1 int64, data2 int64) error {
	if status&0xF0 == noteOn {
		t.notesOn[data1] = time.Now()
		// notify.Debugf("note on:%d", data1)
	} else if status&0xF0 == noteOff {
		delete(t.notesOn, data1)
		//notify.Debugf("note off:%d", data1)
	} else if status&0xF0 == controlChange {
		t.notesOn = map[int64]time.Time{}
		//notify.Debugf("control change:%d", data1)
	}
	return t.out.WriteShort(status, data1, data2)
}

func (t tracingMIDIStream) Close() error {
	return t.out.Close()
}
func (t tracingMIDIStream) Abort() error {
	return t.out.Abort()
}

func (t tracingMIDIStream) log() {
	for nr, when := range t.notesOn {
		notify.Debugf("note %d on at %v", nr, when)
	}
}
