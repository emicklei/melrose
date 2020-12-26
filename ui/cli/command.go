package cli

import (
	"fmt"
	"strings"

	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/notify"
)

var cmdFuncMap = cmdFunctions()

func cmdFunctions() map[string]Command {
	cmds := map[string]Command{}
	cmds[":h"] = Command{Description: "show help, optional on a command or function", Func: showHelp}
	cmds[":v"] = Command{Description: "show variables, optional filter on given prefix", Func: func(ctx core.Context, args []string) notify.Message {
		return dsl.ListVariables(ctx.Variables(), args)
	}}
	cmds[":k"] = Command{Description: "end all running playables", Func: func(ctx core.Context, args []string) notify.Message {
		dsl.StopAllPlayables(ctx)
		return nil
	}}
	cmds[":b"] = Command{Description: "beat settings", Func: handleBeatSetting}
	cmds[":m"] = Command{Description: "MIDI settings", Func: handleMIDISetting}
	cmds[":q"] = Command{Description: "quit"} // no Func because it is handled in the main loop
	cmds[":d"] = Command{Description: "toggle debug lines", Func: handleToggleDebug}
	return cmds
}

// Command represents a REPL action that starts with c colon, ":"
type Command struct {
	Description string
	Sample      string
	Func        func(ctx core.Context, args []string) notify.Message
}

func lookupCommand(cmd string) (Command, bool) {
	if len(cmd) == 0 {
		return Command{}, false
	}
	if cmd, ok := cmdFuncMap[strings.ToLower(cmd)]; ok {
		return cmd, true
	}
	return Command{}, false
}

func handleMIDISetting(ctx core.Context, args []string) notify.Message {
	return ctx.Device().Command(args)
}

func handleBeatSetting(ctx core.Context, args []string) notify.Message {
	l := ctx.Control()
	fmt.Printf("[sequencer] beats per minute (BPM): %v\n", l.BPM())
	fmt.Printf("[sequencer] beats in a bar  (BIAB): %d\n", l.BIAB())
	return nil
}

func handleToggleDebug(ctx core.Context, args []string) notify.Message {
	if core.ToggleDebug() {
		return notify.Infof("debug enabled")
	}
	return notify.Infof("debug not enabled")
}
