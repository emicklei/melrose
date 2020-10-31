package term

import (
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"
	"github.com/emicklei/melrose/notify"
	"github.com/emicklei/tviewplus"
	"github.com/rakyll/portmidi"
)

func (m *Monitor) Open(ctx core.Context) {
	//ctx.Device().SetEchoNotes(true)
	setupConsole(m)
	m.setupDeviceSelections(ctx)
	ctx.Control().SettingNotifier(m.HandleControlSetting)
	ctx.Control().Start() // looper
	startUI(m)
}

func setupConsole(mon *Monitor) {
	notify.Console = notify.ConsoleWriter{
		DeviceIn:      mon.Received,
		DeviceOut:     mon.Sent,
		StandardOut:   mon.Console,
		StandardError: mon.Console,
	}
}

func (m *Monitor) setupDeviceSelections(ctx core.Context) {
	device := ctx.Device().(*midi.DeviceRegistry)
	inputID, outputID := device.IO()
	input := []string{" not active "}
	output := []string{" not active "}
	inputSelectionIndex, outputSelectionIndex := 0, 0 // not active
	inputSelectionIndexToID := map[int]int{0: -1}
	outputSelectionIndexToID := map[int]int{0: -1}
	for i := 0; i < portmidi.CountDevices(); i++ {
		info := portmidi.Info(portmidi.DeviceID(i))
		if info.IsInputAvailable {
			inputSelectionIndexToID[len(input)] = i
			input = append(input, fmt.Sprintf(" %d: %s/%s ", i, info.Interface, info.Name))
			if i == inputID {
				inputSelectionIndex = len(input) - 1
			}
		}
		if info.IsOutputAvailable {
			outputSelectionIndexToID[len(output)] = i
			output = append(output, fmt.Sprintf(" %d: %s/%s ", i, info.Interface, info.Name))
			if i == outputID {
				outputSelectionIndex = len(output) - 1
			}
		}
	}

	m.InputDeviceList.Set(input)
	m.InputDeviceList.Select(inputSelectionIndex) // TODO
	m.InputDeviceList.AddDependent(func(old, new tviewplus.SelectionWithIndex) {
		id := inputSelectionIndexToID[new.Index]
		fmt.Fprintf(notify.Console.StandardOut, "changing input MIDI device to %s\n", new.Value)
		device.ChangeInputDeviceID(id)
	})
	m.OutputDeviceList.Set(output)
	m.OutputDeviceList.Select(outputSelectionIndex)
	m.OutputDeviceList.AddDependent(func(old, new tviewplus.SelectionWithIndex) {
		id := outputSelectionIndexToID[new.Index]
		fmt.Fprintf(notify.Console.StandardOut, "changing output MIDI device to %s\n", new.Value)
		device.ChangeOutputDeviceID(id)
	})

	m.EchoReceivedPitchOnly.AddDependent(func(old, new bool) {
		device.EchoReceivedPitchOnly(new)
	})
}
