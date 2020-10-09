package system

import (
	"flag"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"

	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/server"
)

var (
	Version  = "dev"
	verbose  = flag.Bool("v", false, "verbose logging")
	httpPort = flag.String("http", ":8118", "address on which to listen for HTTP requests")
)

func Setup() (core.Context, error) {
	flag.Parse()

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
