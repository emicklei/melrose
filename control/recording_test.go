package control

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestRecordingStop(t *testing.T) {
	var r interface{}
	r = new(Recording)
	_, ok := r.(core.Stoppable)
	if !ok {
		t.Fail()
	}
}
