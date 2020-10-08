package system

import (
	"flag"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"

	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/server"
)

var (
	version   = "dev"
	verbose   = flag.Bool("v", false, "verbose logging")
	inputFile = flag.String("i", "", "read expressions from a file")
	httpPort  = flag.String("http", ":8118", "address on which to listen for HTTP requests")
	history   = ".melrose.history"
)

func Setup() (core.Context, error) {
	flag.Parse()

	device, err := midi.Open()
	if err != nil {
		return nil, err
	}

	ctx := new(core.PlayContext)
	ctx.EnvironmentVars = map[string]string{}
	ctx.AudioDevice = device
	ctx.VariableStorage = dsl.NewVariableStore()
	ctx.LoopControl = core.NewBeatmaster(ctx, 120)

	if len(*httpPort) > 0 {
		// start DSL server
		go server.NewLanguageServer(ctx, *httpPort).Start()
	}

	return ctx, nil
}
