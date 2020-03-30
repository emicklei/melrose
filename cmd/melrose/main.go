package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/emicklei/liner"
	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/notify"
)

var (
	verbose = flag.Bool("v", false, "verbose logging")

	history  = ".melrose.history"
	varStore = dsl.NewVariableStore()
)

func main() {
	welcome()
	flag.Parse()

	// set audio
	currentDevice := setupAudio("midi")
	defer currentDevice.Close()
	melrose.SetCurrentDevice(currentDevice)

	// start REPL
	line := liner.NewLiner()
	defer line.Close()
	defer tearDown(line)
	setup(line)
	loop(line)
}

func welcome() {
	fmt.Println("\033[1;34mmelrose\033[0m" + " - v0.1")
}

var functionNames = []string{"play"}

func tearDown(line *liner.State) {
	if f, err := os.Create(history); err != nil {
		notify.Print(notify.Errorf("error writing history file:%v", err))
	} else {
		line.WriteHistory(f)
		f.Close()
	}
	fmt.Println("\033[1;34mmelrose\033[0m" + " sings bye!")
}

func setup(line *liner.State) {
	line.SetCtrlCAborts(true)
	line.SetCompleter(completeMe)
	if f, err := os.Open(history); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
}

func loop(line *liner.State) {
	for {
		entry, err := line.Prompt("𝄞 ")
		if err != nil {
			notify.Print(notify.Error(err))
			continue
		}
		if strings.HasPrefix(entry, ":") {
			// special case
			if entry == ":q" {
				goto exit
			}
			if cmd, ok := lookupCommand(entry); ok {
				if msg := cmd.Func(entry); msg != nil {
					notify.Print(msg)
				}
				continue
			}
		}
		if err := dispatch(entry); err != nil {
			notify.Print(notify.Error(err))
			// even on error, add entry to history so we can edit/fix it
		}
		line.AppendHistory(entry)
	}
exit:
}
