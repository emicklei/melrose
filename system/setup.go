package system

import (
	"flag"
	"log"
	"sync"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"
	"github.com/emicklei/melrose/midi/transport"

	"github.com/emicklei/melrose/dsl"
)

var (
	debugLogging = flag.Bool("d", false, "debug logging")
)

func Setup(buildTag string) (core.Context, error) {
	core.BuildTag = buildTag
	flag.Parse()
	if *debugLogging {
		core.ToggleDebug()
	}
	transport.Initializer()
	//checkVersion()

	ctx := new(core.PlayContext)
	ctx.EnvironmentVars = new(sync.Map)
	ctx.VariableStorage = dsl.NewVariableStore()
	ctx.LoopControl = core.NewBeatmaster(ctx, 120)
	ctx.CapabilityFlags = core.NewCapabilities()
	reg, err := midi.NewDeviceRegistry()
	if err != nil {
		log.Fatalln("unable to initialize MIDI")
	}
	ctx.AudioDevice = reg
	return ctx, nil
}
