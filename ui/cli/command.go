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
	cmds[":k"] = Command{Description: "stop all sound and loops", Func: func(ctx core.Context, args []string) notify.Message {
		dsl.StopAllPlayables(ctx)
		ctx.Device().Reset()
		return nil
	}}
	cmds[":b"] = Command{Description: "beat settings", Func: handleBeatSetting}
	cmds[":m"] = Command{Description: "MIDI settings", Func: handleMIDISetting}
	cmds[":q"] = Command{Description: "quit"} // no Func because it is handled in the main loop
	cmds[":d"] = Command{Description: "toggle debug lines", Func: handleToggleDebug}
	cmds[":p"] = Command{Description: "list all running", Func: handleListAllRunning}
	cmds[":e"] = Command{Description: "echo MIDI", Func: handleEchoNotes}
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
	notify.PrintHighlighted("sequencer settings:")
	fmt.Printf("beats per minute (BPM) = %v\n", l.BPM())
	fmt.Printf("beats in a bar  (BIAB) = %d\n", l.BIAB())
	return nil
}

func handleToggleDebug(ctx core.Context, args []string) notify.Message {
	if notify.ToggleDebug() {
		return notify.NewInfof("debug enabled")
	}
	return notify.NewInfof("debug not enabled")
}

func handleListAllRunning(ctx core.Context, args []string) notify.Message {
	running := []core.Stoppable{}
	for _, v := range ctx.Variables().Variables() {
		if s, ok := v.(core.Stoppable); ok {
			if s.IsPlaying() {
				running = append(running, s)
			}
		}
	}
	for _, each := range running {
		fmt.Printf("%s = %s\n", ctx.Variables().NameFor(each), core.Storex(each))
	}
	return nil
}

func handleEchoNotes(ctx core.Context, args []string) notify.Message {
	return ctx.Device().Command(append([]string{"e"}, args...))
}
