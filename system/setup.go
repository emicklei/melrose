package system

import (
	"flag"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"

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

	ctx := new(core.PlayContext)
	ctx.EnvironmentVars = map[string]string{}
	ctx.VariableStorage = dsl.NewVariableStore()
	ctx.LoopControl = core.NewBeatmaster(ctx, 120)

	device, err := midi.Open(ctx)
	if err != nil {
		return nil, err
	}
	ctx.AudioDevice = device

	if len(*httpPort) > 0 {
		// start DSL server
		go server.NewLanguageServer(ctx, *httpPort).Start()
	}

	return ctx, nil
}
