package transport

import "testing"

func TestHandleCallback(t *testing.T) {
	t.Skip() // TODO how to create a mock midiIN
	lis := newRtListener(nil)
	if got, want := len(lis.noteListeners), 0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	lis.handleRtEvent(nil, []byte{0x90, 60, 60}, 0.0)
	e, ok := lis.noteOn[60]
	if !ok {
		t.Fatal()
	}
	if got, want := e.note.MIDI(), 60; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	lis.handleRtEvent(nil, []byte{0x80, 60, 60}, 0.0)
	e, ok = lis.noteOn[60]
	if ok {
		t.Fatal()
	}
}
