package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/js"
	"github.com/emicklei/melrose/notify"
	"github.com/peterh/liner"
)

var (
	version   = "v0.1"
	verbose   = flag.Bool("v", false, "verbose logging")
	inputFile = flag.String("i", "", "read expressions from a file")
	httpPort  = flag.String("http", ":8118", "address on which to listen for HTTP requests")

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

	// process file if given
	if len(*inputFile) > 0 {
		processInputFile(*inputFile)
	}

	if len(*httpPort) > 0 {
		// start DSL server
		vm := js.NewVirtualMachine()
		srv := js.NewLanguageServer(vm, *httpPort)
		go srv.Start()
	}

	// start REPL
	line := liner.NewLiner()
	defer line.Close()
	defer tearDown(line)
	setup(line)
	loop(line)
}

func welcome() {
	fmt.Println("\033[1;34mmelrose\033[0m" + fmt.Sprintf(" - %s - syntax %s", version, dsl.Syntax))
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
	line.SetWordCompleter(completeMe)
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
		entry = strings.TrimSpace(entry)
		if strings.HasPrefix(entry, ":") {
			// special case
			if entry == ":q" || entry == ":Q" {
				goto exit
			}
			args := strings.Split(entry, " ")
			if cmd, ok := lookupCommand(args[0]); ok {
				if msg := cmd.Func(args[1:]); msg != nil {
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
	dsl.StopAllLoops(varStore)
}

func processInputFile(fileName string) {
	data, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		notify.Print(notify.Errorf("unable to read file:%v", err))
		return
	}
	for line, each := range strings.Split(string(data), "\n") {
		entry := strings.TrimSpace(each)
		if err := dispatch(entry); err != nil {
			notify.Print(notify.Errorf("line %d:%v", line, err))
		}
	}
}
