package cli

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/notify"
)

func showHelp(ctx core.Context, args []string) notify.Message {
	var b bytes.Buffer

	if len(args) == 0 {
		fmt.Fprintf(&b, "\nversion %s, syntax: %s\n", core.BuildTag, dsl.SyntaxVersion)
		fmt.Fprintf(&b, "https://melrōse.org \n")
	}

	// detect help for a command or function or it alias
	if len(args) > 0 {
		cmdfunc := strings.TrimSpace(args[0])
		if cmd, ok := cmdFunctions()[cmdfunc]; ok {
			fmt.Fprintf(&b, "%s\n----------\n", cmdfunc)
			fmt.Fprintf(&b, "%s\n\n", cmd.Description)
			fmt.Fprintf(&b, "%s\n", cmd.Sample)
			return notify.NewInfof("%s", b.String())
		}
		var fun dsl.Function
		for k, v := range dsl.EvalFunctions(ctx) {
			if k == cmdfunc {
				fun = v
				break
			}
			if v.Alias == cmdfunc {
				fun = v
				break
			}
		}
		if fun.Title != "" {
			if fun.Alias == "" {
				fmt.Fprintf(&b, "%s\n----------\n", fun.Keyword)
			} else {
				fmt.Fprintf(&b, "%s (or %q)\n----------\n", fun.Keyword, fun.Alias)
			}
			fmt.Fprintf(&b, "%s\n\n", fun.Description)
			fmt.Fprintf(&b, "%s\n", fun.Template)
			return notify.NewInfof("%s", b.String())
		}
	}
	io.WriteString(&b, "\n")
	{
		funcs := dsl.EvalFunctions(ctx)
		keys := []string{}
		width := 0
		for k, f := range funcs {
			if len(f.Title) == 0 {
				continue
			}
			if f.ControlsAudio {
				continue
			}
			if len(k) > width {
				width = len(k)
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			f := funcs[k]
			fmt.Fprintf(&b, "%s --- %s", strings.Repeat(" ", width-len(k))+k, f.Title)
			if f.Alias != "" {
				fmt.Fprintf(&b, " (alias:%s)", f.Alias)
			}
			fmt.Fprintln(&b)
		}
	}
	io.WriteString(&b, "\n")
	{
		funcs := dsl.EvalFunctions(ctx)
		keys := []string{}
		width := 0
		for k, f := range funcs {
			if len(f.Title) == 0 {
				continue
			}
			if !f.ControlsAudio {
				continue
			}
			if len(k) > width {
				width = len(k)
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			f := funcs[k]
			fmt.Fprintf(&b, "%s --- %s\n", strings.Repeat(" ", width-len(k))+k, f.Title)
		}
	}
	io.WriteString(&b, "\n")
	{
		cmds := cmdFunctions()
		keys := []string{}
		for k := range cmds {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			c := cmds[k]
			fmt.Fprintf(&b, "%s --- %s\n", k, c.Description)
		}
	}
	return notify.NewInfof("%s", b.String())
}
