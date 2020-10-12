package term

import (
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
	"github.com/rakyll/portmidi"
)

func (m *Monitor) Open(ctx core.Context) {
	ctx.Device().SetEchoNotes(true)
	setupConsole(m)

	input := []string{}
	output := []string{}
	for i := 0; i < portmidi.CountDevices(); i++ {
		info := portmidi.Info(portmidi.DeviceID(i))
		if info.IsInputAvailable {
			input = append(input, fmt.Sprintf("%d:%s/%s", i, info.Interface, info.Name))
		}
		if info.IsOutputAvailable {
			output = append(output, fmt.Sprintf("%d:%s/%s", i, info.Interface, info.Name))
		}
	}
	m.InputDeviceList.Set(input)
	m.OutputDeviceList.Set(output)

	startUI(m)
}

func setupConsole(mon *Monitor) {
	notify.Console = notify.ConsoleWriter{
		DeviceIn:      WriterStringHolderAdaptor{mon.Received},
		DeviceOut:     WriterStringHolderAdaptor{mon.Sent},
		StandardOut:   WriterStringHolderAdaptor{mon.Console},
		StandardError: WriterStringHolderAdaptor{mon.Console},
	}
}
