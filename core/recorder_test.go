package core

import "testing"

func TestRecordingStop(t *testing.T) {
	var r interface{}
	r = new(Recording)
	_, ok := r.(Stoppable)
	if !ok {
		t.Fail()
	}
}
