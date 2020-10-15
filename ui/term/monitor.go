package term

import (
	"bytes"
	"fmt"

	tvp "github.com/emicklei/tviewplus"
)

type Monitor struct {
	BPM              *tvp.StringHolder
	Beat             *tvp.StringHolder
	Sent             *tvp.StringHolder
	Received         *tvp.StringHolder
	InputDeviceList  *tvp.StringListSelectionHolder
	OutputDeviceList *tvp.StringListSelectionHolder
	Console          *tvp.StringHolder
}

func NewMonitor() *Monitor {
	return &Monitor{
		BPM:              new(tvp.StringHolder),
		Beat:             new(tvp.StringHolder),
		Sent:             new(tvp.StringHolder),
		Received:         new(tvp.StringHolder),
		InputDeviceList:  new(tvp.StringListSelectionHolder),
		OutputDeviceList: new(tvp.StringListSelectionHolder),
		Console:          new(tvp.StringHolder),
	}
}

func (m *Monitor) SetBPM(bpm int) {
	m.BPM.Set(fmt.Sprintf("%d", bpm))
}

func (m *Monitor) SetBeat(b int) {
	buf := new(bytes.Buffer)
	for i := 1; i < b; i++ {
		fmt.Fprintf(buf, " ")
	}
	fmt.Fprintf(buf, "@")
	for i := b + 1; i < 5; i++ {
		fmt.Fprintf(buf, " ")
	}
	m.Beat.Set(buf.String())
}

func (m *Monitor) AppendSent(sent string) {
	m.Sent.Append(sent)
}

func (m *Monitor) SetInputDevices(list []string) {
	m.InputDeviceList.Set(list)
}

func (m *Monitor) SetOutputDevices(list []string) {
	m.OutputDeviceList.Set(list)
}

func (m *Monitor) AppendConsole(text string) {
	// TODO keep at most X chars
	m.Sent.Append(text)
}
