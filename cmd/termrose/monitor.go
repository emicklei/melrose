package main

import (
	"bytes"
	"fmt"
)

type Monitor struct {
	BPM              *StringHolder
	Beat             *StringHolder
	Sent             *StringHolder
	InputDeviceList  *StringListSelectionHolder
	OutputDeviceList *StringListSelectionHolder
}

func NewMonitor() *Monitor {
	return &Monitor{
		BPM:              new(StringHolder),
		Beat:             new(StringHolder),
		Sent:             new(StringHolder),
		InputDeviceList:  new(StringListSelectionHolder),
		OutputDeviceList: new(StringListSelectionHolder),
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
	m.Sent.Set(m.Sent.Value + sent)
}

func (m *Monitor) SetInputDevices(list []string) {
	m.InputDeviceList.Set(list)
}

func (m *Monitor) SetOutputDevices(list []string) {
	m.OutputDeviceList.Set(list)
}
