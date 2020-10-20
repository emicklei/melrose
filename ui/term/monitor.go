package term

import (
	"fmt"

	"github.com/emicklei/melrose/core"
	tvp "github.com/emicklei/tviewplus"
)

type Monitor struct {
	BPM                   *tvp.StringHolder
	BIAB                  *tvp.StringHolder
	Sent                  *tvp.StringHolder
	Received              *tvp.StringHolder
	EchoReceivedPitchOnly *tvp.BoolHolder
	InputDeviceList       *tvp.StringListSelectionHolder
	OutputDeviceList      *tvp.StringListSelectionHolder
	Console               *tvp.StringHolder
}

func NewMonitor() *Monitor {
	return &Monitor{
		BPM:                   new(tvp.StringHolder),
		BIAB:                  new(tvp.StringHolder),
		Sent:                  new(tvp.StringHolder),
		Received:              new(tvp.StringHolder),
		EchoReceivedPitchOnly: new(tvp.BoolHolder),
		InputDeviceList:       new(tvp.StringListSelectionHolder),
		OutputDeviceList:      new(tvp.StringListSelectionHolder),
		Console:               new(tvp.StringHolder),
	}
}

func (m *Monitor) HandleControlSetting(control core.LoopController) {
	bpm := control.BPM()
	biab := control.BIAB()
	m.BPM.Set(fmt.Sprintf("%.2f", bpm))
	m.BIAB.Set(fmt.Sprintf("%d", biab))
}
