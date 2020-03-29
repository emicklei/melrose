package main

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/notify"
)

var cmdFuncMap = cmdFunctions()

func cmdFunctions() map[string]Command {
	cmds := map[string]Command{}
	cmds[":h"] = Command{Description: "show help on commands and functions", Func: showHelp}
	cmds[":s"] = Command{Description: "save memory to disk", Func: varStore.SaveMemoryToDisk}
	cmds[":l"] = Command{Description: "load memory from disk", Func: varStore.LoadMemoryFromDisk}
	cmds[":v"] = Command{Description: "show variables", Func: varStore.ListVariables}
	cmds[":m"] = Command{Description: "show MIDI information", Func: ShowDeviceInfo}
	cmds[":q"] = Command{Description: "quit"} // no Func because it is handled in the main loop
	return cmds
}

type Command struct {
	Description string
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

func showHelp(entry string) notify.Message {
	var b bytes.Buffer
	io.WriteString(&b, "\n")
	{
		funcs := dsl.EvalFunctions(varStore)
		keys := []string{}
		width := 0
		for k, _ := range funcs {
			if len(k) > width {
				width = len(k)
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			f := funcs[k]
			fmt.Fprintf(&b, "%s --- %s\n", strings.Repeat(" ", width-len(k))+k, f.Description)
		}
	}
	io.WriteString(&b, "\n")
	{
		cmds := cmdFunctions()
		keys := []string{}
		for k, _ := range cmds {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			c := cmds[k]
			fmt.Fprintf(&b, "%s --- %s\n", k, c.Description)
		}
	}
	return notify.Infof("%s", b.String())
}
