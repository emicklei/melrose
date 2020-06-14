package main

import (
	"fmt"
	"strings"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/notify"
)

var cmdFuncMap = cmdFunctions()

func cmdFunctions() map[string]Command {
	cmds := map[string]Command{}
	cmds[":h"] = Command{Description: "show help, optional on a command or function", Func: showHelp}
	cmds[":s"] = Command{Description: "save memory to disk, optional use given filename", Func: func(args []string) notify.Message {
		return dsl.SaveMemoryToDisk(globalStore, args)
	}}
	cmds[":l"] = Command{Description: "load memory from disk, optional use given filename", Func: func(args []string) notify.Message {
		return dsl.LoadMemoryFromDisk(globalStore, args)
	}}
	cmds[":v"] = Command{Description: "show variables, optional filter on given prefix", Func: func(args []string) notify.Message {
		return dsl.ListVariables(globalStore, args)
	}}
	cmds[":k"] = Command{Description: "end all running Loops", Func: func(args []string) notify.Message {
		dsl.StopAllLoops(globalStore)
		melrose.Context().LoopControl.Reset()
		melrose.Context().AudioDevice.Reset()
		return nil
	}}
	cmds[":b"] = Command{Description: "Beat settings", Func: handleBeatSetting}
	cmds[":m"] = Command{Description: "MIDI settings", Func: handleMIDISetting}
	cmds[":q"] = Command{Description: "quit"} // no Func because it is handled in the main loop
	return cmds
}

// Command represents a REPL action that starts with c colon, ":"
type Command struct {
	Description string
	Sample      string
	Func        func(args []string) notify.Message
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

func handleMIDISetting(args []string) notify.Message {
	return melrose.Context().AudioDevice.Command(args)
}

func handleBeatSetting(args []string) notify.Message {
	l := melrose.Context().LoopControl
	fmt.Printf("[sequencer] beats per minute (BPM): %v\n", l.BPM())
	fmt.Printf("[sequencer] beats in a bar  (BIAB): %d\n", l.BIAB())
	return nil
}
