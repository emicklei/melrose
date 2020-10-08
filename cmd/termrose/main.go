package main

import (
	"log"

	"github.com/emicklei/melrose/notify"
	"github.com/emicklei/melrose/system"
)

func main() {
	ctx, err := system.Setup()
	if err != nil {
		log.Fatalln(err)
	}
	defer system.TearDown(ctx)
	mon := NewMonitor()
	setupConsole(mon)
	startUI(mon)
}

func setupConsole(mon *Monitor) {
	notify.Console = notify.ConsoleWriter{
		DeviceIn:      WriterStringHolderAdaptor{mon.Received},
		DeviceOut:     WriterStringHolderAdaptor{mon.Sent},
		StandardOut:   WriterStringHolderAdaptor{mon.Console},
		StandardError: WriterStringHolderAdaptor{mon.Console},
	}
}
