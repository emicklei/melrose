package main

import (
	"strings"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
)

var cmdFuncMap = cmdFunctions()

func cmdFunctions() map[string]Command {
	cmds := map[string]Command{}
	cmds[":h"] = Command{Description: "show help on a command or function, e.g :h seq", Func: showHelp}
	cmds[":s"] = Command{Description: "save memory to disk", Func: varStore.SaveMemoryToDisk}
	cmds[":l"] = Command{Description: "load memory from disk", Func: varStore.LoadMemoryFromDisk}
	cmds[":v"] = Command{Description: "show variables", Func: varStore.ListVariables}
	cmds[":m"] = Command{Description: "show MIDI information", Func: ShowDeviceInfo}
	cmds[":q"] = Command{Description: "quit"} // no Func because it is handled in the main loop
	return cmds
}

type Command struct {
	Description string
	Sample      string
	Func        func(entry string) notify.Message
}

func lookupCommand(entry string) (Command, bool) {
	tokens := strings.Split(entry, " ")
	if len(tokens) == 0 {
		return Command{}, false
	}
	if cmd, ok := cmdFuncMap[tokens[0]]; ok {
		return cmd, true
	}
	return Command{}, false
}

func ShowDeviceInfo(entry string) notify.Message {
	// TODO
	melrose.CurrentDevice().PrintInfo()
	return nil
}
