package system

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"
	"github.com/emicklei/melrose/midi/transport"
	"github.com/emicklei/melrose/notify"

	"github.com/emicklei/melrose/dsl"
)

var (
	debugLogging   = flag.Bool("d", false, "debug logging")
	errLogLocation = flag.String("log", "", "log file location")
)

func Setup(buildTag string) (core.Context, error) {
	core.BuildTag = buildTag
	flag.Parse()
	if *errLogLocation == "" {
		err := os.MkdirAll(filepath.Dir(*errLogLocation), os.ModeDir)
		if err != nil {
			notify.Errorf("failed to create log directory %s: %v", *errLogLocation, err)
			return nil, err
		}
		errOut, err := os.OpenFile(*errLogLocation, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err == nil {
			notify.Console.StandardError = io.MultiWriter(errOut, os.Stderr)
		} else {
			notify.Errorf("failed to open error log file %s: %v", *errLogLocation, err)
			return nil, err
		}
	}
	if *debugLogging {
		notify.ToggleDebug()
	}
	transport.Initializer()

	ctx := new(core.PlayContext)
	ctx.EnvironmentVars = new(sync.Map)
	ctx.VariableStorage = dsl.NewVariableStore()
	ctx.LoopControl = core.NewBeatmaster(ctx, 120)
	ctx.CapabilityFlags = core.NewCapabilities()
	reg, err := midi.NewDeviceRegistry()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize MIDI: %w", err)
	}
	ctx.AudioDevice = reg
	return ctx, nil
}
