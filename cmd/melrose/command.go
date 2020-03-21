package main

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

var cmdFuncMap = cmdFunctions()

type Command struct {
	Description string
	Func        func()
}

func cmdFunctions() map[string]Command {
	cmds := map[string]Command{}
	cmds[":h"] = Command{Description: "show help on commands and functions", Func: showHelp}
	cmds[":s"] = Command{Description: "save memory to disk", Func: saveMemoryToDisk}
	cmds[":l"] = Command{Description: "load memory from disk", Func: loadMemoryFromDisk}
	cmds[":v"] = Command{Description: "show variables", Func: listVariables}
	cmds[":d"] = Command{Description: "switch audio device", Func: listVariables}
	cmds[":x"] = Command{Description: "show examples", Func: listVariables}
	cmds[":q"] = Command{Description: "quit"}
	return cmds
}

func showHelp() {
	var b bytes.Buffer
	io.WriteString(&b, "\n")
	{
		funcs := evalFunctions()
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
	printInfo(b.String())
}
