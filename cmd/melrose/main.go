package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/emicklei/melrose"
	"github.com/peterh/liner"
)

var (
	deviceID = flag.String("d", "midi", "set the audio device")
	verbose  = flag.Bool("v", false, "verbose logging")

	history       = ".melrose.history"
	currentDevice melrose.AudioDevice
)

func main() {
	welcome()
	flag.Parse()

	// set audio
	currentDevice = setupAudio(*deviceID)
	defer currentDevice.Close()

	// start REPL
	line := liner.NewLiner()
	defer line.Close()
	defer tearDown(line)
	setup(line)
	loop(line)
}

func welcome() {
	fmt.Println("\033[1;34mmelrose\033[0m" + " - v0.0.1")
}

var functionNames = []string{"play"}

func tearDown(line *liner.State) {
	if f, err := os.Create(history); err != nil {
		printError(fmt.Sprintf("error writing history file:%v", err))
	} else {
		line.WriteHistory(f)
		f.Close()
	}
	fmt.Println("\033[1;34mmelrose\033[0m" + " says bye!")
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
			printError(err)
			continue
		}
		if strings.HasPrefix(entry, ":") {
			// special case
			if entry == ":q" {
				goto exit
			}
			if cmd, ok := cmdFuncMap[entry]; ok {
				cmd.Func()
				continue
			}
		}
		if err := dispatch(entry); err != nil {
			printError(err)
			// even on error, add entry to history so we can edit/fix it
		}
		line.AppendHistory(entry)
	}
exit:
}
