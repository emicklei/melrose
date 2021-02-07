package system

import (
	"flag"
	"log"
	"sync"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"
	"github.com/emicklei/melrose/midi/transport"

	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/server"
)

var (
	debugLogging = flag.Bool("d", false, "debug logging")
	httpPort     = flag.String("http", ":8118", "address on which to listen for HTTP requests")
)

func Setup() (core.Context, error) {
	flag.Parse()
	if *debugLogging {
		core.ToggleDebug()
	}
	transport.Initializer()

	ctx := new(core.PlayContext)
	ctx.EnvironmentVars = new(sync.Map)
	ctx.VariableStorage = dsl.NewVariableStore()
	ctx.LoopControl = core.NewBeatmaster(ctx, 120)
	reg, err := midi.NewDeviceRegistry()
	if err != nil {
		log.Fatalln("unable to initialize MIDI")
	}
	ctx.AudioDevice = reg

	if len(*httpPort) > 0 {
		// start DSL server
		go server.NewLanguageServer(ctx, *httpPort).Start()
	}

	return ctx, nil
}
