package term

import (
	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

func (m *Monitor) Open(ctx core.Context) {
	ctx.Device().SetEchoNotes(true)
	setupConsole(m)
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
