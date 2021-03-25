package system

import (
	"flag"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"
	"github.com/emicklei/melrose/midi/transport"
	"github.com/emicklei/melrose/notify"

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
	checkVersion()

	ctx := new(core.PlayContext)
	ctx.EnvironmentVars = new(sync.Map)
	ctx.VariableStorage = dsl.NewVariableStore()
	ctx.LoopControl = core.NewBeatmaster(ctx, 120)
	reg, err := midi.NewDeviceRegistry()
	if err != nil {
		log.Fatalln("unable to initialize MIDI")
	}
	ctx.AudioDevice = reg
	return ctx, nil
}

func checkVersion() {
	if core.BuildTag == "dev" || core.BuildTag == "wasm" {
		return // ignore
	}
	v := getVersion()
	notify.Infof("you are running version %s, a newer version (%s) is available on http://melrōse.org", core.BuildTag, v)
}

func getVersion() string {
	resp, err := http.Get("https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/versions/version.txt")
	if err != nil {
		if core.IsDebug() {
			notify.Warnf("failed to fetch melrose version:%v", err)
		}
		return "?"
	}
	if resp.StatusCode != 200 || resp.Body == nil {
		if core.IsDebug() {
			notify.Warnf("failed to fetch melrose version:%v", resp)
		}
		return "?"
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		if core.IsDebug() {
			notify.Warnf("failed to fetch melrose version:%v", err)
		}
		return "?"
	}
	return strings.TrimSpace(string(data))
}
